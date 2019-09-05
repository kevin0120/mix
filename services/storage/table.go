package storage

import (
	"time"
)

type Workorders struct {
	Id             int64  `xorm:"pk autoincr notnull 'id'" json:"id"`
	WorkorderID    int64  `xorm:"bigint 'x_workorder_id'" json:"-"`
	HMISN          string `xorm:"varchar(64) 'hmi_sn'" json:"-"`
	WorkcenterCode string `xorm:"varchar(64) 'workcenter_code'" json:"-"`
	Vin            string `xorm:"varchar(64) 'vin'" json:"-"`
	Knr            string `xorm:"varchar(64) 'knr'" json:"-"`
	LongPin        string `xorm:"varchar(64) 'long_pin'" json:"-"`

	MaxOpTime    int    `xorm:"int 'max_op_time'" json:"-"`
	MaxSeq       int    `xorm:"int 'max_seq'" json:"-"`
	Status       string `xorm:"varchar(32) 'status'" json:"status"`
	LastResultID int64  `xorm:"bigint 'last_result_id'" json:"-"`
	//WorkSheet      string    `xorm:"text 'work_sheet'"`
	ImageOPID      int64     `xorm:"bigint 'img_op_id'" json:"-"`
	VehicleTypeImg string    `xorm:"text 'vehicle_type_img'" json:"product_type_img"`
	UpdateTime     time.Time `xorm:"datetime 'update_time'" json:"-"`
	ProductID      int64     `xorm:"bigint 'product_id'" json:"-"`
	WorkcenterID   int64     `xorm:"bigint 'workcenter_id'" json:"-"`
	UserID         int64     `xorm:"bigint 'user_id'" json:"-"`

	JobID int    `xorm:"bigint 'job_id'" json:"-"`
	Mode  string `xorm:"varchar(64) 'mode'" json:"-"`

	Consumes string `xorm:"text 'consumes'" json:"-"`

	// mo相关信息
	MO_EquipemntName  string `xorm:"varchar(64) 'equipment_name'" json:"-"` // 设备名
	MO_FactoryName    string `xorm:"varchar(64) 'factory_name'" json:"-"`   // 工厂代码
	MO_Year           int64  `xorm:"bigint 'year'" json:"-"`
	MO_Pin            int64  `xorm:"bigint 'pin'" json:"-"`
	MO_Pin_check_code int64  `xorm:"bigint 'pin_check_code'" json:"-"`
	MO_AssemblyLine   string `xorm:"varchar(64) 'assembly_line'" json:"-"`
	MO_Lnr            string `xorm:"varchar(64) 'lnr'" json:"-"`
	MO_Model          string `xorm:"varchar(64) 'model'" json:"-"`

	Name string `xorm:"varchar(64) 'name'" json:"name"`
	Desc string `xorm:"varchar(64) 'desc'" json:"desc"`

	Payload        string      `xorm:"text" json:"-"`
	MarshalPayload interface{} `xorm:"-" json:"payload"`
}

type Steps struct {
	Id             int64       `xorm:"pk autoincr notnull 'id'" json:"id"`
	WorkorderID    int64       `xorm:"bigint 'workorder_id'" json:"-"`
	Name           string      `xorm:"varchar(64) 'name'" json:"name"`
	Desc           string      `xorm:"varchar(64) 'desc'" json:"desc"`
	Type           string      `xorm:"varchar(64) 'type'" json:"type"`
	Skippable      bool        `xorm:"varchar(64) 'skippable'" json:"skippable"`
	Undoable       bool        `xorm:"varchar(64) 'undoable'" json:"undoable"`
	Status         string      `xorm:"varchar(32) 'status'" json:"status"`
	Payload        string      `xorm:"text" json:"-"`
	MarshalPayload interface{} `xorm:"-" json:"payload"`
}

type Results struct {
	Id        int64 `xorm:"pk autoincr notnull 'id'"`
	HasUpload bool  `xorm:"bool 'has_upload'"`
	//Seq                int       `xorm:"int 'seq'"`
	//GroupSeq           int       `xorm:"int 'group_sequence'"`
	//ResultId           int64     `xorm:"bigint 'x_result_id'"`
	//WorkorderID        int64     `xorm:"bigint 'x_workorder_id'"`
	//StepID             int64     `xorm:"bigint 'step_id'"`
	//UserID             int64     `xorm:"bigint 'user_id'"`
	//ControllerSN       string    `xorm:"varchar(64) 'controller_sn'"`
	//GunSN              string    `xorm:"varchar(64) 'gun_sn'"`
	//Result             string    `xorm:"varchar(32) 'result'"`
	//Stage              string    `xorm:"varchar(32) 'stage'"`
	//UpdateTime         time.Time `xorm:"datetime 'update_time'"`
	//PSetDefine         string    `xorm:"text 'pset_define'"`
	//ResultValue        string    `xorm:"text 'result_value'"`
	//Count              int       `xorm:"int 'count'"`
	//PSet               int       `xorm:"int 'pset'"`
	//NutNo              string    `xorm:"varchar(64) 'nut_no'"`
	//ConsuProductID     int64     `xorm:"bigint 'consu_product_id'"`
	//ToleranceMinDegree float64   `xorm:"Double 'tolerance_min_degree'"`
	//ToleranceMaxDegree float64   `xorm:"Double 'tolerance_max_degree'"`
	//ToleranceMax       float64   `xorm:"Double 'tolerance_max'"`
	//ToleranceMin       float64   `xorm:"Double 'tolerance_min'"`
	//OffsetX            float64   `xorm:"Double 'offset_x'"`
	//OffsetY            float64   `xorm:"Double 'offset_y'"`
	//MaxRedoTimes       int       `xorm:"int 'max_redo_times'"`
	//Batch              string    `xorm:"varchar(32) 'batch'"`
	//ExInfo             string    `xorm:"text 'exinfo'"`
	//Spent              int64     `xorm:"bigint 'spent'"`
	//TighteningID       string    `xorm:"varchar(128) 'tightening_id'"`

	// 控制器序列号
	ControllerSN string `json:"controller_sn"`

	// 工具序列号
	ToolSN string `json:"tool_sn"`

	// 收到时间
	UpdateTime time.Time `json:"update_time"`

	// job号
	Job int `json:"job"`

	// pset号
	PSet int `json:"pset"`

	// 批次信息
	Batch string `json:"batch"`

	// 当前拧紧次数
	Count int `json:"count"`

	// 拧紧ID
	TighteningID string `json:"tightening_id"`

	// 实际结果
	MeasureResult string `json:"measure_result"`

	// 实际扭矩
	MeasureTorque float64 `json:"measure_torque"`

	// 实际角度
	MeasureAngle float64 `json:"measure_angle"`

	// 实际耗时
	MeasureTime float64 `json:"measure_time"`

	// 拧紧策略
	Strategy string `json:"strategy"`

	// 最大扭矩
	TorqueMax float64 `json:"torque_max"`

	// 最小扭矩
	TorqueMin float64 `json:"torque_min-"`

	// 扭矩阈值
	TorqueThreshold float64 `json:"torque_threshold"`

	// 目标扭矩
	TorqueTarget float64 `json:"torque_target"`

	// 最大角度
	AngleMax float64 `json:"angle_max"`

	// 最小角度
	AngleMin float64 `json:"angle_min"`

	// 目标角度
	AngleTarget float64 `json:"angle_target"`

	// payload字段
}

type ResultsWorkorders struct {
	Results    `xorm:"extends"`
	Workorders `xorm:"extends"`
}

type Curves struct {
	Id         int64     `xorm:"pk autoincr notnull 'id'"`
	ResultID   int64     `xorm:"bigint 'result_id'"`
	Count      int       `xorm:"int 'count'"`
	CurveFile  string    `xorm:"varchar(128) 'curve_file'"`
	CurveData  string    `xorm:"text 'curve_data'"`
	HasUpload  bool      `xorm:"bool 'has_upload'"`
	UpdateTime time.Time `xorm:"datetime 'update_time'"`
}

type Controllers struct {
	Id           int64     `xorm:"pk autoincr notnull 'id'"`
	SN           string    `xorm:"varchar(128) 'controller_sn'"`
	LastID       string    `xorm:"varchar(128) 'last_id'"`
	TriggerStart time.Time `xorm:"datetime 'trigger_start'"`
	TriggerStop  time.Time `xorm:"datetime 'trigger_stop'"`
}

type Guns struct {
	Id          int64  `xorm:"pk autoincr notnull 'id'"`
	GunID       int64  `xorm:"bigint 'gun_id'"`
	Serial      string `xorm:"varchar(128) 'serial'"`
	WorkorderID int64  `xorm:"bigint 'workorder_id'"`
	Seq         int    `xorm:"bigint 'sequence'"`
	Count       int    `xorm:"int 'count'"`
	Mode        string `xorm:"varchar(128) 'mode'"`
	Trace       string `xorm:"text 'trace'"`
}

type RoutingOperations struct {
	OperationID    int64  `xorm:"bigint 'operation_id'"`
	Job            int    `xorm:"bigint 'job'"`
	MaxOpTime      int    `xorm:"int 'max_op_time'"`
	Name           string `xorm:"varchar(256) 'name'"`
	Img            string `xorm:"text 'img'"`
	ProductId      int64  `xorm:"bigint 'product_id'"`
	WorkcenterID   int64  `xorm:"bigint 'workcenter_id'"`
	ProductType    string `xorm:"varchar(256) 'product_type'"`
	WorkcenterCode string `xorm:"varchar(256) 'workcenter_code'"`
	VehicleTypeImg string `xorm:"text 'vehicle_type_img'"`
	Points         string `xorm:"text 'points'"`
}
