package hmi

type LocalResults struct {
	HmiSN        interface{} `json:"hmi_sn,omitempty"`
	Vin          interface{} `json:"vin,omitempty"`
	VehicleType  interface{} `json:"vehicle_type,omitempty"`
	JobID        interface{} `json:"job_id,omitempty"`
	PSetID       interface{} `json:"pset_id,omitempty"`
	ControllerSN interface{} `json:"controller_sn,omitempty"`
	ToolSN       interface{} `json:"tool_sn,omitempty"`
	Result       interface{} `json:"result,omitempty"`
	Torque       interface{} `json:"torque,omitempty"`
	Angle        interface{} `json:"angle,omitempty"`
	Spent        interface{} `json:"spent,omitempty"`
	TimeStamp    interface{} `json:"timestamp,omitempty"`
	Batch        interface{} `json:"batch,omitempty"`
}

type WSOrderReq struct {
	ID            int64  `json:"id"`
	WorkorderCode string `json:"workorder_code"`
	Status        string `json:"status"`
}

type WSOrderReqData struct {
	ID            int64  `json:"id"`
	WorkorderCode string `json:"workorder_code"`
	Data          string `json:"data"`
}

type WSStepReq struct {
	ID       int64  `json:"id"`
	StepCode string `json:"workstep_code"`
	Status   string `json:"status"`
}

type WSStepReqData struct {
	ID       int64  `json:"id"`
	StepCode string `json:"workstep_code"`
	Data     string `json:"data"`
}

type WSOrderReqCode struct {
	Code       string `json:"code"`
	Workcenter string `json:"workcenter"`
}

type WSWorkcenter struct {
	WorkCenter string `json:"workcenter"`
}

type WSLocalResults struct {
	HmiSN   string   `json:"hmi_sn"`
	Filters []string `json:"filters"`
	Limit   int      `json:"limit"`
}
