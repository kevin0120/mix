package odoo

import (
	"errors"
	"sync/atomic"
	"time"

	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	route       string
	configValue atomic.Value
	diag        Diagnostic
	httpClient  *resty.Client
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag:  d,
		route: c.URL + c.Route,
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

func (s *Service) CreateMO(body interface{}) error {
	if s.httpClient == nil {
		return errors.New("Odoo Http client is nil ")
	}
	r := s.httpClient.R().SetBody(body)

	resp, err := r.Post(s.route)
	if err != nil {
		return fmt.Errorf("Create MO Post fail: %s ", err)
	} else {
		if resp.StatusCode() != http.StatusCreated {
			return fmt.Errorf("Create MO Post fail: %s ", resp.Status())
		}
	}

	return nil
}
