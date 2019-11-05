package broker

import (
	"github.com/pkg/errors"
	"strings"
)

var DefaultOptions = map[string]string{"name": "Dummy Server", "token": "Dummy Token"}

type Config struct {
	ConnectUrls []string          `yaml:"connect-urls"`
	Enable      bool              `yaml:"enable"`
	Provider    string            `yaml:"provider"`
	Options     map[string]string `yaml:"connect-options"`
}

func NewConfig() Config {
	c := Config{
		Enable:      false,
		ConnectUrls: []string{"nats://127.0.0.1:4222"},
		Provider:    "nats",
		Options:     DefaultOptions,
	}

	return c
}

func (c Config) Validate() error {
	if c.Enable {
		for _, connect := range c.ConnectUrls {
			if !strings.HasPrefix(connect, "nats://") {
				return errors.Errorf("Broker connect Url Is Invalid, error: %s", connect)
			}
		}
		if c.Provider != "nats" {
			return errors.Errorf("Broker Provider: %s Is Not Support", c.Provider)
		}
	}
	return nil
}
