package odoo

import "fmt"

const (
	DEFAULT_URL = "127.0.0.1:8069"
)


type Config struct {
		URL 	string 		`yaml:"url"`
		Enable 	bool		`yaml:"enable"`
}

func NewConfig()  Config{
	return Config{
		URL:  DEFAULT_URL,
		Enable: true,

	}
}

func (c Config) Validate() error {
	if c.Enable && c.URL == "" {
		return fmt.Errorf("Odoo URL is empty ")
	}
	return nil
}
