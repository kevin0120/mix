package scanner

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"strings"
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
		label = "COM5"
	} else {
		label = fmt.Sprintf("%d:%d", VendorDataLogic, ProductDataLogic)
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
		if runtime.GOOS == "windows" {
			if !strings.HasPrefix(c.EntityLabel, "COM") {
				return errors.New("Platform Windows EntityLabel Is COM Port")
			} else {
				if !strings.Contains(c.EntityLabel, ":") {
					return errors.New("Platform Unix/Linux EntityLabel Is VID:PID Format")
				}
			}
		}
	}
	return nil
}
