package openprotocol

import (
	"github.com/masami10/rush/services/controller"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
}

type Service struct {
	configValue atomic.Value
	name        string
	Controllers map[string]*Controller
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		Controllers: map[string]*Controller{},
		name:        controller.OPENPROTOCOL,
	}

	return s
}

func (p *Service) Parse(msg string) ([]byte, error) {
	return nil, nil
}

func (p *Service) Write(sn string, buf []byte) error {
	return nil
}

func (p *Service) AddNewController(cfg controller.Config) {

}

func (p *Service) Open() error {
	return nil
}

func (p *Service) Close() error {
	return nil
}
