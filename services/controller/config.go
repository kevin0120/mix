package controller

import (
	"fmt"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type ControllerConfig struct {
	SN          string `yaml:"serial_no"`
	Protocol    string `yaml:"protocol"`
	RemoteIP    string `yaml:"remote_ip"`
	Port        int    `yaml:"port"`
	ToolChannel int    `yaml:"channel"`
}

type Config struct {
	Workers		int `yaml:"workers"`
	Configs     []ControllerConfig `yaml:"controllers"`
}

var Protocols []string

func init() {
	Protocols = []string{AUDIPROTOCOL, OPENPROTOCOL}
}

func NewConfig() Config {
	_sn, _ := uuid.NewV4()
	configs := []ControllerConfig{}

	configs = append(configs, ControllerConfig {
			SN:          _sn.String(),
			Protocol:    AUDIPROTOCOL,
			RemoteIP:    "127.0.0.1",
			Port:        4700,
			ToolChannel: DEFAULT_TOOL_CHANNEL,
		})

	return Config {
		Workers: 4,
		Configs:configs,
	}
}

func (c ControllerConfig) Validate() error {
	if c.SN == "" {
		return errors.New("Controller Serial Number must be configuration")
	}
	if c.RemoteIP == "" {
		return errors.New("Controller Remote IP must be configuration")
	}
	if c.Protocol == "" {
		return errors.New("Controller Protocol must be configuration")
	}
	if !utils.StringInSlice(c.Protocol, Protocols) {
		return fmt.Errorf("Protocol %s not in support Protocols: %s ", c.Protocol, Protocols)
	}
	return nil
}

// Configs is the configuration for all [[alertpost]] sections of the kapacitor
// configuration file.
//type Configs []Config

func (cs Config) Validate() error {

	for _, c := range cs.Configs {
		err := c.Validate()
		return errors.Wrap(err, "Validate Controllers Protocol Fail")
	}

	return nil
}
