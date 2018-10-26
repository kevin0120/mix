package storage

import (
	"github.com/masami10/aiis/toml"
	"time"
)

type Config struct {
	Url                string        `yaml:"db_url"`
	DBName             string        `yaml:"db_name"`
	User               string        `yaml:"db_user"`
	Password           string        `yaml:"db_pwd"`
	MaxConnects        int           `yaml:"max_connection"`
	BatchSaveRowsLimit int           `yaml:"batch_save_rows_limit"`
	BatchSaveTimeLimit toml.Duration `yaml:"batch_save_time_limit"`
}

func NewConfig() Config {
	return Config{
		Url:                "127.0.0.1:5432",
		DBName:             "dbname",
		User:               "user",
		Password:           "pwd",
		MaxConnects:        60,
		BatchSaveRowsLimit: 64,
		BatchSaveTimeLimit: toml.Duration(30 * time.Second),
	}
}

func (c Config) Validate() error {

	return nil
}
