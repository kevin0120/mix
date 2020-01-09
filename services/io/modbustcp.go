package io

import (
	"errors"
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/masami10/rush/utils"
	"go.uber.org/atomic"
	"strconv"
	"time"
)

const (
	TIMEOUT = 3 * time.Second
	ReadItv = 300 * time.Millisecond
)

type ModbusTcp struct {
	cfg    ConfigIO
	status atomic.Value

	inputs  atomic.Value
	outputs atomic.Value

	handler *modbus.TCPClientHandler
	client  modbus.Client

	vendor Vendor
	notify IONotify
}

func (s *ModbusTcp) Start() error {
	s.status.Store(IoStatusOffline)
	s.inputs.Store("")
	s.outputs.Store("")

	go s.connect()

	return nil
}

func (s *ModbusTcp) Stop() error {
	if s.handler != nil {
		return s.handler.Close()
	}

	return nil
}

func (s *ModbusTcp) Status() string {
	return s.status.Load().(string)
}

func (s *ModbusTcp) createHandler() {
	if s.handler != nil {
		_ = s.handler.Close()
	}

	s.handler = modbus.NewTCPClientHandler(s.cfg.Address)
	s.handler.Timeout = TIMEOUT
}

func (s *ModbusTcp) connect() {
	s.createHandler()

	for {
		err := s.handler.Connect()

		if err != nil {
			continue
		} else {
			// online
			s.status.Store(IoStatusOnline)
			s.notify.OnStatus(s.cfg.SN, IoStatusOnline)
			s.client = modbus.NewClient(s.handler)
			s.read()
		}
	}
}

func (s *ModbusTcp) read() {
	for {
		_, _, err := s.Read()
		if err != nil {
			// offline
			s.status.Store(IoStatusOffline)
			s.notify.OnStatus(s.cfg.SN, IoStatusOffline)
			s.createHandler()
			return
		}

		time.Sleep(ReadItv)
	}
}

func (s *ModbusTcp) formatIO(results []byte, num uint16) string {
	rt := ""
	resultLen := int(num / 8)
	for i := 0; i < resultLen; i++ {
		strIO := strconv.FormatUint(uint64(results[i]), 2)
		format := fmt.Sprintf("%%0%ds", num)
		strIO = fmt.Sprintf(format, strIO)
		rt += utils.ReverseString(strIO)
	}

	return rt
}

func (s *ModbusTcp) Read() (string, string, error) {
	client := s.client
	var err error
	var result []byte
	result = []byte{}

	inputs := ""
	outputs := ""

	if client == nil {
		return inputs, outputs, errors.New("client is nil")
	}

	// input status
	switch s.vendor.Cfg().InputReadType {
	case ReadTypeDiscretes:
		result, err = client.ReadDiscreteInputs(s.vendor.Cfg().InputAddress, s.vendor.Cfg().InputNum)
	}

	if err != nil {
		return inputs, outputs, err
	}

	inputs = s.formatIO(result, s.vendor.Cfg().InputNum)

	// output status
	switch s.vendor.Cfg().OutputReadType {
	case ReadTypeCoils:
		result, err = client.ReadCoils(s.vendor.Cfg().OutputAddress, s.vendor.Cfg().OutputNum)
	}

	if err != nil {
		return inputs, outputs, err
	}

	outputs = s.formatIO(result, s.vendor.Cfg().OutputNum)

	if s.inputs.Load().(string) != inputs {
		s.inputs.Store(inputs)
		s.notify.OnChangeIOStatus(s.cfg.SN, IoTypeInput, inputs)
	}

	if s.outputs.Load().(string) != outputs {
		s.outputs.Store(outputs)
		s.notify.OnChangeIOStatus(s.cfg.SN, IoTypeOutput, outputs)
	}

	return inputs, outputs, nil
}

func (s *ModbusTcp) Write(index uint16, status uint16) error {
	if s.Status() == IoStatusOffline {
		return errors.New(IoStatusOffline)
	}

	if index > (s.vendor.Cfg().OutputNum - 1) {
		return errors.New("invalid index")
	}

	var err error
	client := s.client

	if status == 1 {
		status = 0xFF00
	}

	switch s.vendor.Cfg().WriteType {
	case WriteTypeSingleCoil:
		_, err = client.WriteSingleCoil(index, status)
	}

	return err
}
