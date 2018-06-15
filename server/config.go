package server

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/masami10/aiis/command"
	"github.com/masami10/aiis/services/diagnostic"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/storage"

	"github.com/masami10/aiis/services/fis"
	"github.com/masami10/aiis/services/odoo"
	"github.com/masami10/aiis/services/pmon"
	"github.com/masami10/aiis/services/rush"
	"github.com/pkg/errors"
)

type Config struct {
	Hostname string `yaml:"hostname"`
	DataDir  string `yaml:"data_dir"`

	Logging diagnostic.Config `yaml:"logging"`

	HTTP httpd.Config `yaml:"httpd"`

	Pmon pmon.Config `yaml:"pmon"`

	Fis fis.Config `yaml:"fis"`

	Odoo odoo.Config `yaml:"odoo"`

	Storage storage.Config `yaml:"storage"`

	Rush rush.Config `yaml:"rush"`

	Commander command.Commander `yaml:"-"`
}

func NewConfig() *Config {
	c := &Config{
		Hostname:  "localhost",
		Commander: command.ExecCommander,
	}

	c.HTTP = httpd.NewConfig()
	c.Pmon = pmon.NewConfig()
	c.Fis = fis.NewConfig()
	c.Logging = diagnostic.NewConfig()
	c.Odoo = odoo.NewConfig()
	c.Storage = storage.NewConfig()
	c.Rush = rush.NewConfig()

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

	if err := c.Pmon.Validate(); err != nil {
		return errors.Wrap(err, "pmon")
	}

	if err := c.Fis.Validate(); err != nil {
		return errors.Wrap(err, "fis")
	}

	if err := c.Odoo.Validate(); err != nil {
		return errors.Wrap(err, "odoo")
	}

	if err := c.Storage.Validate(); err != nil {
		return errors.Wrap(err, "storage")
	}

	if err := c.Rush.Validate(); err != nil {
		return errors.Wrap(err, "rush")
	}

	return nil
}
