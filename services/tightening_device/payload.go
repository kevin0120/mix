package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/storage"
	"time"
)

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
	// 控制器序列号
	ControllerSN string `json:"controller_sn"`

	// 工具序列号
	ToolSN string `json:"tool_sn"`

	// 收到时间
	UpdateTime time.Time `json:"update_time"`

	// job号
	Job int `json:"job"`

	// pset号
	PSet int `json:"pset"`

	// 批次信息
	Batch string `json:"batch"`

	// 当前拧紧次数
	Count int `json:"count"`

	// 拧紧ID
	TighteningID string `json:"tightening_id"`

	// 实际结果
	MeasureResult string `json:"measure_result"`

	// 实际扭矩
	MeasureTorque float64 `json:"measure_torque"`

	// 实际角度
	MeasureAngle float64 `json:"measure_angle"`

	// 实际耗时
	MeasureTime float64 `json:"measure_time"`

	// 拧紧策略
	Strategy string `json:"strategy"`

	// 最大扭矩
	TorqueMax float64 `json:"torque_max"`

	// 最小扭矩
	TorqueMin float64 `json:"torque_min-"`

	// 扭矩阈值
	TorqueThreshold float64 `json:"torque_threshold"`

	// 目标扭矩
	TorqueTarget float64 `json:"torque_target"`

	// 最大角度
	AngleMax float64 `json:"angle_max"`

	// 最小角度
	AngleMin float64 `json:"angle_min"`

	// 目标角度
	AngleTarget float64 `json:"angle_target"`
}

func (r *TighteningResult) ToDBResult() *storage.Results {
	return &storage.Results{
		HasUpload:    false,
		Result:       r.MeasureResult,
		ToolSN:       r.ToolSN,
		ControllerSN: r.ControllerSN,
		TighteningID: r.TighteningID,
	}
}

type TighteningCurve struct {
	// 工具序列号
	ToolSN string `json:"tool_sn"`

	// 拧紧ID
	TighteningID string `json:"tightening_id"`

	// 收到时间
	UpdateTime time.Time `json:"update_time"`

	TighteningCurveContent
}

func (c *TighteningCurve) ToDBCurve() *storage.Curves {
	curveContent, _ := json.Marshal(c.TighteningCurveContent)

	return &storage.Curves{
		HasUpload:  false,
		UpdateTime: c.UpdateTime,
		CurveData:  string(curveContent),
		CurveFile:  fmt.Sprintf("%s_%s.json", c.ToolSN, c.TighteningID),
	}
}

type TighteningCurveContent struct {
	// 实际拧紧结果(ok/nok)
	Result string `json:"result"`

	CUR_M []float64 `json:"cur_m"`
	CUR_W []float64 `json:"cur_w"`
	CUR_T []float64 `json:"cur_t"`
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
