package odoo

import (
	"errors"
	"fmt"
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	Urls []string `yaml:"urls"`
	//Headers         map[string]string `yaml:"headers" override:"headers"`
	//WorkorderGetUrl string            `yaml:"workorder_get_url"`
	//Method          string            `yaml:"method"`
	Timeout   toml.Duration  `yaml:"timeout"`
	Interval  toml.Duration  `yaml:"interval"`
	MaxRetry  int            `yaml:"max_retry"` // api最大尝试次数
	Workers   int            `yaml:"workers"`
	Endpoints []EndpointConf `yaml:"endpoints"`
}

type EndpointConf struct {
	Name    string            `yaml:"name"`
	Path    string            `yaml:"path"`
	Headers map[string]string `yaml:"headers" override:"headers"`
	Method  string            `yaml:"method"`
}

func NewConfig() Config {
	ec := EndpointConf{
		Name:    "test",
		Path:    "/api/v1/test",
		Headers: map[string]string{"Content-Type": "application/json"},
		Method:  "GET",
	}

	c := Config{
		Urls: []string{"http://127.0.0.1:8069"},
		//WorkorderGetUrl: "/api/v1/mrp.workorders?masterpc=%s&hmi=%s&code=%s&limit=1&order=production_date",
		Timeout:  toml.Duration(time.Millisecond * 10),
		Interval: toml.Duration(time.Second * 1),
		//Method:          "GET",
		MaxRetry: 3,
		//Headers:         map[string]string{"Content-Type": "application/json"},
		Endpoints: []EndpointConf{ec},
	}

	return c
}

func (c Config) Validate() error {
	if len(c.Urls) == 0 {
		return errors.New("ODOO service URLs is empty ")
	}

	for _, v := range c.Endpoints {
		if v.Path[0] != '/' {
			return fmt.Errorf("route patterns must begin with a '/' %s", v.Path[0])
		}
	}

	return nil
}

func (c Config) index() ([]*Endpoint, error) {
	m := []*Endpoint{}

	for _, v := range c.Endpoints {
		for _, url := range c.Urls {
			m = append(m, NewEndpoint(url+v.Path, v.Headers, v.Method, v.Name))
		}
	}

	return m, nil
}
