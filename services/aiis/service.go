package aiis

import (
	"sync/atomic"
	"time"

	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
)

type Diagnostic interface {
	Error(msg string, err error)
	PutResultDone()
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
	rush_port string
}

func NewService(c Config, d Diagnostic, rush_port string) *Service {
	e, _ := c.index()
	s := &Service{
		diag:      d,
		endpoints: e,
		rush_port: rush_port,
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

func (s *Service) PutResult(result_id int64, body interface{}) error {

	var err error
	for _, endpoint := range s.endpoints {
		err = s.putResult(body, fmt.Sprintf(endpoint.url, result_id), endpoint.method)
		if err == nil {
			// 如果第一次就成功，推出循环
			return nil
		}
	}
	return errors.Wrap(err, "Put result fail")
}

func (s *Service) putResult(body interface{}, url string, method string) error {
	r := s.httpClient.R().SetBody(body).SetHeader("rush_port", s.rush_port)
	var resp *resty.Response
	var err error

	switch method {
	case "PATCH":
		resp, err = r.Patch(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	case "PUT":
		resp, err = r.Put(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	case "POST":
		resp, err = r.Post(url)
		if err != nil {
			return fmt.Errorf("Result Put fail: %s ", err.Error())
		} else {
			if resp.StatusCode() != http.StatusNoContent {
				return fmt.Errorf("Result Put fail: %d ", resp.StatusCode())
			}
		}
	default:
		return errors.New("Result Put :the Method is wrong")

	}
	s.diag.PutResultDone()
	return nil
}
