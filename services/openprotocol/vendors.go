package openprotocol

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/tightening_device"
)

const (
	IO_CONFIG = "IO_CONFIG"
)

type IOpenProtocolController interface {
	tightening_device.ITighteningController
	defaultControllerGet() *TighteningController                              //可设置特定的默认数值
	UpdateToolStatus(status string)                                           //控制器状态变化影响相关工具的状态变化
	GetToolViaChannel(channel int) (tightening_device.ITighteningTool, error) //根据标识获取工具，通道号或者序列号或者连接(tcp)
	Connect() error                                                           //建立连接
	handlerOldResults() error                                                 //处理未被处理的历史数据
	Protocol() string
	initSubscribeInfos()                                                                                                                      //初始化需要订阅的信息
	ProcessSubscribeControllerInfo()                                                                                                          //执行订阅相关控制器信息
	CurveDataDecoding(original []byte, torqueCoefficient float64, angleCoefficient float64, d Diagnostic) (Torque []float64, Angle []float64) //曲线解析
}

// model[mid][rev]
var VendorModels = map[string]map[string]interface{}{
	tightening_device.ModelDesoutterCvi3: {
		// *MID							  *每个MID对应的REV版本
		MID_0001_START:                   "004",
		MID_0018_PSET:                    "001",
		MID_0014_PSET_SUBSCRIBE:          "001",
		MID_0034_JOB_INFO_SUBSCRIBE:      "004",
		MID_0250_SELECTOR_SUBSCRIBE:      "001",
		MID_0060_LAST_RESULT_SUBSCRIBE:   "998",
		MID_0150_IDENTIFIER_SET:          "001",
		MID_0038_JOB_SELECT:              "002",
		MID_0064_OLD_SUBSCRIBE:           "006",
		MID_0130_JOB_OFF:                 "001",
		MID_0012_PSET_DETAIL_REQUEST:     "002",
		MID_0010_PSET_LIST_REQUEST:       "001",
		MID_0032_JOB_DETAIL_REQUEST:      "003",
		MID_0030_JOB_LIST_REQUEST:        "002",
		MID_0042_TOOL_DISABLE:            "001",
		MID_0043_TOOL_ENABLE:             "001",
		MID_0200_CONTROLLER_RELAYS:       "001",
		MID_0019_PSET_BATCH_SET:          "001",
		MID_0210_INPUT_SUBSCRIBE:         "001",
		MID_0127_JOB_ABORT:               "001",
		MID_0100_MULTI_SPINDLE_SUBSCRIBE: "001",
		MID_0051_VIN_SUBSCRIBE:           "002",
		MID_0070_ALARM_SUBSCRIBE:         "001",
		MID_0040_TOOL_INFO_REQUEST:       "002",

		MID_7408_LAST_CURVE_SUBSCRIBE: "001",

		IO_CONFIG: io.IoConfig{
			InputNum:  8,
			OutputNum: 8,
		},
	},

	tightening_device.ModelDesoutterDeltaWrench: {
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

		IO_CONFIG: io.IoConfig{
			InputNum:  0,
			OutputNum: 0,
		},
	},
}

func GetModel(model string) (map[string]interface{}, error) {
	m, exist := VendorModels[model]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Model %s Not Supported", model))
	}

	return m, nil
}

func GetVendorMid(model string, mid string) (string, error) {
	m, err := GetModel(model)
	if err != nil {
		return "", err
	}

	rev, exist := m[mid]
	if !exist {
		return "", errors.New(fmt.Sprintf("MID %s Not Supported", mid))
	}

	return rev.(string), nil
}
