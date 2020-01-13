package grpc

import (
	"github.com/masami10/rush/services/transport"
	microTransport "github.com/micro/go-micro/transport"
)

type Service struct {
	microTransport.Transport
	addr string
}

func newGRPCTransport(config Config) transport.Transport {
	options := microTransport.Addrs(config.addr)
	return microTransport.NewTransport(options)
}

func NewService(config Config) *Service {
	s := &Service{addr: config.addr}
	s.Transport = newGRPCTransport(config)
	return s
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Connect(addr string) (transport.Client, error) {
	c, err := s.Dial(addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (s *Service)GetServerAddress() []string  {
	return s.Transport.Options().Addrs
}
