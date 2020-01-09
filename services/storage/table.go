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

	MaxOpTime    int   `xorm:"int 'max_op_time'" json:"-"`
	MaxSeq       int   `xorm:"int 'max_seq'" json:"-"`
	LastResultID int64 `xorm:"bigint 'last_result_id'" json:"-"`
	//WorkSheet      string    `xorm:"text 'work_sheet'"`
	ImageOPID      int64     `xorm:"bigint 'img_op_id'" json:"-"`
	VehicleTypeImg string    `xorm:"text 'vehicle_type_img'" json:"-"`
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

	Name           string      `xorm:"varchar(64) 'name'" json:"-"`
	Desc           string      `xorm:"varchar(64) 'desc'" json:"-"`
	Payload        string      `xorm:"text" json:"-"`
	MarshalPayload interface{} `xorm:"-" json:"-"`

	Workorder    string `xorm:"text 'workorder'" json:"-"`
	Code         string `xorm:"pk unique varchar(128) 'code'"  json:"code"`
	Track_code   string `xorm:"varchar(128) 'track_code'" json:"track_code"`
	Product_code string `xorm:"varchar(128) 'product_code'" json:"product_code"  validator:"required"`
	//Workcenter            string    `xorm:"varchar(128) 'workcenter'" json:"workcenter"`
	Date_planned_start    time.Time `xorm:"datetime 'date_planned_start'" json:"date_planned_start"`
	Date_planned_complete time.Time `xorm:"datetime 'date_planned_complete'" json:"date_planned_complete"`
	Status                string    `xorm:"varchar(32) default 'todo' 'status' " json:"status"`
	Product_type_image    string    `json:"product_type_image"`

	Unique_Num int64     `xorm:"bigint 'unique_num'" json:"unique_num"`
	Data       string    `xorm:"text 'data'" json:"data"`
	Created    time.Time `xorm:"created" json:"-"`
	Updated    time.Time `xorm:"updated" json:"-"`

	Steps []Steps `json:"-" validator:"required"`
}

type Steps struct {
	//Id             int64       `xorm:"pk autoincr notnull 'id'" json:"id"`
	////CurrentWorkorderID    int64     `xorm:"bigint 'workorder_id'" json:"-"`
	//Name           string      `xorm:"varchar(64) 'name'" json:"name"`
	//Desc           string      `xorm:"varchar(64) 'desc'" json:"desc"`
	//Type           string      `xorm:"varchar(64) 'type'" json:"type"`
	//Skippable      bool        `xorm:"varchar(64) 'skippable'" json:"skippable"`
	//Undoable       bool        `xorm:"varchar(64) 'undoable'" json:"undoable"`
	//Status         string      `xorm:"varchar(32) 'status'"   json:"status"`
	//Payload        string      `xorm:"text" json:"-"`
	//MarshalPayload interface{} `xorm:"-" json:"payload"`

	Id int64 `xorm:"pk autoincr notnull 'id'" json:"id"`
	//CurrentWorkorderID    int64       `xorm:"bigint 'workorder_id'" json:"-"`
	Name string `xorm:"varchar(64) 'name'" json:"-"`

	Type string `xorm:"varchar(64) 'type'" json:"-"`

	Payload        string      `xorm:"text" json:"payload"`
	MarshalPayload interface{} `xorm:"-" json:"-"`

	WorkorderID int64  `xorm:"bigint 'x_workorder_id'" json:"-"`
	Step        string `xorm:"text 'step'" json:"-"`

	Code           string `xorm:"varchar(128) 'code'" json:"code"`
	Sequence       int64  `xorm:"bigint 'sequence'" json:"sequence"`
	Testtype       string `xorm:"varchar(128) 'test_type'" json:"test_type"`
	FailureMessage string `xorm:"varchar(128) 'failure_msg'" json:"failure_msg"`
	Desc           string `xorm:"varchar(128) 'desc'" json:"desc"`
	Image          string `xorm:"-" json:"image"`
	ImageRef       string `xorm:"varchar(128) 'tightening_image_by_step_code'" json:"tightening_image_by_step_code" validator:"required"`
	Skippable      bool   `xorm:"varchar(64) 'skippable'" json:"skippable"`
	Undoable       bool   `xorm:"varchar(64) 'undoable'" json:"undoable"`
	Data           string `xorm:"text 'data'" json:"data"`
	Status         string `xorm:"varchar(32) default 'ready' 'status'" json:"status"`

	Created time.Time `xorm:"created" json:"-"`
	Updated time.Time `xorm:"updated" json:"-"`
}

type Results struct {
	Id                 int64     `xorm:"pk autoincr notnull 'id'"`
	HasUpload          bool      `xorm:"bool 'has_upload'"`
	Seq                int       `xorm:"int 'seq'"`
	GroupSeq           int       `xorm:"int 'group_sequence'"`
	ResultId           int64     `xorm:"bigint 'x_result_id'"`
	WorkorderID        int64     `xorm:"bigint 'x_workorder_id'"`
	StepID             int64     `xorm:"bigint 'step_id'"`
	UserID             int64     `xorm:"bigint 'user_id'"`
	ControllerSN       string    `xorm:"varchar(64) 'controller_sn'"`
	ToolSN             string    `xorm:"varchar(64) 'gun_sn'"`
	Result             string    `xorm:"varchar(32) 'result'"`
	Stage              string    `xorm:"varchar(32) 'stage'"`
	UpdateTime         time.Time `xorm:"datetime 'update_time'"`
	PSetDefine         string    `xorm:"text 'pset_define'"`
	ResultValue        string    `xorm:"text 'result_value'"`
	Count              int       `xorm:"int 'count'"`
	PSet               int       `xorm:"int 'pset'"`
	NutNo              string    `xorm:"varchar(64) 'nut_no'"`
	ConsuProductID     int64     `xorm:"bigint 'consu_product_id'"`
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
	TighteningID       string    `xorm:"varchar(128) 'tightening_id'"`
	CurveFile          string    `xorm:"varchar(256) 'curve_file'"`
	// payload字段
}

type ResultsWorkorders struct {
	Results    `xorm:"extends"`
	Workorders `xorm:"extends"`
}

type Curves struct {
	Id           int64     `xorm:"pk autoincr notnull 'id'"`
	ResultID     int64     `xorm:"bigint 'result_id'"`
	Count        int       `xorm:"int 'count'"`
	CurveFile    string    `xorm:"varchar(128) 'curve_file'"`
	CurveData    string    `xorm:"text 'curve_data'"`
	HasUpload    bool      `xorm:"bool 'has_upload'"`
	UpdateTime   time.Time `xorm:"datetime 'update_time'"`
	ToolSN       string    `xorm:"varchar(128) 'tool_sn'"`
	TighteningID string    `xorm:"varchar(128) 'tightening_id'"`
}

type Controllers struct {
	Id           int64     `xorm:"pk autoincr notnull 'id'"`
	SN           string    `xorm:"varchar(128) 'controller_sn'"`
	LastID       string    `xorm:"varchar(128) 'last_id'"`
	TriggerStart time.Time `xorm:"datetime 'trigger_start'"`
	TriggerStop  time.Time `xorm:"datetime 'trigger_stop'"`
}

type Tools struct {
	Id                 int64  `xorm:"pk autoincr notnull 'id'"`
	GunID              int64  `xorm:"bigint 'gun_id'"`
	Serial             string `xorm:"varchar(128) 'serial'"`
	CurrentWorkorderID int64  `xorm:"bigint 'workorder_id'"` //当前正在进行的工单
	Seq                int    `xorm:"bigint 'sequence'"`
	Count              int    `xorm:"int 'count'"`
	Total              int    `xorm:"int 'total'"`
	Mode               string `xorm:"varchar(128) 'mode'"`
	UserID             int64  `xorm:"bigint 'user_id'"`
	StepID             int64  `xorm:"bigint 'step_id'"`
}

type RoutingOperations struct {
	Id                int64  `xorm:"pk autoincr notnull 'id'"`
	OperationID       int64  `xorm:"bigint 'operation_id'"`
	Job               int    `xorm:"bigint 'job'"`
	MaxOpTime         int    `xorm:"int 'max_op_time'"`
	Name              string `xorm:"varchar(256) 'name'"`
	Img               string `xorm:"text 'img'"`
	TighteningStepRef string `xorm:"varchar(256) 'tightening_step_ref'"`

	ProductId    int64 `xorm:"bigint 'product_id'"`
	WorkcenterID int64 `xorm:"bigint 'workcenter_id'"`

	ProductType      string `xorm:"varchar(256) 'product_type'"`
	ProductTypeImage string `xorm:"text 'product_type_image'"`

	WorkcenterCode string `xorm:"varchar(256) 'workcenter_code'"`
	VehicleTypeImg string `xorm:"text 'vehicle_type_img'"`
	Points         string `xorm:"text 'points'"`
}
