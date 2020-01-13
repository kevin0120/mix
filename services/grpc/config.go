package grpc

type Config struct {
	addr string        `yaml:"address"`
}

func NewConfig() Config {
	return Config{
		addr: "localhost:8082",
	}
}
