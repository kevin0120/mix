package broker

import (
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/pkg/errors"
	"strings"
)

var DefaultOptions = map[string]string{"name": "Dummy Server", "token": "Dummy Token"}

type Config struct {
	ConnectUrls []string          `yaml:"connect-urls"`
	Enable      bool              `yaml:"enable"`
	Provider    string            `yaml:"provider"`
	Options     map[string]string `yaml:"connect-options"`
	Name        string            `yaml:"name"`
}

func NewConfig() Config {
	n, _ := uuid.NewV4()
	c := Config{
		Enable:      false,
		ConnectUrls: []string{"nats://127.0.0.1:4222"},
		Provider:    "nats",
		Options:     DefaultOptions,
		Name:        n.String(),
	}

	return c
}

func (c Config) Validate() error {
	if c.Enable {
		if c.Name == "" {
			return errors.New("broker Name Is Required! But Now Is Empty")
		}
		for _, connect := range c.ConnectUrls {
			if !strings.HasPrefix(connect, "nats://") {
				return errors.Errorf("broker connect Url Is Invalid, error: %s", connect)
			}
		}
		if c.Provider != "nats" {
			return errors.Errorf("broker Provider: %s Is Not Support", c.Provider)
		}
	}
	return nil
}
