package tightening_device

import (
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
	Register(name string, handler *utils.DispatchHandlerStruct)
}

type ITighteningDevice interface {
	device.IBaseDevice

	//GetDispatcher(string) *utils.Dispatcher
}

type ITighteningController interface {
	ITighteningDevice

	// 定位工具
	GetToolViaSerialNumber(toolSN string) (ITighteningTool, error)

	// 控制输出
	SetOutput(outputs []ControllerOutput) error
}

type ITighteningTool interface {
	ITighteningDevice

	// 工具使能控制
	ToolControl(enable bool) error

	// 设置pset
	SetPSet(pset int) error

	// 设置job
	SetJob(job int) error

	// 模式选择: job/pset
	ModeSelect(mode string) error

	// 取消job
	AbortJob() error

	// 设置pset次数
	SetPSetBatch(pset int, batch int) error

	// pset列表
	GetPSetList() ([]int, error)

	// pset详情
	GetPSetDetail(pset int) (*PSetDetail, error)

	// job列表
	GetJobList() ([]int, error)

	// job详情
	GetJobDetail(job int) (*JobDetail, error)

	// 追溯信息设置
	TraceSet(str string) error
}
