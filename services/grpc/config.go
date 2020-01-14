package grpc

type Config struct {
	addr    string `yaml:"address"`
	workers int    `yaml:"workers"`
}

func NewConfig() Config {
	return Config{
		addr:    "localhost:8082",
		workers: 2,
	}
}

func (c Config) Validate() error {
	return nil
}
