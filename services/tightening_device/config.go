package tightening_device

import "github.com/masami10/rush/services/controller"

type DeviceConfig struct {
	Model    string `yaml:"model"`
	Protocol string `yaml:"protocol"`
	Endpoint string `yaml:"endpoint"`
}

type Config struct {
	Enable  bool                      `yaml:"enable"`
	Devices []controller.DeviceConfig `yaml:"devices"`
}

func NewConfig() Config {

	return Config{
		Enable:  true,
		Devices: []controller.DeviceConfig{},
	}
}

func (c Config) Validate() error {

	return nil
}
