package aiis

const (
	WS_REG    = "reg"
	WS_RESULT = "result"
)

const (
	ODOO_RESULT_PASS = "pass"
	ODOO_RESULT_FAIL = "fail"

	QUALITY_STATE_PASS = "pass"
	QUALITY_STATE_FAIL = "fail"
	QUALITY_STATE_EX   = "exception"

	HANDLER_TYPE_CURVE = "curve"
	HANDLER_TYPE_AIIS  = "aiis"

	ODOO_STATUS_ONLINE  = "online"
	ODOO_STATUS_OFFLINE = "offline"
)

type AIISResult struct {
	// local_id
	ID int64 `json:"id"`

	// db
	Vin              string      `json:"vin"`
	Pset_strategy    string      `json:"pset_strategy"`
	Pset_m_max       float64     `json:"pset_m_max"`
	Pset_m_min       float64     `json:"pset_m_min"`
	Pset_m_threshold float64     `json:"pset_m_threshold"`
	Pset_m_target    float64     `json:"pset_m_target"`
	Pset_w_max       float64     `json:"pset_w_max"`
	Pset_w_min       float64     `json:"pset_w_min"`
	Pset_w_threshold float64     `json:"pset_w_threshold"`
	Pset_w_target    float64     `json:"pset_w_target"`
	CURObjects       []CURObject `json:"cur_objects"`
	QualityState     string      `json:"quality_state"`
	ExceptionReason  string      `json:"exception_reason"`
	Job              string      `json:"job"`
	Control_date     string      `json:"control_date"`
	Measure_torque   float64     `json:"measure_torque"`
	Measure_degree   float64     `json:"measure_degree"`
	Measure_t_don    float64     `json:"measure_t_don"`
	Measure_result   string      `json:"measure_result"`
	Lacking          string      `json:"lacking"`
	Op_time          int         `json:"op_time"`
	One_time_pass    string      `json:"one_time_pass"`
	Final_pass       string      `json:"final_pass"`
	Batch            string      `json:"batch"`

	//db_fk
	UserID       int64 `json:"user_id"`
	ProductID    int64 `json:"product_id"`
	WorkcenterID int64 `json:"workcenter_id"`
	GunID        int64 `json:"gun_id"`
	NutID        int64 `json:"nut_id"`

	// mo相关信息
	MO_EquipemntName  string `json:"equipment_name"` // 设备名
	MO_FactoryName    string `json:"factory_name"`   // 工厂代码
	MO_Year           int64  `json:"year"`
	MO_Pin            int64  `json:"pin"`
	MO_Pin_check_code int64  `json:"pin_check_code"`
	MO_AssemblyLine   string `json:"assembly_line"`
	MO_Lnr            string `json:"lnr"`
	MO_NutNo          string `json:"nut_no"`
	MO_Model          string `json:"model"`

	// others
	Seq            int    `json:"seq"`
	Mode           string `json:"mode"`
	ControllerSN   string `json:"controller_sn"`
	TighteningId   int64  `json:"tightening_id"`
	ToolSN         string `json:"tool_sn"`
	WorkcenterCode string `json:"workcenter_code"`
}

type CURObject struct {
	File string `json:"file"`
	OP   int    `json:"op"`
}

type WSMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WSRegist struct {
	Rush_SN string `json:"rush_sn"`
}

type WSOpResult struct {
	ResultID int64      `json:"result_id"`
	Result   AIISResult `json:"result"`
	Port     string     `json:"port"`
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

type ResultPatch struct {
	ID        int64 `json:"id"`
	HasUpload bool  `json:"has_upload"`
}

type ODOOStatus struct {
	Status string `json:"status"`
}
