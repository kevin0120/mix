package openprotocol

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/tightening_device"
)

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
}

type IOpenProtocolController interface {
	tightening_device.ITighteningController

	// 初始化控制器
	initController(deviceConfig *tightening_device.TighteningDeviceConfig, d Diagnostic, service *Service, dp Dispatcher)

	// Vendor Model定义(MID，IO等)
	GetVendorModel() map[string]interface{}

	//可设置特定的默认数值
	DefaultControllerGet() IOpenProtocolController

	//控制器状态变化影响相关工具的状态变化
	UpdateToolStatus(status string)

	//根据标识获取工具，通道号或者序列号或者连接(tcp)
	GetToolViaChannel(channel int) (tightening_device.ITighteningTool, error)

	//建立连接
	Connect() error

	//处理未被处理的历史数据
	handlerOldResults() error

	// 加载的协议
	Protocol() string

	//初始化需要订阅的信息
	initSubscribeInfos()

	//执行订阅相关控制器信息
	ProcessSubscribeControllerInfo()

	//曲线解析
	CurveDataDecoding(original []byte, torqueCoefficient float64, angleCoefficient float64, d Diagnostic) (Torque []float64, Angle []float64)
}
