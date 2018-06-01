package aiis

import (
	"errors"
	"sync/atomic"
	"time"

	"gopkg.in/resty.v1"
	"fmt"
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
	configValue atomic.Value
	diag        Diagnostic
	endpoints   []*Endpoint
	httpClient  *resty.Client
}

func NewService(c Config, d Diagnostic) *Service {
	e, _ := c.index()
	s := &Service{
		diag:      d,
		endpoints: e,
	}
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
		SetRetryWaitTime(time.Duration(c.PushInterval)).
		SetRetryMaxWaitTime(20 * time.Second)

	s.httpClient = client
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) ProductionCreate(body interface{}) error {
	for _, endpoint := range s.endpoints {
		err := s.productionCreate(body, endpoint)
		if err == nil {
			// 如果第一次就成功，推出循环
			break
		}
	}
	return nil
}

func (s *Service) productionCreate(body interface{},endpoint *Endpoint ) error {
	r := s.httpClient.R().SetBody(body)
	switch endpoint.method {
	case "PATCH":
		_, err := r.Patch(endpoint.url)
		if err != nil {
			return fmt.Errorf("Create Production fail : %s ", err.Error())
		}
	case "PUT":
		_, err := r.Put(endpoint.url)
		if err != nil {
			return fmt.Errorf("Create Production fail : %s ", err.Error())
		}
	case "POST":
		_, err := r.Post(endpoint.url)
		if err != nil {
			return fmt.Errorf("Create Production fail : %s ", err.Error())
		}
	default:
		return errors.New("Create Production :the Method is wrong ")

	}
	return nil
}
