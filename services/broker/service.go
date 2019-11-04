package broker

import (
	"github.com/pkg/errors"
	"sync/atomic"
)

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type IBrokerProvider interface {
	Connect(urls []string) error
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	Provider    IBrokerProvider
}

func NewService(c Config, d Diagnostic) *Service {

	srv := &Service{
		diag: d,
	}

	srv.configValue.Store(c)

	return srv
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	if c.Enable {

	}
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) NewBroker(provider string) (ret IBrokerProvider, err error) {
	c := s.Config()
	switch c.Provider {
	case "nats":
		ret = NewNats(s.diag)
	default:
		err = errors.Errorf("Provider: %d Is Not Support", provider)
		s.diag.Error("New Broker", err)

	}
	return
}
