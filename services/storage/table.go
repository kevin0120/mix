package storage

import "time"

type Workorders struct {
	Id           int64     `xorm:"pk autoincr notnull 'id'"`
	WorkorderID  int64     `xorm:"bigint 'x_workorder_id'"`
	HMISN        string    `xorm:"varchar(64) 'hmi_sn'"`
	//PSet         int       `xorm:"int 'pset'"`
	//NutTotal     float64   `xorm:"double 'nut_total'"`
	Vin          string    `xorm:"varchar(64) 'vin'"`
	Knr          string    `xorm:"varchar(64) 'knr'"`
	LongPin      string    `xorm:"varchar(64) 'long_pin'"`

	MaxOpTime    int       `xorm:"int 'max_op_time'"`
	Status       string    `xorm:"varchar(32) 'status'"`
	LastResultID int64     `xorm:"bigint 'last_result_id'"`
	WorkSheet    string    `xorm:"text 'work_sheet'"`
	UpdateTime   time.Time `xorm:"datetime 'update_time'"`

	// mo相关信息
	MO_EquipemntName  string `xorm:"varchar(64) 'equipment_name'"` // 设备名
	MO_FactoryName    string `xorm:"varchar(64) 'factory_name'"`   // 工厂代码
	MO_Year           int64  `xorm:"bigint 'year'"`
	MO_Pin            int64  `xorm:"bigint 'pin'"`
	MO_Pin_check_code int64  `xorm:"bigint 'pin_check_code'"`
	MO_AssemblyLine   string `xorm:"varchar(64) 'assembly_line'"`
	MO_Lnr            string `xorm:"varchar(64) 'lnr'"`
	MO_Model		  string `xorm:"varchar(64) 'model'"`

}

type Results struct {
	Id           int64     `xorm:"pk autoincr notnull 'id'"`
	Seq			 int	   `xorm:"int 'seq'"`
	ResultId     int64     `xorm:"bigint 'x_result_id'"`
	WorkorderID  int64     `xorm:"bigint 'x_workorder_id'"`
	UserID       int64     `xorm:"bigint 'user_id'"`
	ControllerSN string    `xorm:"varchar(64) 'controller_sn'"`
	GunSN 		 string    `xorm:"varchar(64) 'gun_sn'"`
	Result       string    `xorm:"varchar(32) 'result'"`
	HasUpload    bool      `xorm:"bool 'has_upload'"`
	Stage        string    `xorm:"varchar(32) 'stage'"`
	UpdateTime   time.Time `xorm:"datetime 'update_time'"`
	PSetDefine   string    `xorm:"text 'pset_define'"`
	ResultValue  string    `xorm:"text 'result_value'"`
	Count        int       `xorm:"int 'count'"`
	PSet         int       `xorm:"int 'pset'"`
	NutNo     	 string    `xorm:"varchar(64) 'nut_no'"`
	ToleranceMinDegree 		float64 `xorm:"Double 'tolerance_min_degree'"`
	ToleranceMaxDegree 		float64 `xorm:"Double 'tolerance_max_degree'"`
	ToleranceMax 			float64 `xorm:"Double 'tolerance_max'"`
	ToleranceMin 			float64 `xorm:"Double 'tolerance_min'"`
	OffsetX		 float64	`xorm:"Double 'offset_x'"`
	OffsetY		 float64	`xorm:"Double 'offset_y'"`
	MaxRedoTimes int        `xorm:"int 'max_redo_times'"`
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
