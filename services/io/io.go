package io

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/device"
)

const (
	IO_STATUS_ONLINE  = "online"
	IO_STATUS_OFFLINE = "offline"

	IO_TYPE_INPUT  = "input"
	IO_TYPE_OUTPUT = "output"

	IO_MODBUSTCP = "modbustcp"
)

type IONotify interface {
	OnStatus(sn string, status string)
	OnIOStatus(sn string, t string, status string)
}

type IO interface {
	Status() string
	Start() error
	Stop() error
	Read() (string, string, error)
	Write(uint16, uint16) error
}

type IOModule struct {
	cfg    ConfigIO
	client IO
}

func (s *IOModule) Start(srv *Service) error {

	vendor := VENDOR_MODELS[s.cfg.Model]
	switch vendor.Type() {
	case IO_MODBUSTCP:
		s.client = &ModbusTcp{
			cfg:    s.cfg,
			notify: srv,
			vendor: vendor,
		}

	default:
		return errors.New(fmt.Sprintf("invalid model type: %s", s.cfg.Model))
	}

	return s.client.Start()
}

func (s *IOModule) Status() string {
	return s.client.Status()
}

func (s *IOModule) Read() (string, string, error) {
	return s.client.Read()
}

func (s *IOModule) Write(index uint16, status uint16) error {
	return s.client.Write(index, status)
}

func (s *IOModule) DeviceType(sn string) string {
	return "io"
}

func (s *IOModule) Children() map[string]device.IDevice {
	return map[string]device.IDevice{}
}

func (s *IOModule) Data() interface{} {
	inputs, outputs, err := s.Read()
	if err != nil {
		return nil
	}

	return IoData{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (s *IOModule) Config() interface{} {
	vendor := VENDOR_MODELS[s.cfg.Model]

	return IoConfig{
		InputNum:  vendor.Cfg().InputNum,
		OutputNum: vendor.Cfg().OutputNum,
	}
}
