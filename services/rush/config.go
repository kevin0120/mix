package rush

import (
	"github.com/masami10/aiis/toml"
	"time"
)

type Config struct {
	Enable       bool              `yaml:"enable"`
	Workers      int               `yaml:"workers"`
	Route        string            `yaml:"route"`
	Headers      map[string]string `yaml:"headers" override:"headers"`
	Timeout      toml.Duration     `yaml:"timeout"`
	PushInterval toml.Duration     `yaml:"push_interval"`
	MaxRetry     int               `yaml:"max_retry"` // api最大尝试次数

	WSRoute           string `yaml:"ws_route"`
	WSReadBufferSize  int    `yaml:"ws_read_buf_size"`
	WSWriteBufferSize int    `yaml:"ws_write_buf_size"`

	BatchSaveRowsLimit int           `yaml:"batch_save_rows_limit"`
	BatchSaveTimeLimit toml.Duration `yaml:"batch_save_time_limit"`
}

func NewConfig() Config {
	return Config{
		Enable:       true,
		Workers:      5,
		Route:        "/rush/v1/results", //方法永远为patch
		Timeout:      toml.Duration(time.Millisecond * 10),
		PushInterval: toml.Duration(time.Second * 1),
		MaxRetry:     3,
		Headers:      map[string]string{"Content-Type": "application/json"},

		BatchSaveRowsLimit: 64,
		BatchSaveTimeLimit: toml.Duration(30 * time.Second),
	}
}

func (c Config) Validate() error {
	return nil
}
