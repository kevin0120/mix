package io

import (
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/toml"
	"time"
)

type ConfigIO struct {
	SN      string `yaml:"sn"`
	Model   string `yaml:"model"`
	Address string `yaml:"address"`
}

type Config struct {
	Enable       bool          `yaml:"enable"`
	FlashInteval toml.Duration `yaml:"flash_inteval"`
	IOS          []ConfigIO    `yaml:"ios"`
}

func NewConfig() Config {
	return Config{
		Enable:       true,
		FlashInteval: toml.Duration(time.Second * 1),
		IOS: []ConfigIO{
			{
				SN:      "1",
				Model:   ModelMoxaE1212,
				Address: "modbustcp://127.0.0.1:502",
			},
		},
	}
}

func (c Config) Validate() error {

	for _, io := range c.IOS {
		_, exist := VendorModels[io.Model]
		if !exist {
			return errors.New("Vendor Not Found")
		}
	}

	return nil
}
