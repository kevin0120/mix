package ts002

import (
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	Enable bool `yaml:"enable"`

	// 柜锁控制（当收刷卡登录成功时，控制相关IO输出）
	IOLocker []int `yaml:"io_locker"`

	// 报警控制（当收到MES报警请求或拧紧NOK时，控制相关IO输出, 持续5秒）
	IOAlarm []int `yaml:"io_alarm"`

	// 报警持续时间
	IOAlarmLast toml.Duration `yaml:"io_alarm_last"`

	// Mes接口配置
	MesApiConfig MesApiConfig `yaml:"mes_api"`
}

type MesApiConfig struct {
	APIUrl               string        `yaml:"api_url"`
	Timeout              toml.Duration `yaml:"timeout"`
	RetryCount           int           `yaml:"retry_count"`
	EndpointCardInfo     string        `yaml:"endpoint_cardinfo"`
	EndpointResultUpload string        `yaml:"endpoint_result_upload"`
}

func NewMesApiConfig() MesApiConfig {
	return MesApiConfig{
		APIUrl:  "http://127.0.0:1:8000",
		Timeout: toml.Duration(time.Duration(10 * time.Second)),
		RetryCount: 5,
		EndpointCardInfo: "/api/v1/cardinfo",
		EndpointResultUpload: "/api/v1/results",
	}
}

func NewConfig() Config {
	return Config{
		Enable:       true,
		IOLocker:     []int{1},
		IOAlarm:      []int{1},
		MesApiConfig: NewMesApiConfig(),
	}
}

func (c Config) Validate() error {
	return nil
}
