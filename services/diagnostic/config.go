package diagnostic

import "gopkg.in/yaml.v2"

type Config struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

func NewConfig() Config {
	return Config{
		File:  "STDERR",
		Level: "DEBUG",
	}
}
