package odoo

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	ODOO_STATUS_ONLINE  = "online"
	ODOO_STATUS_OFFLINE = "offline"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
	CreateWOSuccess(id int64)
}

type Endpoint struct {
	url     string
	headers map[string]string
	method  string
	name    string
}

func NewEndpoint(url string, headers map[string]string, method string, name string) *Endpoint {
	return &Endpoint{
		url:     url,
		headers: headers,
		method:  method,
		name:    name,
	}
}

type Service struct {
	diag         Diagnostic
	methods      Methods
	DB           *storage.Service
	HTTPDService interface {
		GetHandlerByName(version string) (*httpd.Handler, error)
	}
	httpClient        *resty.Client
	endpoints         []*Endpoint
	configValue       atomic.Value
	status            string
	WS                *wsnotify.Service
	Aiis              *aiis.Service
	workordersChannel chan interface{}
	wg                sync.WaitGroup
	closing           chan struct{}
}

func NewService(c Config, d Diagnostic) *Service {
	e, _ := c.index()
	s := &Service{
		diag:              d,
		methods:           Methods{},
		endpoints:         e,
		status:            ODOO_STATUS_OFFLINE,
		workordersChannel: make(chan interface{}, c.Workers),
		closing:           make(chan struct{}),
	}

	s.methods.service = s
	s.configValue.Store(c)

	return s
}

func (s *Service) UpdateStatus(status string) {
	if s.status != status {
		s.status = status
		s.PushStatus()
	}
}

func (s *Service) PushStatus() {
	odooStatus := wsnotify.WSOdooStatus{
		Status: s.status,
	}

	str, _ := json.Marshal(odooStatus)
	s.WS.WSSend(wsnotify.WS_EVENT_ODOO, string(str))
}

func (s *Service) GetEndpoints(name string) []Endpoint {

	endpoints := []Endpoint{}
	for _, v := range s.endpoints {
		if v.name == name {
			endpoints = append(endpoints, *v)
		}
	}

	return endpoints
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(time.Duration(c.Timeout))
	client.SetContentLength(true)
	// Headers for all request
	client.
		SetRetryCount(c.MaxRetry).
		SetRetryWaitTime(time.Duration(c.Interval)).
		SetRetryMaxWaitTime(20 * time.Second)

	s.httpClient = client

	handler, err := s.HTTPDService.GetHandlerByName(httpd.BasePath)
	if err != nil {
		return errors.Wrap(err, "Odoo server get Httpd default Handler fail")
	}

	r := httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/results",
		HandlerFunc: s.methods.getResults,
	}
	handler.AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PATCH",
		Pattern:     "/results/{id:int}",
		HandlerFunc: s.methods.patchResult,
	}
	handler.AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/workorders",
		HandlerFunc: s.methods.postWorkorders,
	}
	handler.AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/mrp.routing.workcenter",
		HandlerFunc: s.methods.putSyncRoutingOpertions,
	}
	handler.AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/maintenance",
		HandlerFunc: s.methods.postMaintenance,
	}
	handler.AddRoute(r)

	s.Aiis.OnOdooStatus = s.OnStatus
	s.Aiis.SyncGun = s.GetGunID

	for i := 0; i < s.Config().Workers; i++ {
		s.wg.Add(1)
		go s.taskSaveWorkorders()
	}

	return nil
}

func (s *Service) Close() error {
	for i := 0; i < s.Config().Workers; i++ {
		s.closing <- struct{}{}
	}

	s.wg.Wait()

	return nil
}

func (s *Service) GetWorkorder(masterpcSn string, hmiSn string, workcenterCode, code string) ([]byte, error) {

	var err error
	var body []byte
	endpoints := s.GetEndpoints("getWorkorder")
	for _, endpoint := range endpoints {
		url := fmt.Sprintf(endpoint.url, code)
		if hmiSn != "" {
			url += fmt.Sprintf("&hmi=%s", hmiSn)
		}

		if workcenterCode != "" {
			url += fmt.Sprintf("&workcenter=%s", workcenterCode)
		}

		body, err = s.getWorkorder(url, endpoint.method)
		if err == nil {
			// 如果第一次就成功，推出循环
			return body, nil
		}
	}
	return nil, errors.Wrap(err, "Get workorder fail")
}

func (s *Service) getWorkorder(url string, method string) ([]byte, error) {
	r := s.httpClient.R()
	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = r.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Get workorder fail: %s", err.Error())
		} else {
			status := resp.StatusCode()
			if status != http.StatusOK {
				return nil, fmt.Errorf("Get workorder fail: %d", status)
			} else {
				return resp.Body(), nil
			}
		}
	default:
		return nil, errors.New("Get workorder :the Method is wrong")

	}

	return nil, nil
}

func (s *Service) GetGun(serial string) (ODOOGun, error) {

	var err error
	var gun ODOOGun
	endpoints := s.GetEndpoints("getGun")
	for _, endpoint := range endpoints {
		body, err := s.getGun(fmt.Sprintf(endpoint.url, serial), endpoint.method)
		if err == nil {
			// 如果第一次就成功，推出循环
			var guns []ODOOGun
			err = json.Unmarshal(body, &guns)
			if err != nil || len(guns) == 0 {
				return gun, errors.Wrap(err, "Get gun fail")
			}

			return guns[0], nil
		}
	}

	return gun, errors.Wrap(err, "Get gun fail")
}

func (s *Service) GetGunID(serial string) (int64, error) {

	var err error
	//var gun ODOOGun
	endpoints := s.GetEndpoints("getGun")
	for _, endpoint := range endpoints {
		body, err := s.getGun(fmt.Sprintf(endpoint.url, serial), endpoint.method)
		if err == nil {
			// 如果第一次就成功，推出循环
			var guns []ODOOGun
			err = json.Unmarshal(body, &guns)
			if err != nil || len(guns) == 0 {
				return 0, errors.Wrap(err, "Get gun fail")
			}

			return guns[0].ID, nil
		}
	}

	return 0, errors.Wrap(err, "Get gun fail")
}

func (s *Service) getGun(url string, method string) ([]byte, error) {
	r := s.httpClient.R()
	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = r.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Get gun fail: %s", err.Error())
		} else {
			status := resp.StatusCode()
			if status != http.StatusOK {
				return nil, fmt.Errorf("Get gun fail: %d", status)
			} else {
				return resp.Body(), nil
			}
		}
	default:
		return nil, errors.New("Get gun :the Method is wrong")

	}

	return nil, nil
}

func (s *Service) CreateWorkorders(workorders []ODOOWorkorder) ([]storage.Workorders, error) {

	var finalErr error = nil
	dbWorkorders := make([]storage.Workorders, len(workorders))

	for i, v := range workorders {

		o := storage.Workorders{}

		if len(v.Consumes) == 0 {
			// 忽略没有消耗品的工单
			continue
		}

		exist, _ := s.DB.WorkorderExists(v.ID)
		if exist {
			// 忽略已存在的工单
			o, err := s.DB.GetWorkorder(v.ID,false)
			if err != nil {
				continue
			}
			dbWorkorders[i] = o
		}else {
			o.Status = "ready"
			o.WorkorderID = v.ID
			o.HMISN = v.HMI.UUID
			o.WorkcenterCode = v.Workcenter.Code
			o.Knr = v.KNR
			o.LongPin = v.LongPin
			o.Vin = v.VIN
			o.MaxOpTime = v.Max_op_time
			//o.WorkSheet = v.Worksheet
			o.UserID = 1
			o.ImageOPID = v.ImageOPID
			o.VehicleTypeImg = v.VehicleTypeImg
			o.UpdateTime = time.Now()
			o.JobID, _ = strconv.Atoi(v.Job)

			o.MO_Year = v.MO_Year
			o.MO_Pin_check_code = v.MO_Pin_check_code
			o.MO_Pin = v.MO_Pin
			o.MO_FactoryName = v.MO_FactoryName
			o.MO_AssemblyLine = v.MO_AssemblyLine
			o.MO_EquipemntName = v.MO_EquipemntName
			o.MO_Lnr = v.MO_Lnr
			o.MO_Model = v.MO_Model
			sConsumes, _ := json.Marshal(v.Consumes)
			o.Consumes = string(sConsumes)

			//results := []storage.Results{}
			//result_count := 0
			//ignore := false
			//for k, consu := range v.Consumes {
			//	if len(consu.ResultIDs) == 0 {
			//		// 忽略没有结果的消耗品
			//		continue
			//	}
			//
			//	r := storage.Results{}
			//	r.ControllerSN = consu.ControllerSN
			//	r.GunSN = consu.GunSN
			//	r.PSet, _ = strconv.Atoi(consu.PSet)
			//	r.ToleranceMax = consu.ToleranceMax
			//	r.ToleranceMin = consu.ToleranceMin
			//	r.ToleranceMaxDegree = consu.ToleranceMaxDegree
			//	r.ToleranceMinDegree = consu.ToleranceMinDegree
			//	r.NutNo = consu.NutNo
			//	r.Batch = fmt.Sprintf("%d/%d", k+1, len(v.Consumes))
			//
			//	//r.WorkorderID = o.WorkorderID
			//	r.Result = storage.RESULT_NONE
			//	r.HasUpload = false
			//	r.Stage = storage.RESULT_STAGE_INIT
			//	r.UpdateTime = time.Now()
			//	r.PSetDefine = ""
			//	r.ResultValue = ""
			//	r.Count = 1
			//	r.UserID = 1
			//
			//	if len(consu.ResultIDs) == 0 {
			//		ignore = true
			//		break
			//	}
			//
			//	for _, result_id := range consu.ResultIDs {
			//		result_count++
			//
			//		r.OffsetX = consu.X
			//		r.OffsetY = consu.Y
			//
			//		r.Seq = result_count
			//		r.ResultId = result_id
			//		r.MaxRedoTimes = consu.Max_redo_times
			//		results = append(results, r)
			//	}
			//}

			//o.LastResultID = results[len(results)-1].Id

			e := s.DB.InsertWorkorder(&o, nil, true, false, true)
			if e != nil {
				finalErr = e
			}
			dbWorkorders[i] = o
		}

	}

	return dbWorkorders, finalErr
}

func (s *Service) OnStatus(status string) {
	s.UpdateStatus(status)
}

func (s *Service) Status() string {
	return s.status
}

func (s *Service) taskSaveWorkorders() {
	for {
		select {
		case payload := <-s.workordersChannel:
			s.handleSaveWorkorders(payload)

		case <-s.closing:
			s.diag.Info("taskSaveWorkorders closed")
			s.wg.Done()
			return
		}
	}
}

func (s *Service) handleSaveWorkorders(payload interface{}) {
	//defer debug.FreeOSMemory() //快速释放不必要的内存

	workorders := payload.(*[]ODOOWorkorder)
	s.CreateWorkorders(*workorders)
}

func (s *Service) TryCreateMaintenance(body interface{}) error {
	var err error = nil
	endpoints := s.GetEndpoints("tryCreateMaintenance")
	for _, endpoint := range endpoints {
		url := endpoint.url

		r := s.httpClient.R().SetBody(body)
		var resp *resty.Response

		switch endpoint.method {
		case "POST":
			resp, err = r.Post(url)
			if err != nil {
				return fmt.Errorf("TryCreateMaintenance: %s", err.Error())
			} else {
				status := resp.StatusCode()
				if status > 400 {
					return fmt.Errorf("TryCreateMaintenance return status code: %d", status)
				} else if status == http.StatusCreated {
					s.diag.Info("create Maintenance")
					return nil
				} else {
					return nil
				}
			}
		default:
			return errors.New("TryCreateMaintenance :the Method is wrong")
		}

	}
	return err
}
