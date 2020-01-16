package broker

import (
	"fmt"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/transport"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"time"
)

type Service struct {
	diag          Diagnostic
	configValue   atomic.Value
	provider      IBrokerProvider
	dispatcherBus Dispatcher
	dispatcherMap dispatcherbus.DispatcherMap
	opened        bool
	closing       chan struct{}
	status        atomic.String
}

func NewService(c Config, d Diagnostic, dp Dispatcher) *Service {

	s := &Service{
		diag:          d,
		closing:       make(chan struct{}, 1),
		dispatcherBus: dp,
	}
	s.configValue.Store(c)

	p := s.newBroker(c.Provider)
	s.provider = p

	s.setupGblDispatcher()

	s.status.Store(utils.STATUS_OFFLINE)

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) setupGblDispatcher() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DispatcherBrokerStatus: utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) doOpen() {
	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)
	s.doConnect(false) // 初始化所有连接状态为未连接
	go s.connectProc()
	s.opened = true
}

func (s *Service) Open() error {
	c := s.Config()
	if !c.Enable {
		return nil
	}
	s.doOpen()
	return nil
}

func (s *Service) TransportForceOpen() error {
	if s.opened {
		return nil
	}
	c := s.Config()
	if err := c.Validate(); err != nil {
		s.diag.Error("TransportForceOpen", err)
		return err
	}
	s.doOpen()
	return nil
}

func (s *Service) Close() error {
	if !s.opened {
		return nil
	}
	s.closing <- struct{}{}
	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	if s.provider != nil {
		return s.provider.Close()
	}
	return nil
}

func (s *Service) dispatcherBrokerStatus(status string) {
	if s.dispatcherBus == nil {
		s.diag.Error("dispatcherBrokerStatus Error", errors.New("dispatcherBus Is Empty"))
	}
	if err := s.dispatcherBus.Dispatch(dispatcherbus.DispatcherBrokerStatus, status); err != nil {
		s.diag.Error("dispatcherBrokerStatus", err)
	}

}

func (s *Service) doConnect(opened bool) {
	s.diag.Debug(fmt.Sprintf("broker Service Is Opened: %v", opened))
	status := utils.STATUS_OFFLINE
	if opened {
		status = utils.STATUS_ONLINE
	}
	s.dispatcherBrokerStatus(status)
}

func (s *Service) connectProc() {
	for {
		select {
		case <-time.After(1 * time.Second):
			if err := s.provider.Connect(s.Config().ConnectUrls); err != nil {
				continue
			} else {
				s.dispatcherBrokerStatus(utils.STATUS_ONLINE)
				s.diag.Debug(fmt.Sprintf("broker Service Is Opened"))
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

func (s *Service) OnMessage(subject string, handler transport.OnMsgHandler) error {
	return s.subscribe(subject, handler)
}

func (s *Service) SendMessage(subject string, data []byte) error {
	return s.publish(subject, data)
}

func (s *Service) GetServerAddress() []string {
	c := s.Config()
	return c.ConnectUrls
}

func (s *Service) subscribe(subject string, handler transport.OnMsgHandler) error {
	p := s.provider
	if p == nil {
		return errors.New("Can Not Create broker Subscribe, Cause provider Is Empty")
	}

	return p.Subscribe(subject, handler)
}

func (s *Service) publish(subject string, data []byte) error {
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
	s.dispatcherBrokerStatus(status)
}

func (s *Service) Status() string {
	return s.status.Load()
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error("Dispatch Failed", err)
	}
}
