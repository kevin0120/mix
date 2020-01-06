package broker

import (
	"fmt"
	"github.com/masami10/rush/services/dispatcherbus"
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
	diag                Diagnostic
	configValue         atomic.Value
	Provider            IBrokerProvider
	opened              bool
	DispatcherBus       Dispatcher
	BrokerDispatcherMap dispatcherbus.DispatcherMap
	//BrokerStatusDisptcher *utils.Dispatcher
	closing chan struct{}
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag:    d,
		closing: make(chan struct{}, 1),
	}
	s.configValue.Store(c)

	p := s.newBroker(c.Provider)
	s.Provider = p

	s.initGblDispatcher()

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initGblDispatcher() {
	s.BrokerDispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DISPATCHER_BROKER_STATUS: utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) Open() error {
	c := s.Config()
	if !c.Enable {
		return nil
	}

	s.DispatcherBus.LaunchDispatchersByHandlerMap(s.BrokerDispatcherMap)
	s.doConnect(false) // 初始化所有连接状态为未连接
	go s.connectProc()
	return nil
}

func (s *Service) Close() error {
	s.closing <- struct{}{}
	s.DispatcherBus.ReleaseDispatchersByHandlerMap(s.BrokerDispatcherMap)
	if s.Provider != nil {
		return s.Provider.Close()
	}
	return nil
}

func (s *Service) doConnect(opened bool) {
	s.opened = true
	s.diag.Debug(fmt.Sprintf("Broker Service Is Opened: %v", opened))
	status := utils.STATUS_OFFLINE
	if opened {
		status = utils.STATUS_ONLINE
	}
	s.DispatcherBus.Dispatch(dispatcherbus.DISPATCHER_BROKER_STATUS, status)
}

func (s *Service) connectProc() {
	for {
		select {
		case <-time.After(1 * time.Second):
			if err := s.Provider.Connect(s.Config().ConnectUrls); err != nil {
				continue
			} else {
				s.doConnect(true)
				return
			}

		case <-s.closing:
			s.doConnect(false)
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
