package aiis

import (
	"github.com/masami10/rush/services/DispatcherBus"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
	PutResultDone()
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherBus.DispatcherMap)
	Release(name string, handler string) error
}
