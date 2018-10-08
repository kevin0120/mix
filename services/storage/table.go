package storage

import "time"

type Workorders struct {
	Id          int64  `xorm:"pk autoincr notnull 'id'"`
	WorkorderID int64  `xorm:"bigint 'x_workorder_id'"`
	HMISN       string `xorm:"varchar(64) 'hmi_sn'"`
	Vin     string `xorm:"varchar(64) 'vin'"`
	Knr     string `xorm:"varchar(64) 'knr'"`
	LongPin string `xorm:"varchar(64) 'long_pin'"`

	MaxOpTime      int       `xorm:"int 'max_op_time'"`
	MaxSeq         int       `xorm:"int 'max_seq'"`
	Status         string    `xorm:"varchar(32) 'status'"`
	LastResultID   int64     `xorm:"bigint 'last_result_id'"`
	WorkSheet      string    `xorm:"text 'work_sheet'"`
	VehicleTypeImg string    `xorm:"text 'vehicle_type_img'"`
	UpdateTime     time.Time `xorm:"datetime 'update_time'"`
	ProductID      int64     `xorm:"bigint 'product_id'"`
	WorkcenterID   int64     `xorm:"bigint 'workcenter_id'"`
	UserID         int64     `xorm:"bigint 'user_id'"`

	JobID int    `xorm:"bigint 'job_id'"`
	Mode  string `xorm:"varchar(64) 'mode'"`

	// mo相关信息
	MO_EquipemntName  string `xorm:"varchar(64) 'equipment_name'"` // 设备名
	MO_FactoryName    string `xorm:"varchar(64) 'factory_name'"`   // 工厂代码
	MO_Year           int64  `xorm:"bigint 'year'"`
	MO_Pin            int64  `xorm:"bigint 'pin'"`
	MO_Pin_check_code int64  `xorm:"bigint 'pin_check_code'"`
	MO_AssemblyLine   string `xorm:"varchar(64) 'assembly_line'"`
	MO_Lnr            string `xorm:"varchar(64) 'lnr'"`
	MO_Model          string `xorm:"varchar(64) 'model'"`
}

type Results struct {
	Id                 int64     `xorm:"pk autoincr notnull 'id'"`
	Seq                int       `xorm:"int 'seq'"`
	GroupSeq           int       `xorm:"int 'group_sequence'"`
	ResultId           int64     `xorm:"bigint 'x_result_id'"`
	WorkorderID        int64     `xorm:"bigint 'x_workorder_id'"`
	UserID             int64     `xorm:"bigint 'user_id'"`
	ControllerSN       string    `xorm:"varchar(64) 'controller_sn'"`
	GunSN              string    `xorm:"varchar(64) 'gun_sn'"`
	Result             string    `xorm:"varchar(32) 'result'"`
	HasUpload          bool      `xorm:"bool 'has_upload'"`
	Stage              string    `xorm:"varchar(32) 'stage'"`
	UpdateTime         time.Time `xorm:"datetime 'update_time'"`
	PSetDefine         string    `xorm:"text 'pset_define'"`
	ResultValue        string    `xorm:"text 'result_value'"`
	Count              int       `xorm:"int 'count'"`
	PSet               int       `xorm:"int 'pset'"`
	NutNo              string    `xorm:"varchar(64) 'nut_no'"`
	ToleranceMinDegree float64   `xorm:"Double 'tolerance_min_degree'"`
	ToleranceMaxDegree float64   `xorm:"Double 'tolerance_max_degree'"`
	ToleranceMax       float64   `xorm:"Double 'tolerance_max'"`
	ToleranceMin       float64   `xorm:"Double 'tolerance_min'"`
	OffsetX            float64   `xorm:"Double 'offset_x'"`
	OffsetY            float64   `xorm:"Double 'offset_y'"`
	MaxRedoTimes       int       `xorm:"int 'max_redo_times'"`
	Batch              string    `xorm:"varchar(32) 'batch'"`
	ExInfo             string    `xorm:"text 'exinfo'"`
	Spent              int64     `xorm:"bigint 'spent'"`
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
	Id     int64  `xorm:"pk autoincr notnull 'id'"`
	GunID  int64  `xorm:"bigint 'gun_id'"`
	Serial string `xorm:"varchar(128) 'serial'"`
}

type RoutingOperations struct {
	OperationID    int64  `xorm:"bigint 'operation_id'"`
	Job            int    `xorm:"bigint 'job'"`
	MaxOpTime      int    `xorm:"int 'max_op_time'"`
	Name           string `xorm:"varchar(256) 'name'"`
	Img            string `xorm:"text 'img'"`
	ProductId      int64  `xorm:"bigint 'product_id'"`
	ProductType    string `xorm:"varchar(256) 'product_type'"`
	WorkcenterCode string `xorm:"varchar(256) 'workcenter_code'"`
	VehicleTypeImg string `xorm:"text 'vehicle_type_img'"`
	Points         string `xorm:"text 'points'"`
}
