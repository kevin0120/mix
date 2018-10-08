package openprotocol

import (
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	KeepAlivePeriod toml.Duration `yaml:"keep_alive_period"`
	ReadBufferSize  int           `yaml:"read_buf"`
	ReqTimeout      toml.Duration `yaml:"req_time_out"`
	SkipJob         int           `yaml:"skip_job"`
	IOTrigger       int           `yaml:"io_trigger"`
	DataIndex       int           `yaml:"data_index"`
}

func NewConfig() Config {
	return Config{
		KeepAlivePeriod: toml.Duration(time.Second * 10), //keep alive周期
		ReadBufferSize:  65535,
		ReqTimeout:      toml.Duration(time.Microsecond * 3000),
		SkipJob:         0,
		IOTrigger:       0,
		DataIndex:       1,
	}
}

func (c Config) Validate() error {
	return nil
}
