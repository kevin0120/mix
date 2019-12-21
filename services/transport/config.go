package transport

import "github.com/pkg/errors"

const (
	GRPCTransport         = "grpc"
	BrokerTransport       = "broker"
	DefaultGRPCServerAddr = "127.0.0.1:9093"
)

var SupportProviders = []string{GRPCTransport, BrokerTransport}

type Config struct {
	Provider string   `yaml:"provider"`
	Address  []string `yaml:"transport_address"`
}

func NewConfig() Config {
	return Config{
		Provider: GRPCTransport,
		Address:  []string{DefaultGRPCServerAddr},
	}
}

func isSupportProvider(provider string) bool {
	for _, p := range SupportProviders {
		if p == provider {
			return true
		}
	}
	return false
}

func (c Config) Validate() error {
	if !isSupportProvider(c.Provider) {
		return errors.Errorf("Provider: %s Is Not Support", c.Provider)
	}
	return nil
}
