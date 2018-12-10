package minio

import (
	"github.com/masami10/aiis/toml"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Config struct {
	URL         string        `yaml:"url"`
	Bucket      string        `yaml:"bucket"`
	Access      string        `yaml:"access"`
	Secret      string        `yaml:"secret"`
	Secure      bool          `yaml:"secure"`
	ReuploadItv toml.Duration `yaml:"reupload_interval"`
	Enable      bool          `yaml:"enable"`
}

func NewConfig() Config {
	return Config{
		URL:         "127.0.0.1:9000",
		Bucket:      "bucket",
		Access:      "access",
		Secret:      "secret",
		Secure:      false,
		ReuploadItv: toml.Duration(time.Duration(24 * time.Hour)),
		Enable:      true,
	}
}

func (c Config) Validate() error {
	_, _, err := net.SplitHostPort(c.URL)
	if err != nil {
		return errors.Wrapf(err, "invalid minio  address %s", c.URL)
	}
	return nil
}
