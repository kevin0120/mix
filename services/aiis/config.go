package aiis

import (
	"errors"
	"fmt"
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	Urls          []string          `yaml:"urls"`
	Headers       map[string]string `yaml:"headers" override:"headers"`
	ResultPushUrl string            `yaml:"result_push_url"`
	PushMethod    string            `yaml:"push_method"`
	Timeout       toml.Duration     `yaml:"timeout"`
	PushInterval  toml.Duration     `yaml:"push_interval"`
	MaxRetry      int               `yaml:"max_retry"` // api最大尝试次数
}

func NewConfig() Config {
	c := Config{
		Urls:          []string{"http://127.0.0.1:9092"},
		ResultPushUrl: "/aiis/v1/operation.results/%d",
		Timeout:       toml.Duration(time.Millisecond * 10),
		PushInterval:  toml.Duration(time.Second * 1),
		PushMethod:    "PUT",
		MaxRetry:      5,
		Headers:       map[string]string{"Content-Type": "application/json"},
	}

	return c
}

func (c Config) Validate() error {
	if len(c.Urls) == 0 {
		return errors.New("Aiis service URLs is empty ")
	}
	if c.ResultPushUrl[0] != '/' {
		return fmt.Errorf("route patterns must begin with a '/' %s", c.ResultPushUrl)
	}

	m := []string{"PATCH", "PUT", "POST"}

	for _, v := range m {
		if c.PushMethod != v {
			break
		}
		return fmt.Errorf("PushMethod: %s not in  %s", c.PushMethod, m)
	}

	return nil
}

func (c Config) index() ([]*Endpoint, error) {
	m := []*Endpoint{}

	for _, url := range c.Urls {
		m = append(m, NewEndpoint(url+c.ResultPushUrl, c.Headers, c.PushMethod))
	}

	return m, nil
}
