package odoo

import (
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/httpd"
	"gopkg.in/resty.v1"
	"sync/atomic"
	"time"
	"fmt"
	"net/http"
	"github.com/pkg/errors"
	"strconv"
	"encoding/json"
)

type Diagnostic interface {
	Error(msg string, err error)
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
	diag        Diagnostic
	methods	Methods
	DB	   *storage.Service
	Httpd  *httpd.Service
	httpClient  *resty.Client
	endpoints   []*Endpoint
	configValue atomic.Value
}

func NewService(c Config, d Diagnostic) *Service {
	e, _ := c.index()
	s := &Service{
		diag:      d,
		methods: Methods{},
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

	var r httpd.Route

	r = httpd.Route{
		RouteType:	httpd.ROUTE_TYPE_HTTP,
		Method:  "GET",
		Pattern: "/results",
		HandlerFunc: s.methods.getResults,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:	httpd.ROUTE_TYPE_HTTP,
		Method:  "PATCH",
		Pattern: "/results/{id:int}",
		HandlerFunc: s.methods.patchResult,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:	httpd.ROUTE_TYPE_HTTP,
		Method:  "POST",
		Pattern: "/workorders",
		HandlerFunc: s.methods.postWorkorders,
	}
	s.Httpd.Handler[0].AddRoute(r)

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) GetWorkorder(masterpa_sn string, hmi_sn string, code string) ([]byte, error) {

	var err error
	var body []byte
	for _, endpoint := range s.endpoints {
		body, err = s.getWorkorder(fmt.Sprintf(endpoint.url, masterpa_sn, hmi_sn, code), endpoint.method)
		if err == nil {
			// 如果第一次就成功，推出循环
			return body, nil
		}
	}
	return nil, errors.Wrap(err, "Get workorder fail")
}

func (s *Service) getWorkorder(url string , method string) ([]byte, error) {
	r := s.httpClient.R()
	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = r.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Get workorder fail: %s", err.Error())
		} else {
			if resp.StatusCode() != http.StatusOK {
				return nil, fmt.Errorf("Get workorder fail: %d", resp.StatusCode())
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
	var final_err error = nil
	var db_workorders []storage.Workorders
	for _, v := range workorders {
		o := storage.Workorders{}
		o.Status = v.Status
		o.WorkorderID = v.ID
		o.PSet, _ = strconv.Atoi(v.PSet)
		o.HMISN = v.HMI.UUID
		o.Knr = v.KNR
		o.LongPin = v.LongPin
		o.NutTotal = v.NutTotal
		o.Vin = v.VIN
		o.MaxOpTime = v.Max_op_time
		o.MaxRedoTimes = v.Max_redo_times
		worksheet, _ := json.Marshal(v.Worksheet)
		o.WorkSheet = string(worksheet)

		ids, _ := json.Marshal(v.Result_IDs)
		o.ResultIDs = string(ids)

		e := s.DB.InsertWorkorder(o)
		if e != nil {
			final_err = e
		}

		db_workorders = append(db_workorders, o)
	}

	return db_workorders, final_err
}