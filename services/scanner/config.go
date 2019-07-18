package scanner

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
)

type DeviceConfig struct {
	VendorID  int `yaml:"vendorID"`
	ProductID int `yaml:"productID"`
}

type Config struct {
	Enable      bool   `yaml:"enable"`
	EntityLabel string `yaml:"entity_label"` // windows: COM4, linux: VID(int):PID(int) like 1234:4567
}

func NewConfig() Config {
	var label string
	if runtime.GOOS == "windows" {
		label = "COM4"
	} else {
		label = fmt.Sprintf("%d:%d", VendorHoneyWell, ProductHoneyWell)
	}
	return Config{
		Enable:      true,
		EntityLabel: label,
	}
}

func (c Config) Validate() error {
	if c.Enable {
		if c.EntityLabel == "" {
			return errors.New("EntityLabel Is Empty")
		}
	}
	return nil
}
