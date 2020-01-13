package transport

import (
	"github.com/pkg/errors"
)

type Service struct {
	ITransportService
	diag     Diagnostic
	provider string
	addrs    []string
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag:     d,
		provider: c.Provider,
		addrs:    c.Address,
	}
	return s
}

func (s *Service) BindTransportByProvider(ss IServer) error {
	switch s.provider {
	case BrokerTransport:
		p := ss.GetServiceByName(BrokerTransport)
		if p == nil {
			err := errors.Errorf("%s Service Is Not Found", BrokerTransport)
			return err
		}
		s.ITransportService = p.(ITransportService)
	default:
		err := errors.Errorf("Provider: %s Is Not Support Right Now", s.provider)
		s.diag.Error("BindTransportByProvider", err)
	}
	return nil
}
