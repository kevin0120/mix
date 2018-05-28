package rushdb

import "time"

type Workorders struct {
	WorkorderID		int			`xorm:"int 'workorder_id'"`
	HMISN			string		`xorm:"varchar(64) 'hmi_sn'"`
	PSet			int			`xorm:"int 'pset'"`
	NutTotal		float64		`xorm:"double 'nut_total'"`
	Vin				string		`xorm:"varchar(64) 'vin'"`
	Knr				string		`xorm:"varchar(64) 'knr'"`
	MaxRedoTimes	int			`xorm:"int 'max_redo_times'"`
	MaxOpTime		int			`xorm:"int 'max_op_time'"`
	Status			string		`xorm:"varchar(32) 'status'"`
	ResultIDs		string		`xorm:"text 'result_ids'"`
	WorkSheet		string		`xorm:"text 'work_sheet'"`
}

type Results struct {
	ResultId		int			`xorm:"int 'result_id'"`
	WorkorderID		int			`xorm:"int 'workorder_id'"`
	ControllerSN	string		`xorm:"varchar(64) 'controller_sn'"`
	Result			string		`xorm:"varchar(32) 'result'"`
	HasUpload		bool		`xorm:"bool 'has_upload'"`
	Stage 			string		`xorm:"varchar(32) 'stage'"`
	UpdateTime		time.Time	`xorm:"date 'update_time'"`
	PSetDefine		string		`xorm:"text 'pset_define'"`
	ResultValue		string		`xorm:"text 'result_value'"`
	Count			int			`xorm:"int 'count'"`
}

type Curves struct {
	ResultID		int			`xorm:"int 'result_id'"`
	Count			int			`xorm:"int 'count'"`
	CurveFile		string		`xorm:"varchar(32) 'curve_file'"`
	CurveData		string		`xorm:"text 'curve_data'"`
	HasUpload		bool		`xorm:"bool 'has_upload'"`
}
