package transport

import "github.com/pkg/errors"

var SupportProviders = []string{GRPCTransport, BrokerTransport, HttpTransport}

type Config struct {
	Provider string `yaml:"provider"`
}

func NewConfig() Config {
	return Config{
		Provider: BrokerTransport,
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
