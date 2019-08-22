package controller

import "github.com/masami10/rush/services/aiis"

const (
	AUDIPROTOCOL         = "Audi/VW"
	OPENPROTOCOL         = "OpenProtocol"
	DEFAULT_TOOL_CHANNEL = 1
	AUTO_MODE            = "auto"

	WS_TIGHTENING_RESULT = "WS_TIGHTENING_RESULT"
)

const (
	STATUS_ONLINE  = "online"
	STATUS_OFFLINE = "offline"

	STATUS_ENABLE  = "enable"
	STATUS_DISABLE = "disable"

	ERR_CONTROLER_NOT_FOUND = "controller not found"
	ERR_TOOL_NOT_FOUND      = "tool not found"
	ERR_CONTROLER_TIMEOUT   = "controller timeout"
	ERR_NOT_FOUND           = "not found"
	ERR_PSET_ERROR          = "pset error"
	ERR_KNOWN               = "error known"
	ERR_NOT_SUPPORTED       = "not supported"
)

const (
	EVT_TOOL_CONNECTED    = "connected"
	EVT_TOOL_DISCONNECTED = "disconnected"
)

const (
	STRATEGY_AD  = "AD"
	STRATEGY_AW  = "AW"
	STRATEGY_ADW = "ADW"
	STRATEGY_LN  = "LN"
	STRATEGY_AK2 = "AK2"
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
	Seq           int        `json:"sequence"`
	GroupSeq      int        `json:"group_sequence"`
	Count         int        `json:"count"`
	PSetDefine    PSetDefine `json:"pset_define"`
	GunSN         string     `json:"gun_sn"`

	ResultValue  ResultValue `json:"result_value"`
	MaxRedoTime  int
	TighteningID string

	NeedPushAiis bool
	NeedPushHmi  bool

	ExceptionReason string

	Raw    string
	StepID int64
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

type Consume struct {
	Seq                int     `json:"sequence"`
	GroupSeq           int     `json:"group_sequence"`
	ControllerSN       string  `json:"controller_sn"`
	GunSN              string  `json:"gun_sn"`
	ToleranceMinDegree float64 `json:"tolerance_min_degree"`
	ToleranceMaxDegree float64 `json:"tolerance_max_degree"`
	ToleranceMax       float64 `json:"tolerance_max"`
	ToleranceMin       float64 `json:"tolerance_min"`
	PSet               string  `json:"pset"`
	NutNo              string  `json:"nut_no"`
	Max_redo_times     int     `json:"max_redo_times"`
	X                  float64 `json:"offset_x"`
	Y                  float64 `json:"offset_y"`
	ResultIDs          []int64 `json:"result_ids"`
}
