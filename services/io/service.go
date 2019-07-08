package io

import (
	"fmt"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	configValue atomic.Value
	ios         map[string]*IOModule
	diag        Diagnostic

	IONotify
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
		ios:  map[string]*IOModule{},
	}

	s.configValue.Store(c)

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {

	cfgs := s.config().IOS
	for _, v := range cfgs {
		s.ios[v.SN] = &IOModule{
			cfg: &v,
		}

		err := s.ios[v.SN].Start(s)
		if err != nil {
			s.diag.Error("start io failed", err)
		}
	}

	return nil
}

func (s *Service) Close() error {

	return nil
}

func (s *Service) Read(sn string) (string, string, error) {
	return s.ios[sn].Read()
}

func (s *Service) Write(sn string, index uint16, status uint16) error {
	return s.ios[sn].Write(index, status)
}

func (s *Service) OnStatus(sn string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s status:%s", sn, status))
}

func (s *Service) OnIOStatus(sn string, t string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s type:%s status:%s", sn, t, status))
}
