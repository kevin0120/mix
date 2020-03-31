package odoo

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"gopkg.in/resty.v1"
	"net/http"
	"time"
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
	diag              Diagnostic
	dispatcherBus     Dispatcher
	dispatcherMap     dispatcherbus.DispatcherMap
	storageService    IStorageService
	httpd             HTTPService
	httpClient        *resty.Client
	endpoints         []*Endpoint
	configValue       atomic.Value
	status            string
	workordersChannel chan interface{}
	closing           chan struct{}
}

func NewService(c Config, d Diagnostic, dp Dispatcher, storage IStorageService, httpd HTTPService) *Service {
	e, _ := c.endpoints()
	s := &Service{
		diag:              d,
		endpoints:         e,
		status:            utils.STATUS_OFFLINE,
		workordersChannel: make(chan interface{}, c.Workers),
		closing:           make(chan struct{}, 1),
		dispatcherBus:     dp,
		storageService:    storage,
		httpd:             httpd,
	}

	s.configValue.Store(c)

	s.setupGlbDispatcher()
	s.ensureHttpClient()
	s.setupHttpHandlers()

	return s
}

func (s *Service) GetEndpoints(name string) []Endpoint {

	var endpoints []Endpoint
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

func (s *Service) ensureHttpClient() *resty.Client {
	if s.httpClient != nil {
		return s.httpClient
	}
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
	return client
}

func (s *Service) setupGlbDispatcher() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DispatcherMaintenanceInfo: utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherOrderNew:        utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherToolMaintenance: utils.CreateDispatchHandlerStruct(s.onToolMaintenance),
	}
}

func (s *Service) Open() error {
	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)
	go s.taskSaveWorkorders()
	return nil
}

func (s *Service) Close() error {
	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	s.closing <- struct{}{}
	return nil
}

func (s *Service) setupHttpHandlers() {

	r := httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/workorders",
		HandlerFunc: s.postWorkorders,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/mrp.routing.workcenter",
		HandlerFunc: s.putSyncRoutingOpertions,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/mrp.routing.workcenter.delete",
		HandlerFunc: s.deleteRoutingOpertions,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/maintenance",
		HandlerFunc: s.postMaintenance,
	}
	s.httpd.AddNewHttpHandler(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "DELETE",
		Pattern:     "/mrp.routing.workcenter/all",
		HandlerFunc: s.deleteAllRoutingOpertions,
	}
	s.httpd.AddNewHttpHandler(r)
}

func (s *Service) handleWorkorder(data []byte) {
	s.workordersChannel <- data
}

func (s *Service) GetWorkorder(hmiSn string, workcenterCode, code string) ([]byte, error) {

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
			s.handleWorkorder(body)
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
			return nil, fmt.Errorf("Get Workorder Fail: %s ", err.Error())
		} else {
			status := resp.StatusCode()
			if status != http.StatusOK {
				return nil, fmt.Errorf("Get Workorder Fail: %d ", status)
			} else {
				return resp.Body(), nil
			}
		}
	default:
		return nil, errors.New("Get Workorder :The Method Is Wrong")

	}
}

func (s *Service) taskSaveWorkorders() {
	for {
		select {
		case payload := <-s.workordersChannel:

			orderOut, err := s.handleSaveWorkorders(payload)
			if err != nil {
				s.diag.Error("Save Workorder Failed", err)
				break
			}

			s.doDispatch(dispatcherbus.DispatcherOrderNew, orderOut)

		case <-s.closing:
			s.diag.Info("taskSaveWorkorders closed")
			return
		}
	}
}

func (s *Service) handleSaveWorkorders(payload interface{}) (interface{}, error) {

	code, err := s.storageService.WorkorderIn(payload.([]byte))
	if err != nil {
		return nil, err
	}

	orderOut, _ := s.storageService.WorkorderOut(code, 0)
	return orderOut, nil
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
				return fmt.Errorf("CreateMaintenance: %s", err.Error())
			} else {
				status := resp.StatusCode()
				if status > 400 {
					return fmt.Errorf("CreateMaintenance return status code: %d", status)
				} else if status == http.StatusCreated {
					s.diag.Info("create Maintenance")
					return nil
				} else {
					return nil
				}
			}
		default:
			return errors.New("CreateMaintenance :the Method is wrong")
		}

	}
	return err
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error("Dispatch Failed", err)
	}
}

// 收到工具保养通知
func (s *Service) onToolMaintenance(data interface{}) {
	if data == nil {
		return
	}

	ti := data.(tightening_device.ToolMaintenanceInfo)
	if err := s.TryCreateMaintenance(ti); err != nil {
		s.diag.Error("Create Maintenance Failed ", err)
	}
}

func (s *Service) GetUserByUID(uid string) (*OdooUser, error) {
	r := s.httpClient.R()
	url := fmt.Sprintf("%s/api/v1/res.users/%s", s.Config().Urls[0], uid)
	resp, err := r.Get(url)

	user := OdooUser{}
	if err != nil {
		return nil, err
	} else {
		if resp.StatusCode() != http.StatusOK {
			return nil, errors.New(string(resp.Body()))
		} else {
			_ = json.Unmarshal(resp.Body(), &user)
			return &user, nil
		}
	}
}
