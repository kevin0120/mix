package tightening_device

type TighteningProtocol interface {
	CreateController(cfg *DeviceConfig) TighteningController
}
