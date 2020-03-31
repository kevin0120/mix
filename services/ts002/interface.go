package ts002

import (
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
}

type IHttpService interface {
	GetHandlerByName(version string) (*httpd.Handler, error)
}

type IIOService interface {
	GetIOSerialNumberByIdx(index int) string
	Write(sn string, index uint16, status uint16) error
}

type IDispatcher interface {
	Dispatch(name string, data interface{}) error
	Register(name string, handler *utils.DispatchHandlerStruct)
}

type ITightening interface {
	ToolPSetByIP(req *tightening_device.PSetSet) error
}

type IStorage interface {
	UpdateToolLocation(toolSN string, location string) error
	GetToolLocation(toolSN string) (string, error)
}
