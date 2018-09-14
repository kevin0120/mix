package aiis

type AIISResult struct {
	Pset_m_threshold float64 `json:"pset_m_threshold"`
	Pset_w_threshold float64 `json:"pset_w_threshold"`
	UserID           int64   `json:"user_id"`
	Pset_m_max       float64 `json:"pset_m_max"`
	Pset_m_min       float64 `json:"pset_m_min"`
	Control_date     string  `json:"control_date"`
	Pset_w_max       float64 `json:"pset_w_max"`
	Pset_strategy    string  `json:"pset_strategy"`
	Pset_m_target    float64 `json:"pset_m_target"`
	Measure_degree   float64 `json:"measure_degree"`
	Measure_t_don    float64 `json:"measure_t_don"`
	Measure_torque   float64 `json:"measure_torque"`
	Measure_result   string  `json:"measure_result"`
	Op_time          int     `json:"op_time"`
	Pset_w_min       float64 `json:"pset_w_min"`
	Pset_w_target    float64 `json:"pset_w_target"`
	Final_pass       string  `json:"final_pass"`
	One_time_pass    string  `json:"one_time_pass"`
	QualityState     string  `json:"quality_state"`
	ExceptionReason  string  `json:"exception_reason"`
	Seq              int     `json:"seq"`
	ProductID        int64   `json:"product_id"`
	WorkcenterID     int64   `json:"workcenter_id"`
	GunID            int64   `json:"gun_id"`
	Batch			string `json:"batch"`

	CURObjects []CURObject `json:"cur_objects"`

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

	Vin string `json:"vin"`
}

type CURObject struct {
	File string `json:"file"`
	OP   int    `json:"op"`
}
