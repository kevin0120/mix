package masterplc

type Config struct {
	Port 			int				`yaml:"port"`
	Enable          bool          `yaml:"enable"`
	MaxConnections  int           `yaml:"max_connection"`
	ReadBufferSize  int           `yaml:"read_buf"`
	Workers int					`yaml:"workers"`
}

func NewConfig() Config {

	return Config{
		Port: 			4720,
		Enable:          false,
		MaxConnections:  4096,
		ReadBufferSize:  65535,
		Workers: 4,
	}
}

func (c Config) Validate() error {

	return nil
}