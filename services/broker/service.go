package broker

import (
	"github.com/masami10/rush/utils"
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
	diag                  Diagnostic
	configValue           atomic.Value
	Provider              IBrokerProvider
	opened                bool
	BrokerStatusDisptcher *utils.Dispatcher
	closing               chan struct{}
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag:                  d,
		opened:                false,
		BrokerStatusDisptcher: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		closing:               make(chan struct{}, 1),
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
		s.BrokerStatusDisptcher.Start()
		go s.connectProc()
	}
	return nil
}

func (s *Service) Close() error {
	s.closing <- struct{}{}
	s.BrokerStatusDisptcher.Release()
	if s.Provider != nil {
		return s.Provider.Close()
	}
	return nil
}

func (s *Service) connectProc() {
	for {
		select {
		case <-time.After(1 * time.Second):
			if err := s.Provider.Connect(s.Config().ConnectUrls); err != nil {
				continue
			} else {
				s.opened = true
				s.diag.Debug("Broker Service Started")
				s.BrokerStatusDisptcher.Dispatch(s.opened)
				return
			}

		case <-s.closing:
			return
		}
	}
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
		return errors.New("Can Not Create Broker SubscribeControllerInfo, Cause Provider Is Empty")
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
