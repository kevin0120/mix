package odoo

import (
	"errors"
	"sync/atomic"
	"time"

	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
)

const (
	HEALTHZ_ITV    = 3 // 秒
	STATUS_ONLINE  = "online"
	STATUS_OFFLINE = "offline"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type OnStatus func(status string)

type Service struct {
	route       string
	configValue atomic.Value
	diag        Diagnostic
	httpClient  *resty.Client
	status      string
	OnStatus    OnStatus
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag:   d,
		route:  c.URL + c.Route,
		status: STATUS_OFFLINE,
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

	go s.taskHealthz()

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
		return fmt.Errorf("Create MO Post fail: %s\n", err)
	} else {
		if resp.StatusCode() != http.StatusCreated {
			return fmt.Errorf("Create MO Post fail: %s\n", resp.Status())
		}
	}

	return nil
}

func (s *Service) Healthz() error {
	if s.httpClient == nil {
		return errors.New("Odoo Http client is nil ")
	}
	r := s.httpClient.R()

	url := fmt.Sprintf("%s/api/v1/healthz", s.Config().URL)
	_, err := r.Get(url)
	if err != nil {
		return fmt.Errorf("check status failed", err)
	}

	return nil
}

func (s *Service) UpdateStatus(status string) {
	if s.status != status {
		s.status = status

		// 推送状态
		if s.OnStatus != nil {
			s.OnStatus(s.status)
		}
	}
}

func (s *Service) taskHealthz() {
	for {
		err := s.Healthz()
		if err != nil {
			// 离线
			s.UpdateStatus(STATUS_OFFLINE)
		} else {
			// 在线
			s.UpdateStatus(STATUS_ONLINE)
		}

		time.Sleep(HEALTHZ_ITV * time.Second)
	}
}

func (s *Service) Status() string {
	return s.status
}
