package odoo

import (
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/hmi"
)

type ResultPatch struct {
	HasUpload bool `json:"has_upload"`
}

type ODOOMOCreated struct {
	ID           int    `json:"id"`
	KNR          string `json:"knr"`
	VIN          string `json:"vin"`
	ProductID    int    `json:"product_id"`
	Result_IDs   []int  `json:"result_ids"`
	AssembleLine int    `json:"assembly_line_id"`
}

type ODOOWorkorder struct {
	ID             int64   `json:"id"`
	Status         string  `json:"status"`
	NutTotal       float64 `json:"nut_total"`
	PSet           string  `json:"pset"`
	Max_redo_times int     `json:"max_redo_times"`
	Max_op_time    int     `json:"max_op_time"`

	HMI struct {
		ID   int    `json:"id"`
		UUID string `json:"uuid"`
	} `json:"hmi"`

	Result_IDs []int64       `json:"result_ids"`
	Worksheet  hmi.Worksheet `json:"worksheet"`
	KNR        string        `json:"knr"`
	VIN        string        `json:"vin"`
}

type ODOOCurveAppend struct {
	File string `json:"file"`
	OP   int    `json:"op"`
}

type ODOOResultSync struct {
	ID               int64   `json:"id"`
	Pset_m_threshold float64 `json:"pset_m_threshold"`
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

	CURObjects []aiis.CURObject `json:"cur_objects"`
}

type ODOOPR struct {
	Pr_value string `json:"pr_value"`
	Pr_group string `json:"pr_group"`
}

type ODOOMO struct {
	Pin_check_code int    `json:"pin_check_code"`
	Equipment_name string `json:"equipment_name"`
	Factory_name   string `json:"factory_name"`
	Pin            int    `json:"pin"`
	Year           int    `json:"year"`
	Assembly_line  string `json:"assembly_line"`
	Model          string `json:"model"`
	//Model2 string `json:"model"`
	Vin                string   `json:"vin"`
	Lnr                string   `json:"lnr"`
	Date_planned_start string   `json:"date_planned_start"`
	Prs                []ODOOPR `json:"prs"`
}
