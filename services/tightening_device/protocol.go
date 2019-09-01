package tightening_device

type ITighteningProtocol interface {

	// 协议名称
	Name() string

	// 创建控制器
	CreateController(cfg *TighteningDeviceConfig) (ITighteningController, error)

	// 判断协议是否支持
	Support(cfg *TighteningDeviceConfig) error
}
