package broker

import (
	"github.com/pkg/errors"
	"strings"
)

type Config struct {
	ConnectUrls []string `yaml:"connect-urls"`
	Enable      bool `yaml:"enable"`
	Provider string `yaml:"provider"`
}

func NewConfig() Config {
	c := Config{
		Enable:      false,
		ConnectUrls: []string{"nats://127.0.0.1:4222"},
		Provider: "nats",
	}

	return c
}

func (c Config) Validate() error {
	if c.Enable {
		for _, connect := range c.ConnectUrls {
			if !strings.HasPrefix(connect, "nats://") {
				return errors.Errorf("Broker Connect Url Is Invalid, error: %s", connect)
			}
		}
		if c.Provider != "nats"{
			return errors.Errorf("Broker Provider: %s Is Not Support", c.Provider)
		}
	}
	return nil
}
