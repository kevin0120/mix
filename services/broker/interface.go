package broker

import "github.com/masami10/rush/services/dispatcherBus"

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherBus.DispatcherMap)
	Release(name string, handler string) error
}
