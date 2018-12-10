package storage

import (
	"github.com/masami10/aiis/services/aiis"
	"time"
)

//odoo中的operation_result数据库 模型
//type OperationResultModel struct {
//	Id             int64     `xorm:"pk SERIAL notnull 'id'"`
//	CreateTime     time.Time `xorm:"DATE notnull 'time'"`
//	PsetMThreshold float32   `xorm:"numeric 'pset_m_threshold'"`
//	PsetMMax       float32   `xorm:"numeric 'pset_m_max'"`
//	FinalPass      string    `xorm:"varchar(32) 'final_pass'"`
//	OneTimePass    string    `xorm:"varchar(32) 'one_time_pass'"`
//	//PointId         int64     `xorm:" INTEGER 'point_id'"`
//	//WorkorderId     int64     `xorm:" INTEGER 'workorder_id'"`
//	MeasureDegree float32 `xorm:"numeric 'measure_degree'"`
//	MeasureTorque float32 `xorm:"numeric 'measure_torque'"`
//	MeasureResult string  `xorm:"varchar(32) 'measure_result'"`
//	PsetWMax      float32 `xorm:"numeric 'pset_w_max'"`
//	UserId        int64   `xorm:" INTEGER null 'user_id'"`
//	//ConsuProductId  int64     `xorm:" INTEGER 'consu_product_id'"`
//	WorkcenterId int64   `xorm:" INTEGER null 'workcenter_id'"`
//	Sent         int     `xorm:" bool 'sent'"`
//	PsetWTarget  float32 `xorm:"numeric 'pset_w_target'"`
//	//ProductionId    int64     `xorm:" INTEGER 'production_id'"`
//	Lacking      string `xorm:"varchar(32) 'lacking'"`
//	QualityState string `xorm:"varchar(32) 'quality_state'"`
//	PsetStrategy string `xorm:"varchar(32) 'pset_strategy'"`
//	//AssemblyLineId  int64     `xorm:" INTEGER 'assembly_line_id'"`
//	PsetMTarget     float32   `xorm:"numeric 'pset_m_target'"`
//	PsetWMin        float32   `xorm:"numeric 'pset_w_min'"`
//	ProductId       int64     `xorm:" INTEGER null 'product_id'"`
//	ControlDate     time.Time `xorm:"TIMESTAMP 'control_date'"`
//	Name            string    `xorm:"varchar(128) 'name'"`
//	PsetWThreshold  float32   `xorm:"numeric 'pset_w_threshold'"`
//	CurObjects      string    `xorm:"varchar(256) 'cur_objects'"`
//	PsetMMin        float32   `xorm:"numeric 'pset_m_min'"`
//	MeasureTDon     float32   `xorm:"numeric 'measure_t_don'"`
//	OpTime          int       `xorm:" INTEGER 'op_time'"`
//	ExceptionReason string    `xorm:"varchar(32) 'exception_reason'"`
//	GunID           int64     `xorm:" INTEGER null 'gun_id'"`
//	//ConsuBomLineID  int64	  `xorm:" INTEGER 'consu_bom_line_id'"`
//	Vin   string `xorm:"varchar(128) 'vin'"`
//	Batch string `xorm:"varchar(128) 'batch'"`
//}

var KEYS = []string{
	"pset_m_threshold",
	"pset_m_max",
	"control_date",
	"pset_w_max",
	"user_id",
	"one_time_pass",
	"pset_strategy",
	"pset_w_threshold",
	"cur_objects",
	"pset_m_target",
	"pset_m_min",
	"final_pass",
	"measure_degree",
	"measure_t_don",
	"measure_torque",
	"measure_result",
	"op_time",
	"pset_w_min",
	"pset_w_target",
	"lacking",
	"quality_state",
	"exception_reason",
	"id",
	"sent",
	"batch",
	"time",
	"workcenter_id",
	"gun_id",
	"product_id",
	"consu_product_id",
	"vin",
	"job",
}

type OperationResultModel struct {
	Id              int64     `gorm:"column:id;PRIMARY_KEY"`
	CreateTime      time.Time `gorm:"column:time"`
	PsetMThreshold  float64   `gorm:"column:pset_m_threshold"`
	PsetMMax        float64   `gorm:"column:pset_m_max"`
	FinalPass       string    `gorm:"column:final_pass"`
	OneTimePass     string    `gorm:"column:one_time_pass"`
	MeasureDegree   float64   `gorm:"column:measure_degree"`
	MeasureTorque   float64   `gorm:"column:measure_torque"`
	MeasureResult   string    `gorm:"column:measure_result"`
	PsetWMax        float64   `gorm:"column:pset_w_max"`
	UserId          int64     `gorm:"column:user_id"`
	WorkcenterId    int64     `gorm:"column:workcenter_id"`
	Sent            int       `gorm:"column:sent"`
	PsetWTarget     float64   `gorm:"column:pset_w_target"`
	Lacking         string    `gorm:"column:lacking"`
	QualityState    string    `gorm:"column:quality_state"`
	PsetStrategy    string    `gorm:"column:pset_strategy"`
	PsetMTarget     float64   `gorm:"column:pset_m_target"`
	PsetWMin        float64   `gorm:"column:pset_w_min"`
	ProductId       int64     `gorm:"column:product_id"`
	ControlDate     time.Time `gorm:"column:control_date"`
	Name            string    `gorm:"column:name"`
	PsetWThreshold  float64   `gorm:"column:pset_w_threshold"`
	CurObjects      string    `gorm:"column:cur_objects"`
	PsetMMin        float64   `gorm:"column:pset_m_min"`
	MeasureTDon     float64   `gorm:"column:measure_t_don"`
	OpTime          int       `gorm:"column:op_time"`
	ExceptionReason string    `gorm:"column:exception_reason"`
	GunID           int64     `gorm:"column:gun_id"`
	Vin             string    `gorm:"column:vin"`
	Batch           string    `gorm:"column:batch"`
}

type ResultObject struct {
	OR     map[string]interface{}
	ID     int64
	Send   int
	Port   string
	IP     string
	Stream *aiis.RPCAiis_RPCNodeServer
}

type OperationResult struct {
	ID              int64       `json:"id"`
	PsetMThreshold  float64     `json:"pset_m_threshold"`
	PsetMMax        float64     `json:"pset_m_max"`
	ControlDate     time.Time   `json:"control_date"`
	PsetWMax        float64     `json:"pset_w_max"`
	UserId          int64       `json:"user_id"`
	OneTimePass     string      `json:"one_time_pass"`
	PsetStrategy    string      `json:"pset_strategy"`
	PsetWThreshold  float64     `json:"pset_w_threshold"`
	CurObjects      []CURObject `json:"cur_objects"`
	PsetMTarget     float64     `json:"pset_m_target"`
	PsetMMin        float64     `json:"pset_m_min"`
	FinalPass       string      `json:"final_pass"`
	MeasureDegree   float64     `json:"measure_degree"`
	MeasureTDone    float64     `json:"measure_t_don"` //操作所用的时间
	MeasureTorque   float64     `json:"measure_torque"`
	MeasureResult   string      `json:"measure_result"`
	OPTime          int         `json:"op_time"`
	PsetWMin        float64     `json:"pset_w_min"`
	PsetWTarget     float64     `json:"pset_w_target"`
	QualityState    string      `json:"quality_state"`
	ExceptionReason string      `json:"exception_reason"`
	Seq             int         `json:"seq"`
	ProductID       int64       `json:"product_id"`
	WorkcenterID    int64       `json:"workcenter_id"`
	Vin             string      `json:"vin"`
	Job             string      `json:"job"`
	GunID           int64       `json:"gun_id"`
	Batch           string      `json:"batch"`
	Mode            string      `json:"mode"`
	ControllerSN    string      `json:"controller_sn"`
	TighteningId    int64       `json:"tightening_id"`
	ToolSN          string      `json:"tool_sn"`
	WorkcenterCode  string      `json:"workcenter_code"`
	NutID int64 `json:"consu_product_id"`

	// mo相关信息
	EquipemntName  string `json:"equipment_name"` // 设备名
	FactoryName    string `json:"factory_name"`   // 工厂代码
	Year           int64  `json:"year"`
	Pin            int64  `json:"pin"`
	Pin_check_code int64  `json:"pin_check_code"`
	AssemblyLine   string `json:"assembly_line"`
	Lnr            string `json:"lnr"`
	NutNo          string `json:"nut_no"`
	Model          string `json:"model"`

	Lacking string `json:"lacking"`
}

type CURObject struct {
	File string `json:"file"`
	OP   int    `json:"op"`
}
