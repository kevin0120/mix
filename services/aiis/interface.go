package aiis

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/storage"
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
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
}

type IStorageService interface {
	UpdateResultByCount(id int64, count int, flag bool) error
	GetResultByID(id int64) (*storage.Results, error)
	GetWorkOrder(id int64, raw bool) (storage.Workorders, error)
	ListUnUploadResults() ([]storage.Results, error)
}