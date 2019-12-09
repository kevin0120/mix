package io

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/device"
	"sync"
	"time"
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

	flashInterval time.Duration
	closing       chan struct{}
	flashes       map[uint16]uint16
	mtx           sync.Mutex
	opened        bool
	diag          Diagnostic
}

func (s *IOModule) Start(srv *Service) error {

	s.flashes = map[uint16]uint16{}
	s.closing = make(chan struct{}, 1)
	s.mtx = sync.Mutex{}

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

	go s.flashProc()
	s.opened = true

	return s.client.Start()
}

func (s *IOModule) Stop() error {
	if s.opened {
		s.closing <- struct{}{}

		return s.client.Stop()
	}
	return nil
}

func (s *IOModule) Status() string {
	return s.client.Status()
}

func (s *IOModule) Read() (string, string, error) {
	return s.client.Read()
}

func (s *IOModule) Write(index uint16, status uint16) error {
	switch status {
	case OUTPUT_STATUS_OFF:
		// 从flash列表删除
		s.removeFlash(index)
	case OUTPUT_STATUS_FLASH:
		// 加入flash列表
		s.addFlash(index)
		return nil
	}

	return s.client.Write(index, status)
}

func (s *IOModule) DeviceType() string {
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

func (s *IOModule) addFlash(output uint16) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.flashes[output] = output
}

func (s *IOModule) removeFlash(output uint16) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.flashes, output)
}

func (s *IOModule) getFlashes() map[uint16]uint16 {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.flashes
}

func (s *IOModule) flashProc() {

	status := OUTPUT_STATUS_OFF
	flag := -1
	for {
		select {
		case <-time.After(s.flashInterval):
			// 状态0<->1变换
			flag *= -1
			status += flag

			flashes := s.getFlashes()
			for _, v := range flashes {
				err := s.Write(v, uint16(status))
				if err != nil {
					s.diag.Error("Write Failed", err)
				}
			}

		case <-s.closing:
			return
		}
	}
}
