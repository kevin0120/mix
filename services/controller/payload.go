package controller

import "github.com/masami10/rush/services/aiis"

const (
	AUDIPROTOCOL         = "Audi/VW"
	OPENPROTOCOL         = "OpenProtocol"
	DEFAULT_TOOL_CHANNEL = 1
	AUTO_MODE            = "auto"
)

const (
	STATUS_ONLINE  = "online"
	STATUS_OFFLINE = "offline"

	ERR_CONTROLER_NOT_FOUND = "controller not found"
	ERR_CONTROLER_TIMEOUT   = "controller timeout"
	ERR_NOT_FOUND           = "not found"
	ERR_PSET_ERROR          = "pset error"
	ERR_KNOWN               = "error known"
)

const (
	STRATEGY_AD  = "AD"
	STRATEGY_AW  = "AW"
	STRATEGY_ADW = "ADW"
)

const (
	CONTROLLER_CVI3  = "cvi3"
	CONTROLLER_CVIL2 = "cvil2"
)

//type ControllerCurve struct {
//	ResultID     int64
//	CurveFile    string
//	CurveContent ControllerCurveFile
//	Count        int
//	UpdateTime   string
//}
//
//type ControllerCurveFile struct {
//	Result string    `json:"result"`
//	CUR_M  []float64 `json:"cur_m"`
//	CUR_W  []float64 `json:"cur_w"`
//	CUR_T  []float64 `json:"cur_t"`
//}

type ControllerResult struct {
	Result_id     int64  `json:"result_id"`
	Controller_SN string `json:"controller_sn"`
	Workorder_ID  int64  `json:"workorder_id"`
	UserID        int64  `json:"user_id"`
	CurFile       aiis.CURObject
	Result        string     `json:"result"`
	Dat           string     `json:"dat"`
	PSet          int        `json:"pset"`
	Batch         string     `json:"batch"`
	Count         int        `json:"count"`
	PSetDefine    PSetDefine `json:"pset_define"`
	GunSN         string

	ResultValue  ResultValue `json:"result_value"`
	TighteningID string

	NeedPushAiis bool
	NeedPushHmi  bool

	ExceptionReason string
}

type PSetDefine struct {
	Strategy string  `json:"strategy"`
	Mp       float64 `json:"M+"`
	Mm       float64 `json:"M-"`
	Ms       float64 `json:"MS"`
	Ma       float64 `json:"MA"`
	Wp       float64 `json:"W+"`
	Wm       float64 `json:"W-"`
	Wa       float64 `json:"WS"`
}

type ResultValue struct {
	Mi float64 `json:"MI"`
	Wi float64 `json:"WI"`
	Ti float64 `json:"TI"`
}

type CurveObject struct {
	File  string
	Count int
}
