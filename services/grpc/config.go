package grpc

const (
	DefaultGRPCServerAddr = "127.0.0.1:9093"
)

type Config struct {
	Address string `yaml:"address"`
	Workers int    `yaml:"workers"`
}

func NewConfig() Config {
	return Config{
		Address: DefaultGRPCServerAddr,
		Workers: 2,
	}
}

func (c Config) Validate() error {
	return nil
}
