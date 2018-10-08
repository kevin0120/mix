package changan

import (
	"github.com/masami10/aiis/toml"
	"github.com/satori/go.uuid"
	"time"
)

type ConfigAndonDB struct {
	Url         string `yaml:"db_url"`
	DBName      string `yaml:"db_name"`
	User        string `yaml:"db_user"`
	Password    string `yaml:"db_pwd"`
	MaxConnects int    `yaml:"max_connection"`
}

type Config struct {
	AndonAddr       string        `yaml:"andon_addr"`
	Enable          bool          `yaml:"enable"`
	KeepAlivePeriod toml.Duration `yaml:"keepalive_period"`
	ReadTimeout     toml.Duration `yaml:"read_timeout"`
	ReadBufferSize  int           `yaml:"read_buffer_size"`
	GUID            string        `yaml:"GUID"`
	DB              ConfigAndonDB `yaml:"db"`
}

func NewConfig() Config {

	_guid := uuid.NewV4()

	return Config{
		AndonAddr:       "tcp://192.168.1.4:9002",
		Enable:          true,
		KeepAlivePeriod: toml.Duration(3 * time.Second),
		ReadTimeout:     toml.Duration(10 * time.Second),
		ReadBufferSize:  65535,
		GUID:            _guid.String(),
		DB: ConfigAndonDB{
			Url:         "127.0.0.1:1433",
			DBName:      "cvinet",
			User:        "sa",
			Password:    "sa",
			MaxConnects: 64,
		},
	}
}

func (c Config) Validate() error {

	return nil
}
