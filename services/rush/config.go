package rush

import (
	"github.com/masami10/aiis/toml"
	"time"
)

const (
	DEFAULT_URL = "http://127.0.0.1:8082"
)

type Config struct {
	Enable  bool `yaml:"enable"`
	Workers int  `yaml:"workers"`
	URL          string            `yaml:"url"`
	Route        string            `yaml:"route"`
	Headers      map[string]string `yaml:"headers" override:"headers"`
	Timeout      toml.Duration     `yaml:"timeout"`
	PushInterval toml.Duration     `yaml:"push_interval"`
	MaxRetry     int               `yaml:"max_retry"` // api最大尝试次数
}

func NewConfig() Config {
	return Config{
		Enable:  true,
		Workers: 5,
		URL:          DEFAULT_URL,
		Route:        "/rush/v1/results", //方法永远为patch
		Timeout:      toml.Duration(time.Millisecond * 10),
		PushInterval: toml.Duration(time.Second * 1),
		MaxRetry:     3,
		Headers:      map[string]string{"Content-Type": "application/json"},
	}
}

func (c Config) Validate() error {
	return nil
}
