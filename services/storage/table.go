package storage

import "time"

type Workorders struct {
	Id   			int64	`xorm:"pk autoincr notnull 'id'"`
	WorkorderID  	int64     `xorm:"bigint 'x_workorder_id'"`
	HMISN        	string  `xorm:"varchar(64) 'hmi_sn'"`
	PSet         	int     `xorm:"int 'pset'"`
	NutTotal     	float64 `xorm:"double 'nut_total'"`
	Vin          	string  `xorm:"varchar(64) 'vin'"`
	Knr          	string  `xorm:"varchar(64) 'knr'"`
	LongPin         string  `xorm:"varchar(64) 'long_pin'"`
	MaxRedoTimes 	int     `xorm:"int 'max_redo_times'"`
	MaxOpTime    	int     `xorm:"int 'max_op_time'"`
	Status       	string  `xorm:"varchar(32) 'status'"`
	ResultIDs    	string  `xorm:"text 'result_ids'"`
	WorkSheet    	string  `xorm:"text 'work_sheet'"`
}

type Results struct {
	Id   			int64		`xorm:"pk autoincr notnull 'id'"`
	ResultId     	int64       `xorm:"bigint 'x_result_id'"`
	WorkorderID  	int64       `xorm:"bigint 'x_workorder_id'"`
	UserID  		int64       `xorm:"bigint 'user_id'"`
	ControllerSN 	string    `xorm:"varchar(64) 'controller_sn'"`
	Result       	string    `xorm:"varchar(32) 'result'"`
	HasUpload    	bool      `xorm:"bool 'has_upload'"`
	Stage        	string    `xorm:"varchar(32) 'stage'"`
	UpdateTime   	time.Time `xorm:"datetime 'update_time'"`
	PSetDefine   	string    `xorm:"text 'pset_define'"`
	ResultValue  	string    `xorm:"text 'result_value'"`
	Count        	int       `xorm:"int 'count'"`
}

type Curves struct {
	Id   		int64	 `xorm:"pk autoincr notnull 'id'"`
	ResultID  	int64    `xorm:"bigint 'result_id'"`
	Count     	int    `xorm:"int 'count'"`
	CurveFile 	string `xorm:"varchar(128) 'curve_file'"`
	CurveData 	string `xorm:"text 'curve_data'"`
	HasUpload 	bool   `xorm:"bool 'has_upload'"`
}
