package openprotocol

import (
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	KeepAlivePeriod   toml.Duration `yaml:"keep_alive_period"`
	ReadBufferSize    int           `yaml:"read_buf"`
	ReqTimeout        toml.Duration `yaml:"req_time_out"`
	SkipJobs          []int         `yaml:"skip_job"`
	IOTrigger         int           `yaml:"io_trigger"`
	DataIndex         int           `yaml:"data_index"`
	VinIndex          []int         `yaml:"vin_index"`
	GetToolInfoPeriod toml.Duration `yaml:"tool_info_period"`
}

func NewConfig() Config {

	return Config{
		KeepAlivePeriod:   toml.Duration(time.Second * 10), //keep alive周期
		ReadBufferSize:    65535,
		ReqTimeout:        toml.Duration(time.Microsecond * 3000),
		SkipJobs:          []int{250},
		IOTrigger:         0,
		DataIndex:         1,
		VinIndex:          []int{0, 1},
		GetToolInfoPeriod: toml.Duration(time.Hour * 12), // 半天
	}
}

func (c Config) Validate() error {
	return nil
}
