package changan

import (
	"github.com/masami10/aiis/toml"
	"time"
	"github.com/satori/go.uuid"
)

type Config struct {
	AndonAddr       string `yaml:"andon_addr"`
	Enable          bool   `yaml:"enable"`
	KeepAlivePeriod toml.Duration `yaml:"keepalive_period"`
	ReadTimeout     toml.Duration `yaml:"read_timeout"`
	ReadBufferSize  int `yaml:"read_buffer_size"`
	GUID		 	string `yaml:"GUID"`
}

func NewConfig() Config {

	_guid,_:= uuid.NewV4()


	return Config{
		AndonAddr:       "tcp://192.168.1.1:9002",
		Enable:          true,
		KeepAlivePeriod: toml.Duration(3 * time.Second),
		ReadTimeout:     toml.Duration(10 * time.Second),
		ReadBufferSize:  65535,
		GUID: _guid.String(),
	}
}

func (c Config) Validate() error {

	return nil
}
