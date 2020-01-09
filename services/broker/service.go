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
	diag          Diagnostic
	configValue   atomic.Value
	provider      IBrokerProvider
	opened        bool
	dispatcherBus Dispatcher
	dispatcherMap dispatcherbus.DispatcherMap
	closing       chan struct{}
	status        atomic.Value
}

func NewService(c Config, d Diagnostic, dp Dispatcher) *Service {

	s := &Service{
		diag:          d,
		closing:       make(chan struct{}, 1),
		dispatcherBus: dp,
	}
	s.configValue.Store(c)
	s.status.Store(utils.STATUS_OFFLINE)

	p := s.newBroker(c.Provider)
	s.provider = p

	s.setupGblDispatcher()

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) setupGblDispatcher() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DISPATCHER_BROKER_STATUS: utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) Open() error {
	c := s.Config()
	if !c.Enable {
		return nil
	}

	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)
	s.doConnect(false) // 初始化所有连接状态为未连接
	go s.connectProc()
	return nil
}

func (s *Service) Close() error {
	s.closing <- struct{}{}
	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	if s.provider != nil {
		return s.provider.Close()
	}
	return nil
}

func (s *Service) doConnect(opened bool) {
	s.opened = true
	s.diag.Debug(fmt.Sprintf("broker Service Is Opened: %v", opened))
	status := utils.STATUS_OFFLINE
	if opened {
		status = utils.STATUS_ONLINE
	}
	s.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_BROKER_STATUS, status)
}

func (s *Service) connectProc() {
	for {
		select {
		case <-time.After(1 * time.Second):
			if err := s.provider.Connect(s.Config().ConnectUrls); err != nil {
				continue
			} else {
				s.opened = true
				s.diag.Debug(fmt.Sprintf("broker Service Is Opened: %v", s.opened))
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

	ret.SetStatusHandler(s.onStatus)
	return
}

func (s *Service) Subscribe(subject string, handler SubscribeHandler) error {
	p := s.provider
	if p == nil {
		return errors.New("Can Not Create broker Subscribe, Cause provider Is Empty")
	}

	return p.Subscribe(subject, handler)
}

func (s *Service) Publish(subject string, data []byte) error {
	p := s.provider
	if p == nil {
		return errors.New("Can Not Create broker Publish, Cause provider Is Empty")
	}

	return p.Publish(subject, data)
}

func (s *Service) Request(subject string, data []byte, timeOut time.Duration) ([]byte, error) {
	p := s.provider
	if p == nil {
		return nil, errors.New("Can Not Create broker Publish, Cause provider Is Empty")
	}

	return p.DoRequest(subject, data, timeOut)
}

func (s *Service) onStatus(status string) {
	s.status.Store(status)
	s.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_BROKER_STATUS, status)
}

func (s *Service) Status() string {
	return s.status.Load().(string)
}
