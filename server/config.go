package server

import (
	"fmt"
	"github.com/masami10/rush/command"
	"github.com/masami10/rush/services/diagnostic"
	"github.com/masami10/rush/services/httpd"
	"os"
	"os/user"
	"path/filepath"

	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Config struct {
	Hostname string `yaml:"hostname"`
	DataDir  string `yaml:"data_dir"`
	SN       string `yaml:"serial_no"`

	DocPath string `yaml:"doc_path"`

	Logging diagnostic.Config `yaml:"logging"`

	HTTP httpd.Config `yaml:"httpd"`

	Minio minio.Config `yaml:"minio"`

	Aiis aiis.Config `yaml:"aiis"`

	Ws wsnotify.Config `yaml:"websocket"`

	Odoo odoo.Config `yaml:"odoo"`

	Storage storage.Config `yaml:"storage"`

	AudiVW audi_vw.Config `yaml:"audi/vw"`

	OpenProtocol openprotocol.Config `yaml:"openprotocol"`

	Contollers controller.Config `yaml:"controller_service"`

	Commander command.Commander `yaml:"-"`
}

func NewConfig() *Config {
	sn, _ := uuid.NewV4()
	c := &Config{
		Hostname:  "localhost",
		SN:        sn.String(),
		Commander: command.ExecCommander,
	}

	c.HTTP = httpd.NewConfig()
	c.Minio = minio.NewConfig()
	c.Aiis = aiis.NewConfig()
	c.Ws = wsnotify.NewConfig()
	c.Storage = storage.NewConfig()
	c.Logging = diagnostic.NewConfig()
	c.AudiVW = audi_vw.NewConfig()
	c.OpenProtocol = openprotocol.NewConfig()
	c.Odoo = odoo.NewConfig()

	c.Contollers = controller.NewConfig()

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

	c.DocPath = filepath.Join(homeDir, "doc", c.DocPath)
	c.DataDir = filepath.Join(homeDir, ".rush", c.DataDir)

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

	if err := c.Minio.Validate(); err != nil {
		return errors.Wrap(err, "minio")
	}

	if err := c.Aiis.Validate(); err != nil {
		return errors.Wrap(err, "aiis")
	}

	if err := c.Odoo.Validate(); err != nil {
		return errors.Wrap(err, "odoo")
	}

	if err := c.Ws.Validate(); err != nil {
		return errors.Wrap(err, "websocket")
	}

	if err := c.Storage.Validate(); err != nil {
		return errors.Wrap(err, "storage")
	}

	if err := c.Contollers.Validate(); err != nil {
		return errors.Wrap(err, "controller")
	}

	return nil
}
