package odoo

import (
	"github.com/masami10/rush/services/aiis"
	"time"
)

type ODOOMOCreated struct {
	ID           int    `json:"id"`
	KNR          string `json:"knr"`
	VIN          string `json:"vin"`
	ProductID    int    `json:"product_id"`
	Result_IDs   []int  `json:"result_ids"`
	AssembleLine int    `json:"assembly_line_id"`
}

type ODOOConsume struct {
	Seq                int64   `json:"sequence"`
	ControllerSN       string  `json:"controller_sn"`
	GunSN              string  `json:"gun_sn"`
	ToleranceMinDegree float64 `json:"tolerance_min_degree"`
	ToleranceMaxDegree float64 `json:"tolerance_max_degree"`
	ToleranceMax       float64 `json:"tolerance_max"`
	ToleranceMin       float64 `json:"tolerance_min"`
	PSet               string  `json:"pset"`
	NutNo              string  `json:"nut_no"`
	Max_redo_times     int     `json:"max_redo_times"`
	X                  float64 `json:"offset_x"`
	Y                  float64 `json:"offset_y"`
	ResultIDs          []int64 `json:"result_ids"`
}

type ODOOWorkorder struct {
	ID          int64  `json:"id"`
	Status      string `json:"status"`
	Max_op_time int    `json:"max_op_time"`

	HMI struct {
		ID   int    `json:"id"`
		UUID string `json:"uuid"`
	} `json:"hmi"`

	Workcenter struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"workcenter"`

	VehicleTypeImg string `json:"vehicleTypeImg"`
	//Worksheet      string        `json:"worksheet"`
	ImageOPID  int64         `json:"img_op_id"`
	KNR        string        `json:"knr"`
	VIN        string        `json:"vin"`
	LongPin    string        `json:"long_pin"`
	UpdateTime time.Time     `json:"update_time"`
	Consumes   []ODOOConsume `json:"consumes"`
	Job        string        `json:"job"`

	// mo相关信息
	MO_EquipemntName  string `json:"equipment_name"` // 设备名
	MO_FactoryName    string `json:"factory_name"`   // 工厂代码
	MO_Year           int64  `json:"year"`
	MO_Pin            int64  `json:"pin"`
	MO_Pin_check_code int64  `json:"pin_check_code"`
	MO_AssemblyLine   string `json:"assembly_line"`
	MO_Lnr            string `json:"lnr"`
	MO_Model          string `json:"model"`
}

type ODOOPoints struct {
	X float64 `json:"offset_x"`
	Y float64 `json:"offset_y"`
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

	Final_pass      string `json:"final_pass"`
	One_time_pass   string `json:"one_time_pass"`
	QualityState    string `json:"quality_state"`
	ExceptionReason string `json:"exception_reason"`
	UserID          int64  `json:"user_id"`
	Batch           string `json:"batch"`
	Vin             string `json:"vin"`

	LocalID int64 `json:"local_id"`
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

type ODOOGun struct {
	ID     int64  `json:"id"`
	Serial string `json:"serial"`
}

type RoutingOperation struct {
	OperationID    int64       `json:"id"`
	Job            int         `json:"job"`
	MaxOpTime      int         `json:"max_op_time"`
	Name           string      `json:"name"`
	Img            string      `json:"img"`
	ProductId      int64       `json:"product_id"`
	WorkcenterID   int64       `json:"workcenter_id"`
	ProductType    string      `json:"product_type"`
	WorkcenterCode string      `json:"workcenter_code"`
	VehicleTypeImg string      `json:"vehicleTypeImg"`
	Points         interface{} `json:"points"`
}
