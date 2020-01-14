package transport

import "github.com/pkg/errors"

const (
	GRPCTransport         = "grpc"
	BrokerTransport       = "broker"
	HttpTransport         = "http"
)

var SupportProviders = []string{GRPCTransport, BrokerTransport, HttpTransport}

type Config struct {
	Provider string   `yaml:"provider"`
}

func NewConfig() Config {
	return Config{
		Provider: GRPCTransport,
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
