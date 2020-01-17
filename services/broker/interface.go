package broker

import "github.com/masami10/rush/services/dispatcherbus"

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
	ReleaseDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
}

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}
