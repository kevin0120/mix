package odoo

import (
	"fmt"

	"github.com/masami10/aiis/toml"
	"time"
)

const (
	DEFAULT_URL = "127.0.0.1:8069"
)

type Config struct {
	URL          string            `yaml:"url"`
	Route        string            `yaml:"route"`
	Headers      map[string]string `yaml:"headers" override:"headers"`
	Timeout      toml.Duration     `yaml:"timeout"`
	PushInterval toml.Duration     `yaml:"push_interval"`
	MaxRetry     int               `yaml:"max_retry"` // api最大尝试次数
	Enable       bool              `yaml:"enable"`
}

func NewConfig() Config {
	return Config{
		URL:          DEFAULT_URL,
		Route:        "/api/v1/mrp.productions", //方法永远为post
		Enable:       true,
		Timeout:      toml.Duration(time.Millisecond * 10),
		PushInterval: toml.Duration(time.Second * 1),
		MaxRetry:     3,
		Headers:      map[string]string{"Content-Type": "application/json"},
	}
}

func (c Config) Validate() error {
	if c.Enable && c.URL == "" {
		return fmt.Errorf("Odoo URL is empty ")
	}
	if c.Enable && c.Route[0] != '/' {
		return fmt.Errorf("route patterns must begin with a '/' %s", c.Route)
	}
	return nil
}
