package odoo

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/utils"
)

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	Release(name string, handler string) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	ReleaseDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Register(name string, handler *utils.DispatchHandlerStruct)
}

type IStorageService interface {
	WorkorderExists(id int64) (bool, error)
	GetWorkOrder(id int64, raw bool) (storage.Workorders, error)
	InsertWorkorder(workorder *storage.Workorders, results *[]storage.Results, checkWorkorder bool, checkResult bool, rawid bool) error
	WorkorderIn(in []byte) (string, error)
	WorkorderOut(order string, workorderID int64) (interface{}, error)
	DeleteRoutingOperations(rds []storage.RoutingOperationDelete) error
	DeleteAllRoutingOperations() error
	GetRoutingOperations(name string, model string) (storage.RoutingOperations, error)
	Store(data interface{}) error
	UpdateRoutingOperations(ro *storage.RoutingOperations) error
}

type HTTPService interface {
	AddNewHttpHandler(r httpd.Route)
}
