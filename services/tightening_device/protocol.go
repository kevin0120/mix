package tightening_device

type ITighteningProtocol interface {

	// 协议名称
	Name() string

	// 创建控制器
	CreateController(cfg *TighteningDeviceConfig, dp Dispatcher) (ITighteningController, error)
}
