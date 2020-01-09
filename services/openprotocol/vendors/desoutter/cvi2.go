package desoutter

import (
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/openprotocol"
)

type CVI2Controller struct {
	openprotocol.TighteningController
}

func (c *CVI2Controller) New() openprotocol.IOpenProtocolController {
	controller := CVI2Controller{}
	controller.SetInstance(&controller)
	return &controller
}

func (c *CVI2Controller) GetVendorModel() map[string]interface{} {
	vendorModels := map[string]interface{}{
		// *MID							*每个MID对应的REV版本
		openprotocol.MID_0001_START:                 "001",
		openprotocol.MID_0018_PSET:                  "001",
		openprotocol.MID_0014_PSET_SUBSCRIBE:        "001",
		openprotocol.MID_0060_LAST_RESULT_SUBSCRIBE: "001",
		openprotocol.MID_0012_PSET_DETAIL_REQUEST:   "001",
		openprotocol.MID_0010_PSET_LIST_REQUEST:     "001",
		openprotocol.MID_0042_TOOL_DISABLE:          "001",
		openprotocol.MID_0043_TOOL_ENABLE:           "001",
		openprotocol.MID_0019_PSET_BATCH_SET:        "001",
		openprotocol.MID_0070_ALARM_SUBSCRIBE:       "001",
		openprotocol.MID_0040_TOOL_INFO_REQUEST:     "001",
		openprotocol.MID_0210_INPUT_SUBSCRIBE:       "001",
		openprotocol.MID_0051_VIN_SUBSCRIBE:         "001",

		openprotocol.IoModel: io.IoConfig{
			InputNum:  0,
			OutputNum: 0,
		},
	}

	return vendorModels
}

// TODO: 可重写所有TighteningController中的方法
