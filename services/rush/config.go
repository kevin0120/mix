package rush

type Config struct {
	Enable 		bool 	`yaml:"enable"`
	Workers		int		`yaml:"workers"`
}

func NewConfig() Config {
	return Config{
		Enable: true,
		Workers: 5,
	}
}

func (c Config) Validate() error {
	return nil
}
