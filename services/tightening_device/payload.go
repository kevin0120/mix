package tightening_device

// type
const (
	TIGHTENING_DEVICE_TYPE_CONTROLLER = "controller"
	TIGHTENING_DEVICE_TYPE_TOOL       = "tool"
)

const (
	WS_TIGHTENING_DEVICE = "WS_TIGHTENING_DEVICE"

	WS_TOOL_JOB         = "WS_TOOL_JOB"
	WS_TOOL_PSET        = "WS_TOOL_PSET"
	WS_TOOL_RESULT      = "WS_TOOL_RESULT"
	WS_TOOL_ENABLE      = "WS_TOOL_ENABLE"
	WS_TOOL_PSET_LIST   = "WS_TOOL_PSET_LIST"
	WS_TOOL_PSET_DETAIL = "WS_TOOL_PSET_DETAIL"
	WS_TOOL_JOB_LIST    = "WS_TOOL_JOB_LIST"
	WS_TOOL_JOB_DETAIL  = "WS_TOOL_JOB_DETAIL"
)

const (
	TIGHTENING_DEVICE_ONLINE  = "online"
	TIGHTENING_DEVICE_OFFLINE = "offline"
)

type Reply struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}

type JobSet struct {
	ToolSN string `json:"tool_sn"`
	StepID int    `json:"step_id"`
	Job    int    `json:"job"`
	UserID uint   `json:"user_id"`
}

type PSetSet struct {
	ToolSN   string `json:"tool_sn"`
	StepID   int    `json:"step_id"`
	UserID   uint   `json:"user_id"`
	PSet     int    `json:"pset"`
	Sequence uint   `json:"sequence"`
	Count    int    `json:"count"`
}

type ToolEnable struct {
	ToolSN string `json:"tool_sn"`
	Enable bool   `json:"enable"`
}
