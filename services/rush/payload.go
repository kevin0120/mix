package rush

import (
	"github.com/masami10/aiis/services/storage"
)

const (
	WS_REG    = "reg"
	WS_RESULT = "result"
)

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

type WSMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WSRegist struct {
	Rush_SN string `json:"rush_sn"`
}

type WSOpResult struct {
	ResultID int64                   `json:"result_id"`
	Result   storage.OperationResult `json:"result"`
	Port     string                  `json:"port"`
}
