package hmi

import (
	"github.com/masami10/rush/services/storage"
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
	Release(name string, handler string) error
	Register(name string, handler *utils.DispatchHandlerStruct)
}

type INotifyService interface {
	GetWorkCenter() string
	NotifyAll(evt string, payload string)
}

type IStorageService interface {
	Workorders(par []byte) ([]storage.Workorders, error)
	WorkorderOut(order string, workorderID int64) (interface{}, error)
	UpdateWorkorder(workorder *storage.Workorders) (*storage.Workorders, error)
	UpdateStep(step *storage.Steps) (*storage.Steps, error)
	UpdateStepData(step *storage.Steps) (*storage.Steps, error)
	UpdateOrderData(order *storage.Workorders) (*storage.Workorders, error)
}

type IBackendService interface {
	GetWorkorder(masterpcSn string, hmiSn string, workcenterCode, code string) ([]byte, error)
}
