package openprotocol

type WrenchController struct {
	TighteningController
}

func (c *WrenchController) defaultControllerGet() IOpenProtocolController {
	c.TighteningController.defaultControllerGet()
	return c
}

// TODO: 可重写所有TighteningController中的方法
