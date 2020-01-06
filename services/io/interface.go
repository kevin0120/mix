package io

import (
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/utils"
)

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	Release(name string, handler string) error
	Register(name string, handler *utils.DispatchHandlerStruct)
}

type IONotify interface {
	OnStatus(sn string, status string)
	OnChangeIOStatus(sn string, t string, status string)
}

type IO interface {
	Status() string
	Start() error
	Stop() error
	Read() (string, string, error)
	Write(uint16, uint16) error
}

type IDeviceService = device.IDeviceService
