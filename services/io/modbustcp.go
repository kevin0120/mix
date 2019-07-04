package io

import (
	"errors"
	"fmt"
	"github.com/goburrow/modbus"
	"github.com/masami10/rush/utils"
	"strconv"
	"sync/atomic"
	"time"
)

const ()

const (
	TIMEOUT  = 3 * time.Second
	READ_ITV = 500 * time.Millisecond
)

type ModbusTcp struct {
	cfg    *ConfigIO
	status atomic.Value

	inputs  atomic.Value
	outputs atomic.Value

	handler *modbus.TCPClientHandler
	client  modbus.Client

	vendor Vendor
	notify IONotify
}

func (s *ModbusTcp) Start() error {
	s.status.Store(IO_STATUS_OFFLINE)
	s.inputs.Store("")
	s.outputs.Store("")

	go s.connect()

	return nil
}

func (s *ModbusTcp) Stop() error {
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
			s.status.Store(IO_STATUS_ONLINE)
			s.notify.OnStatus(s.cfg.SN, IO_STATUS_ONLINE)
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
			s.status.Store(IO_STATUS_OFFLINE)
			s.notify.OnStatus(s.cfg.SN, IO_STATUS_OFFLINE)
			s.createHandler()
			return
		}

		time.Sleep(READ_ITV)
	}
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
	case READ_TYPE_DISCRETES:
		result, err = client.ReadDiscreteInputs(s.vendor.Cfg().InputAddress, s.vendor.Cfg().InputNum)
	}

	if err != nil {
		return inputs, outputs, err
	}

	inputs = strconv.FormatUint(uint64(result[0]), 2)

	// output status
	switch s.vendor.Cfg().OutputReadType {
	case READ_TYPE_COILS:
		result, err = client.ReadCoils(0, 8)
	}

	if err != nil {
		return inputs, outputs, err
	}

	outputs = strconv.FormatUint(uint64(result[0]), 2)

	inputFormat := fmt.Sprintf("%%0%ds", s.vendor.Cfg().InputNum)
	inputs = fmt.Sprintf(inputFormat, inputs)
	inputs = utils.ReverseString(inputs)

	outputFormat := fmt.Sprintf("%%0%ds", s.vendor.Cfg().OutputNum)
	outputs = fmt.Sprintf(outputFormat, outputs)
	outputs = utils.ReverseString(outputs)

	if s.inputs.Load().(string) != inputs {
		s.inputs.Store(inputs)
		s.notify.OnIOStatus(s.cfg.SN, IO_TYPE_INPUT, inputs)
	}

	if s.outputs.Load().(string) != outputs {
		s.outputs.Store(outputs)
		s.notify.OnIOStatus(s.cfg.SN, IO_TYPE_OUTPUT, outputs)
	}

	return inputs, outputs, nil
}

func (s *ModbusTcp) Write(index uint16, status uint16) error {
	if index > (s.vendor.Cfg().OutputNum - 1) {
		return errors.New("invalid index")
	}

	var err error
	client := s.client

	if status == 1 {
		status = 0xFF00
	}

	switch s.vendor.Cfg().WriteType {
	case WRITE_TYPE_SINGLE_COIL:
		_, err = client.WriteSingleCoil(index, index)
	}

	return err
}
