package tightening_device

// type
const (
	TIGHTENING_DEVICE_TYPE_CONTROLLER = "controller"
	TIGHTENING_DEVICE_TYPE_TOOL       = "tool"
	MODE_PSET                         = "pset"
	MODE_JOB                          = "job"
)

const (
	WS_TIGHTENING_DEVICE = "WS_TIGHTENING_DEVICE"

	WS_TOOL_JOB         = "WS_TOOL_JOB"
	WS_TOOL_PSET        = "WS_TOOL_PSET"
	WS_TOOL_RESULT      = "WS_TOOL_RESULT"
	WS_TOOL_ENABLE      = "WS_TOOL_ENABLE"
	WS_TOOL_MODE_SELECT = "WS_TOOL_MODE_SELECT"
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

type ToolModeSelect struct {
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	Mode         string `json:"mode"`
}

type PSetDetail struct {
	PSetID            int     `json:"pset"`
	PSetName          string  `json:"pset_name"`
	RotationDirection string  `json:"rotation_direction"`
	BatchSize         int     `json:"batch_size"`
	TorqueMin         float64 `json:"torque_min"`
	TorqueMax         float64 `json:"torque_max"`
	TorqueTarget      float64 `json:"torque_target"`
	AngleMin          float64 `json:"angle_min"`
	AngleMax          float64 `json:"angle_max"`
	AngleTarget       float64 `json:"angle_target"`
}

type JobDetail struct {
	JobID             int       `json:"job"`
	JobName           string    `json:"job_name"`
	OrderStrategy     string    `json:"order_strategy"`
	CountType         string    `json:"count_type"`
	LockAtJobDone     bool      `json:"lock_at_job_done"`
	UseLineControl    bool      `json:"use_line_control"`
	RepeatJob         bool      `json:"repeat_job"`
	LooseningStrategy string    `json:"loosening_strategy"`
	Steps             []JobStep `json:"steps"`
}

type JobStep struct {
	StepName  string `json:"step_name"`
	ChannelID int    `json:"channel_id"`
	PSetID    int    `json:"pset_id"`
	BatchSize int    `json:"batch_size"`
	Socket    int    `json:"socket"`
}
