package pmon

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	PmonDir string `yaml:"pmon-dir"`
	Enable  bool   `yaml:"enable"`
}

func NewConfig() Config {
	var homeDir string
	// By default, store meta and data files in current users home directory
	u, err := user.Current()
	if err == nil {
		homeDir = u.HomeDir
	} else if os.Getenv("HOME") != "" {
		homeDir = os.Getenv("HOME")
	} else {
		homeDir = "/etc"
	}

	pmondir := filepath.Join(homeDir, "pmon")
	return Config{
		PmonDir: pmondir,
		Enable:  false,
	}
}

func (c Config) Validate() error {
	if !c.Enable {
		return nil
	}
	if c.PmonDir == "" {
		return fmt.Errorf("must specify PMON lib path")
	}
	if filepath.IsAbs(c.PmonDir) {
		return errors.New("dir must be an related path")
	}
	return nil
}
