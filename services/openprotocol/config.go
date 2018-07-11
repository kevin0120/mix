package openprotocol

import (
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	KeepAlivePeriod toml.Duration `yaml:"keep_alive_period"`
}

func NewConfig() Config {
	return Config{
		KeepAlivePeriod: toml.Duration(time.Second * 10), //keep alive周期
	}
}

func (c Config) Validate() error {
	return nil
}
