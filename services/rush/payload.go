package rush

import "time"

type OperationResult struct {
	PsetMThreshold float32     `json:"pset_m_threshold"`
	PsetMMax       float32     `json:"pset_m_max"`
	ControlDate    time.Time   `json:"control_date"`
	PsetWMax       float32     `json:"pset_w_max"`
	UserId         int64       `json:"user_id"`
	OneTimePass    string      `json:"one_time_pass"`
	PsetStrategy   string      `json:"pset_strategy"`
	PsetWThreshold float32     `json:"pset_w_threshold"`
	CurObjects     []CURObject `json:"cur_objects"`
	PsetMTarget    float32     `json:"pset_m_target"`
	PsetMMin       float32     `json:"pset_m_min"`
	FinalPass      string      `json:"final_pass"`
	MeasureDegree  float32     `json:"measure_degree"`
	MeasureTDone   float32     `json:"measure_t_don"` //操作所用的时间
	MeasureTorque  float32     `json:"measure_torque"`
	MeasureResult  string      `json:"measure_result"`
	OPTime         int         `json:"op_time"`
	PsetWMin       float32     `json:"pset_w_min"`
	PsetWTarget    float32     `json:"pset_w_target"`
}

//type AIISResult struct {
//	Pset_m_threshold float64 `json:"pset_m_threshold"`
//	Pset_m_max       float64 `json:"pset_m_max"`
//	Pset_m_min       float64 `json:"pset_m_min"`
//	Control_date     string  `json:"control_date"`
//	Pset_w_max       float64 `json:"pset_w_max"`
//	Pset_strategy    string  `json:"pset_strategy"`
//	Pset_m_target    float64 `json:"pset_m_target"`
//	Measure_degree   float64 `json:"measure_degree"`
//	Measure_t_don    float64 `json:"measure_t_don"` //操作所用的时间
//	Measure_torque   float64 `json:"measure_torque"`
//	Measure_result   string  `json:"measure_result"`
//	Op_time          int     `json:"op_time"`
//	Pset_w_min       float64 `json:"pset_w_min"`
//	Pset_w_target    float64 `json:"pset_w_target"`
//	Final_pass       string  `json:"final_pass"`
//	One_time_pass    string  `json:"one_time_pass"`
//
//	CURObjects []CURObject `json:"cur_objects"`
//}

type CURObject struct {
	File string `json:"file"`
	OP   int    `json:"op"`
}
