package tightening_device

const (
	DISPATCH_RESULT = "DISPATCH_RESULT"
	DISPATCH_CURVE  = "DISPATCH_CURVE"
)

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
	ControllerSN string `json:"controller_sn"`
	ToolSN       string `json:"tool_sn"`
	Dat          string `json:"dat"`
	Result       string `json:"result"`
	Job          int    `json:"job"`
	PSet         int    `json:"pset"`
	Batch        string `json:"batch"`
	Count        int    `json:"count"`

	Strategy string  `json:"strategy"`
	Mp       float64 `json:"M+"`
	Mm       float64 `json:"M-"`
	Ms       float64 `json:"MS"`
	Ma       float64 `json:"MA"`
	Wp       float64 `json:"W+"`
	Wm       float64 `json:"W-"`
	Wa       float64 `json:"WS"`

	Mi float64 `json:"MI"`
	Wi float64 `json:"WI"`
	Ti float64 `json:"TI"`

	TighteningID string
}

type TighteningCurve struct {
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
