package openprotocol

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
)

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
}

type IStorageService interface {
	FindTargetResultForJobManual(raw_workorder_id int64) (storage.Results, error)
	UpdateTool(gun *storage.Tools) error
	ClearToolResultAndCurve(toolSN string) error
	GetTool(serial string) (storage.Tools, error)
	GetStep(id int64) (storage.Steps, error)
	UpdateIncompleteCurveAndSaveResult(result *storage.Results) error
	UpdateIncompleteResultAndSaveCurve(curve *storage.Curves) error
}

type IBackendService interface {
	TryCreateMaintenance(body interface{}) error
	GetConsumeBySeqInStep(step *storage.Steps, seq int) (*odoo.StepComsume, error)
}

type IOpenProtocolController interface {
	tightening_device.ITighteningController

	// 初始化控制器
	initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service, dp Dispatcher)

	// Vendor Model定义(MID，IO等)
	GetVendorModel() map[string]interface{}

	//控制器状态变化影响相关工具的状态变化
	UpdateToolStatus(status string)

	//根据标识获取工具，通道号或者序列号或者连接(tcp)
	GetToolViaChannel(channel int) (tightening_device.ITighteningTool, error)

	//处理未被处理的历史数据
	handlerOldResults() error

	// 加载的协议
	Protocol() string

	//曲线解析
	CurveDataDecoding(original []byte, torqueCoefficient float64, angleCoefficient float64, d Diagnostic) (Torque []float64, Angle []float64)

	New() IOpenProtocolController
}
