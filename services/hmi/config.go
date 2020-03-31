package hmi

type Config struct {
	Enable bool `yaml:"enable"`
}

func NewConfig() Config {
	return Config{
		Enable: false,
	}
}

func (c Config) Validate() error {
	return nil
}
