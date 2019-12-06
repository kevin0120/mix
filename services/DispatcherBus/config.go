package DispatcherBus

type Config struct {
}

func NewConfig() Config {
	return Config{}
}

func (c Config) Validate() error {
	return nil
}
