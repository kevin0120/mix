package transport

import (
	"github.com/pkg/errors"
)

type Service struct {
	ITransportService
	diag     Diagnostic
	provider string
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag:     d,
		provider: c.Provider,
	}
	return s
}

func (s *Service) BindTransportByProvider(ss IServer) error {
	switch s.provider {
	case BrokerTransport, GRPCTransport:
		p := ss.GetServiceByName(s.provider)
		if p == nil {
			err := errors.Errorf("%s Service Is Not Found", s.provider)
			return err
		}
		s.ITransportService = p.(ITransportService)
	default:
		err := errors.Errorf("Provider: %s Is Not Support Right Now", s.provider)
		s.diag.Error("BindTransportByProvider", err)
	}
	return nil
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) SetStatusHandler(handler StatusHandler) error {
	t := s.ITransportService
	if t == nil {
		return errors.New("Please BindTransportByProvider First")
	}
	return t.SetStatusHandler(handler)
}
