package io

import ()

const (
	MODEL_MOXA_E1212 = "MOXA_E1212"
)

const (
	READ_TYPE_DISCRETES         = "READ_TYPE_DISCRETES"
	READ_TYPE_COILS             = "READ_TYPE_COILS"
	READ_TEYP_HOLDING_REGISTERS = "READ_TEYP_HOLDING_REGISTERS"
	READ_TEYP_INPUT_REGISTERS   = "READ_TEYP_INPUT_REGISTERS"

	WRITE_TYPE_SINGLE_COIL     = "WRITE_TYPE_SINGLE_COIL"
	WRITE_TYPE_MULTI_COIL      = "WRITE_TYPE_MULTI_COIL"
	WRITE_TYPE_SINGLE_REGISTER = "WRITE_TYPE_SINGLE_REGISTER"
	WRITE_TYPE_MULTI_REGISTER  = "WRITE_TYPE_MULTI_REGISTER"
)

type VendorCfg struct {
	InputNum      uint16
	InputAddress  uint16
	InputReadType string

	OutputNum      uint16
	OutputAddress  uint16
	OutputReadType string

	WriteType string
}

type Vendor interface {
	Type() string
	Cfg() VendorCfg
}

var VENDOR_MODELS = map[string]Vendor{
	MODEL_MOXA_E1212: &MOXA{model: MODEL_MOXA_E1212},
}

type MOXA struct {
	model string
}

func (s *MOXA) Type() string {
	switch s.model {
	case MODEL_MOXA_E1212:
		return IO_MODBUSTCP
	}

	return IO_MODBUSTCP
}

// inputNum, outputNum
func (s *MOXA) Cfg() VendorCfg {
	switch s.model {
	case MODEL_MOXA_E1212:
		return VendorCfg{
			InputNum:       8,
			InputAddress:   0,
			InputReadType:  READ_TYPE_DISCRETES,
			OutputNum:      8,
			OutputAddress:  0,
			OutputReadType: READ_TYPE_COILS,
			WriteType:      WRITE_TYPE_SINGLE_COIL,
		}
	}

	return VendorCfg{}
}
