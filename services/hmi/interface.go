package hmi

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
	Disconnect(id string)
	Close()
	Closed()
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
	Register(name string, handler *utils.DispatchHandlerStruct)
}

type INotifyService interface {
	GetWorkCenter() string
	NotifyAll(evt string, payload string)
}