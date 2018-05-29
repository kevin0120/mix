package server

import (
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/diagnostic"
	"github.com/masami10/aiis/command"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration format for the aiisd binary.
type Config struct {
	HTTP           	httpd.Config      `yaml:"http"`
	Logging        diagnostic.Config  `yaml:"logging"`
	Commander 	   	command.Commander `yaml:"-"`
}

func NewConfig() *Config {
	c := &Config{
		Commander: command.ExecCommander,
	}

	c.HTTP = httpd.NewConfig()
	c.Logging = diagnostic.NewConfig()
	
	return c
}

func NewDemoConfig() (*Config, error) {
	c := NewConfig()

	return c, nil
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if err := c.HTTP.Validate(); err != nil {
		return errors.Wrap(err, "http")
	}

	return nil
}