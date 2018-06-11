package odoo

import (
	"errors"
	"sync/atomic"
	"time"

	"gopkg.in/resty.v1"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	httpClient  *resty.Client
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag: d,
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

func (s *Service) CreateMO() error {
	if s.httpClient == nil {
		return errors.New("Odoo Http client is nil")
	}

	return nil
}
