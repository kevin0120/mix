package server

import (
	"github.com/masami10/aiis/command"
	"os/user"
	"os"
	"fmt"
	"path/filepath"
	"github.com/masami10/aiis/services/diagnostic"
	"github.com/masami10/aiis/services/httpd"

	"github.com/pkg/errors"
)



type Config struct {
	Hostname               string `yaml:"hostname"`
	DataDir                string `yaml:"data_dir"`

	Logging        			diagnostic.Config `yaml:"logging"`

	HTTP 					httpd.Config `yaml:"httpd"`

	Commander command.Commander `yaml:"-"`
}


func NewConfig() *Config {
	c := &Config{
		Hostname:  "localhost",
		Commander: command.ExecCommander,
	}

	c.HTTP = httpd.NewConfig()
	c.Logging = diagnostic.NewConfig()

	return c
}

func NewDemoConfig() (*Config, error) {
	c := NewConfig()

	var homeDir string
	// By default, store meta and data files in current users home directory
	u, err := user.Current()
	if err == nil {
		homeDir = u.HomeDir
	} else if os.Getenv("HOME") != "" {
		homeDir = os.Getenv("HOME")
	} else {
		return nil, fmt.Errorf("failed to determine current user for storage")
	}

	c.DataDir = filepath.Join(homeDir, ".aiis", c.DataDir)

	return c, nil
}

func (c *Config) Validate() error {
	if c.Hostname == "" {
		return fmt.Errorf("must configure valid hostname")
	}
	if c.DataDir == "" {
		return fmt.Errorf("must configure valid data dir")
	}

	if err := c.HTTP.Validate(); err != nil {
		return errors.Wrap(err, "http")
	}

	return nil
}
