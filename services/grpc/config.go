package grpc

const (
	DefaultGRPCServerAddr = "127.0.0.1:9093"
)

type Config struct {
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address"`
	Workers int    `yaml:"workers"`
}

func NewConfig() Config {
	return Config{
		Enable:  true,
		Address: DefaultGRPCServerAddr,
		Workers: 2,
	}
}

func (c Config) Validate() error {
	return nil
}
