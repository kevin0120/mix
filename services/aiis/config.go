package aiis

import (
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	Enable              bool          `yaml:"enable"`
	Timeout             toml.Duration `yaml:"timeout"`
	WSResultRoute       string        `yaml:"ws_result_route"`
	ResultUploadInteval toml.Duration `yaml:"result_upload_inteval"`
	KeepAlive           toml.Duration `yaml:"keep_alive"`
	Recheck             bool          `yaml:"recheck"`
	TransportType       string        `yaml:"transport_type"`
}

func NewConfig() Config {
	c := Config{
		Enable:              true,
		Timeout:             toml.Duration(5 * time.Second),
		WSResultRoute:       "ws://127.0.0.1/aiis/v1/ws/results",
		ResultUploadInteval: toml.Duration(time.Duration(1 * time.Hour)),
		KeepAlive:           toml.Duration(time.Second * 3),
		Recheck:             true,
		TransportType:       TRANSPORT_TYPE_GRPC,
	}

	return c
}

func (c Config) Validate() error {

	return nil
}
