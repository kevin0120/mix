package openprotocol

type CVI2Controller struct {
	TighteningController
}

func (c *CVI2Controller) defaultControllerGet() IOpenProtocolController {
	c.TighteningController.defaultControllerGet()
	return c
}

// TODO: 可重写所有TighteningController中的方法
