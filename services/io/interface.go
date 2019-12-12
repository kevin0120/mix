package io

import "github.com/masami10/rush/services/dispatcherbus"

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
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
