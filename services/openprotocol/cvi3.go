package openprotocol

type CVI3Controller struct {
	TighteningController
}

func (c *CVI3Controller) defaultControllerGet() IOpenProtocolController {
	c.TighteningController.defaultControllerGet()
	return c
}

// TODO: 可重写所有TighteningController中的方法
