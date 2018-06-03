package wsnotify

type Config struct {
	ReadBufferSize  int `yaml:"read_buf_size"`
	WriteBufferSize int `yaml:"write_buf_size"`
}

func NewConfig() Config {
	c := Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return c
}

func (c Config) Validate() error {
	return nil
}
