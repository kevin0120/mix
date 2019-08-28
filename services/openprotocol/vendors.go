package openprotocol

import "github.com/masami10/rush/services/tightening_device"

// model[mid][rev]
var VendorModels = map[string]map[string]string{
	tightening_device.MODEL_DESOUTTER_CVI3: {
		MID_0001_START:                   "003",
		MID_0018_PSET:                    "001",
		MID_0014_PSET_SUBSCRIBE:          "001",
		MID_0034_JOB_INFO_SUBSCRIBE:      "003",
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
		MID_0040_TOOL_INFO_REQUEST:       "004",

		MID_7408_LAST_CURVE_SUBSCRIBE: "001",
	},

	tightening_device.MODEL_DESOUTTER_DELTA_WRENCH: {
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
	},
}

func GetVendorMid(model string, mid string) string {
	m, has := VendorModels[model]
	if !has {
		return ""
	}

	rev, has := m[mid]
	if !has {
		return ""
	}

	return rev
}
