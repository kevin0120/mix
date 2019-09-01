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

type ControllerOutput struct {
	OutputNo int    `json:"no"`
	Status   string `json:"status"`
}

type TighteningResult struct {
}

type TighteningCurve struct {
}

// api接口
type Reply struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}

type JobSet struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	StepID       int    `json:"step_id"`
	Job          int    `json:"job"`
	UserID       uint   `json:"user_id"`
}

type PSetSet struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	StepID       int64  `json:"step_id"`
	WorkorderID  int64  `json:"workorder_id"`
	UserID       int64  `json:"user_id"`
	PSet         int    `json:"pset"`
	Sequence     uint   `json:"sequence"`
	Count        int    `json:"count"`
	Total        int    `json:"total"`
}

type ToolEnable struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	Enable       bool   `json:"enable"`
}
