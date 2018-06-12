package audi_vw

import (
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	Port            int           `yaml:"port"`
	MaxConnections  int           `yaml:"max_connection"`
	ReadBufferSize  int           `yaml:"read_buf"`
	ReqTimeout      toml.Duration `yaml:"req_time_out"`
	KeepAlivePeriod toml.Duration `yaml:"keep_alive_period"`
}

func NewConfig() Config {
	return Config{
		Port:            4710,
		MaxConnections:  1024,
		ReadBufferSize:  65535,
		ReqTimeout:      toml.Duration(time.Microsecond * 3000),
		KeepAlivePeriod: toml.Duration(time.Second * 5), //协议层keepalive 周期
	}
}

func (c Config) Validate() error {
	return nil
}
