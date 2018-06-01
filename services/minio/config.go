package minio

import (
	"net"
	"github.com/pkg/errors"
)

type Config struct {
	URL string		`yaml:"url"`
	Bucket string	`yaml:"bucket"`
	Access string	`yaml:"access"`
	Secret string	`yaml:"secret"`
	Secure bool		`yaml:"secure"`
}


func NewConfig() Config {
	return Config{
		URL:     		"127.0.0.1:9000",
		Bucket:      	"bucket",
		Access: 		"access",
		Secret:     	"secret",
		Secure: 		false,
	}
}

func (c Config) Validate() error {
	_, _, err := net.SplitHostPort(c.URL)
	if err != nil {
		return errors.Wrapf(err, "invalid minio  address %s", c.URL)
	}
	return nil
}

