package openprotocol

import (
	"github.com/masami10/rush/services/io"
)

type WrenchController struct {
	TighteningController
}

func (c *WrenchController) DefaultControllerGet() IOpenProtocolController {
	c.TighteningController.DefaultControllerGet()
	return c
}

func (c *WrenchController) GetVendorModel() map[string]interface{} {
	vendorModels := map[string]interface{}{
		// *MID							*每个MID对应的REV版本
		MID_0001_START:                 "003",
		MID_0018_PSET:                  "001",
		MID_0014_PSET_SUBSCRIBE:        "001",
		MID_0060_LAST_RESULT_SUBSCRIBE: "006",
		MID_0064_OLD_SUBSCRIBE:         "006",
		MID_0012_PSET_DETAIL_REQUEST:   "002",
		MID_0010_PSET_LIST_REQUEST:     "001",
		MID_0042_TOOL_DISABLE:          "001",
		MID_0043_TOOL_ENABLE:           "001",
		MID_0019_PSET_BATCH_SET:        "001",
		MID_0070_ALARM_SUBSCRIBE:       "001",
		MID_0040_TOOL_INFO_REQUEST:     "005",

		MID_7408_LAST_CURVE_SUBSCRIBE: "001",

		IO_MODEL: io.IoConfig{
			InputNum:  0,
			OutputNum: 0,
		},
	}

	return vendorModels
}

// TODO: 可重写所有TighteningController中的方法
