package tightening_device

import (
	"github.com/masami10/rush/services/device"
)

type ITighteningDevice interface {
	device.IDevice
}

type ITighteningController interface {
	ITighteningDevice

	// 启动
	Start() error

	// 停止
	Stop() error

	// 定位工具
	GetTool(toolSN string) (ITighteningTool, error)

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
}
