package broker

import (
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value
	Provider    IBrokerProvider
	opened      bool
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag:   d,
		opened: false,
	}
	s.configValue.Store(c)

	p := s.newBroker(c.Provider)
	s.Provider = p

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	c := s.Config()
	if c.Enable {
		if err := s.Provider.Connect(c.ConnectUrls); err != nil {
			s.Provider = NewDefaultBroker()
			return err
		} else {
			s.opened = true
		}
	}
	return nil
}

func (s *Service) Close() error {
	if s.Provider != nil {
		return s.Provider.Close()
	}
	return nil
}

func (s *Service) newBroker(provider string) (ret IBrokerProvider) {
	c := s.Config()
	switch provider {
	case "nats":
		ret = NewNats(s.diag, c)
	default:
		ret = NewDefaultBroker()
	}
	return
}

func (s *Service) Subscribe(subject string, handler SubscribeHandler) error {
	p := s.Provider
	if p == nil {
		return errors.New("Can Not Create Broker Subscribe, Cause Provider Is Empty")
	}

	return p.Subscribe(subject, handler)
}

func (s *Service) Publish(subject string, data []byte) error {
	p := s.Provider
	if p == nil {
		return errors.New("Can Not Create Broker Publish, Cause Provider Is Empty")
	}

	return p.Publish(subject, data)
}

func (s *Service) Request(subject string, data []byte, timeOut time.Duration) ([]byte, error) {
	p := s.Provider
	if p == nil {
		return nil, errors.New("Can Not Create Broker Publish, Cause Provider Is Empty")
	}

	return p.DoRequest(subject, data, timeOut)
}
