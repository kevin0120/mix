package aiis

const (
	ServiceAiis = "aiis"
)

type TransportPayload struct {
	Method string      `json:"method"`
	Data   interface{} `json:"data"`
}

type PublishResult struct {
	ID int64 `json:"id"`

	// db
	Vin             string      `json:"vin"`
	PsetStrategy    string      `json:"pset_strategy"`
	PsetMMax        float64     `json:"pset_m_max"`
	PsetMMin        float64     `json:"pset_m_min"`
	PsetMThreshold  float64     `json:"pset_m_threshold"`
	PsetMTarget     float64     `json:"pset_m_target"`
	PsetWMax        float64     `json:"pset_w_max"`
	PsetWMin        float64     `json:"pset_w_min"`
	PsetWThreshold  float64     `json:"pset_w_threshold"`
	PsetWTarget     float64     `json:"pset_w_target"`
	CURObjects      []CURObject `json:"cur_objects"`
	QualityState    string      `json:"quality_state"`
	ExceptionReason string      `json:"exception_reason"`
	Job             string      `json:"job"`
	ControlDate     string      `json:"control_date"`
	MeasureTorque   float64     `json:"measure_torque"`
	MeasureDegree   float64     `json:"measure_degree"`
	MeasureTDon     float64     `json:"measure_t_don"`
	MeasureResult   string      `json:"measure_result"`
	Lacking         string      `json:"lacking"`
	OpTime          int         `json:"op_time"`
	OneTimePass     string      `json:"one_time_pass"`
	FinalPass       string      `json:"final_pass"`
	Batch           string      `json:"batch"`

	//db_fk
	UserID       int64 `json:"user_id"`
	ProductID    int64 `json:"product_id"`
	WorkcenterID int64 `json:"workcenter_id"`
	GunID        int64 `json:"gun_id"`
	NutID        int64 `json:"consu_product_id"`

	// mo相关信息
	MoEquipemntname string `json:"equipment_name"` // 设备名
	MoFactoryname   string `json:"factory_name"`   // 工厂代码
	MoYear          int64  `json:"year"`
	MoPin           int64  `json:"pin"`
	MoPinCheckCode  int64  `json:"pin_check_code"`
	MoAssemblyline  string `json:"assembly_line"`
	MoLnr           string `json:"lnr"`
	MoNutno         string `json:"nut_no"`
	MoModel         string `json:"model"`

	// others
	Seq            int    `json:"seq"`
	Mode           string `json:"mode"`
	ControllerSN   string `json:"controller_sn"`
	TighteningId   string  `json:"tightening_id"`
	ToolSN         string `json:"tool_sn"`
	WorkcenterCode string `json:"workcenter_code"`
	Stage          string `json:"stage"`
	WorkorderID    int64  `json:"workorder_id"`
	WorkorderName  string `json:"workorder_name"`

	Payload string `json:"payload"`
}

type CURObject struct {
	File string `json:"file"`
	OP   int    `json:"op"`
}

type ResultPatch struct {
	ID        int64 `json:"id"`
	HasUpload bool  `json:"has_upload"`
}

// 服务状态(aiis, odoo, 外部系统等)
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
