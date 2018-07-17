package odoo

import (
	"fmt"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/storage"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	CreateWOSuccess(id int64)
}

type Endpoint struct {
	url     string
	headers map[string]string
	method  string
}

func NewEndpoint(url string, headers map[string]string, method string) *Endpoint {
	return &Endpoint{
		url:     url,
		headers: headers,
		method:  method,
	}
}

type Service struct {
	diag         Diagnostic
	methods      Methods
	DB           *storage.Service
	HTTPDService interface {
		GetHandlerByName(version string) (*httpd.Handler, error)
	}
	httpClient  *resty.Client
	endpoints   []*Endpoint
	configValue atomic.Value
}

func NewService(c Config, d Diagnostic) *Service {
	e, _ := c.index()
	s := &Service{
		diag:      d,
		methods:   Methods{},
		endpoints: e,
	}

	s.methods.service = s
	s.configValue.Store(c)

	return s
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
	client.SetHeaders(c.Headers)
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

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) GetWorkorder(masterpa_sn string, hmi_sn string, code string) ([]byte, error) {

	var err error
	var body []byte
	for _, endpoint := range s.endpoints {
		body, err = s.getWorkorder(fmt.Sprintf(endpoint.url, hmi_sn, code), endpoint.method)
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

func (s *Service) CreateWorkorders(workorders []ODOOWorkorder) ([]storage.Workorders, error) {
	var finalErr error = nil
	dbWorkorders := make([]storage.Workorders, len(workorders))

	for i, v := range workorders {
		if len(v.Consumes) == 0 {
			// 忽略没有消耗品的工单
			continue
		}

		o := storage.Workorders{}
		o.Status = v.Status
		o.WorkorderID = v.ID
		o.HMISN = v.HMI.UUID
		o.Knr = v.KNR
		o.LongPin = v.LongPin
		o.Vin = v.VIN
		o.MaxOpTime = v.Max_op_time
		o.WorkSheet = v.Worksheet
		o.UpdateTime = v.UpdateTime
		o.JobID, finalErr = strconv.Atoi(v.Job)
		if finalErr != nil {
			o.JobID = 0
		}

		o.MO_Year = v.MO_Year
		o.MO_Pin_check_code = v.MO_Pin_check_code
		o.MO_Pin = v.MO_Pin
		o.MO_FactoryName = v.MO_FactoryName
		o.MO_AssemblyLine = v.MO_AssemblyLine
		o.MO_EquipemntName = v.MO_EquipemntName
		o.MO_Lnr = v.MO_Lnr
		o.MO_Model = v.MO_Model

		results := []storage.Results{}
		result_count := 0
		for _, consu := range v.Consumes {
			if len(consu.ResultIDs) == 0 {
				// 忽略没有结果的消耗品
				continue
			}

			r := storage.Results{}
			r.ControllerSN = consu.ControllerSN
			r.GunSN = consu.GunSN
			r.PSet, _ = strconv.Atoi(consu.PSet)
			r.ToleranceMax = consu.ToleranceMax
			r.ToleranceMin = consu.ToleranceMin
			r.ToleranceMaxDegree = consu.ToleranceMaxDegree
			r.ToleranceMinDegree = consu.ToleranceMinDegree
			r.NutNo = consu.NutNo

			r.WorkorderID = o.WorkorderID
			r.Result = audi_vw.RESULT_NONE
			r.HasUpload = false
			r.Stage = audi_vw.RESULT_STAGE_INIT
			r.UpdateTime = time.Now()
			r.PSetDefine = ""
			r.ResultValue = ""
			r.Count = 1

			for _, result_id := range consu.ResultIDs {
				result_count++

				r.OffsetX = consu.X
				r.OffsetY = consu.Y

				r.Seq = result_count
				r.ResultId = result_id
				r.MaxRedoTimes = consu.Max_redo_times
				results = append(results, r)
			}
		}

		o.LastResultID = results[len(results)-1].ResultId

		e := s.DB.InsertWorkorder(&o, &results)
		if e != nil {
			finalErr = e
		}

		dbWorkorders[i] = o
	}

	return dbWorkorders, finalErr
}
