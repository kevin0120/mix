package storage

import (
	"time"
)

//odoo中的operation_result数据库 模型
type OperationResultModel struct {
	Id              int64     `xorm:"pk SERIAL notnull 'id'"`
	CreateTime      time.Time `xorm:"DATE notnull 'time'"`
	PsetMThreshold  float32   `xorm:"numeric 'pset_m_threshold'"`
	PsetMMax        float32   `xorm:"numeric 'pset_m_max'"`
	FinalPass       string    `xorm:"varchar(32) 'final_pass'"`
	OneTimePass     string    `xorm:"varchar(32) 'one_time_pass'"`
	PointId         int64     `xorm:" INTEGER 'point_id'"`
	WorkorderId     int64     `xorm:" INTEGER 'workorder_id'"`
	MeasureDegree   float32   `xorm:"numeric 'measure_degree'"`
	MeasureTorque   float32   `xorm:"numeric 'measure_torque'"`
	MeasureResult   string    `xorm:"varchar(32) 'measure_result'"`
	PsetWMax        float32   `xorm:"numeric 'pset_w_max'"`
	UserId          int64     `xorm:" INTEGER 'user_id'"`
	ConsuProductId  int64     `xorm:" INTEGER 'consu_product_id'"`
	WorkcenterId    int64     `xorm:" INTEGER 'workcenter_id'"`
	Sent            int       `xorm:" bool 'sent'"`
	PsetWTarget     float32   `xorm:"numeric 'pset_w_target'"`
	ProductionId    int64     `xorm:" INTEGER 'production_id'"`
	Lacking         string    `xorm:"varchar(32) 'lacking'"`
	QualityState    string    `xorm:"varchar(32) 'quality_state'"`
	PsetStrategy    string    `xorm:"varchar(32) 'pset_strategy'"`
	AssemblyLineId  int64     `xorm:" INTEGER 'assembly_line_id'"`
	PsetMTarget     float32   `xorm:"numeric 'pset_m_target'"`
	PsetWMin        float32   `xorm:"numeric 'pset_w_min'"`
	ProductId       int64     `xorm:" INTEGER 'product_id'"`
	ControlDate     time.Time `xorm:"TIMESTAMP 'control_date'"`
	Name            string    `xorm:"varchar(128) 'name'"`
	PsetWThreshold  float32   `xorm:"numeric 'pset_w_threshold'"`
	CurObjects      string    `xorm:"varchar(256) 'cur_objects'"`
	PsetMMin        float32   `xorm:"numeric 'pset_m_min'"`
	MeasureTDon     float32   `xorm:"numeric 'measure_t_don'"`
	OpTime          int       `xorm:" INTEGER 'op_time'"`
	ExceptionReason string    `xorm:"varchar(32) 'exception_reason'"`
}
