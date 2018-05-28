package rushdb

import "time"

type Workorders struct {
	Workorder_id	int			`xorm:"int 'workorder_id'"`
	HMI_sn			string		`xorm:"varchar(64) 'hmi_sn'"`
	PSet			int			`xorm:"int 'pset'"`
	Nut_total		float64		`xorm:"double 'nut_total'"`
	Vin				string		`xorm:"varchar(64) 'vin'"`
	Knr				string		`xorm:"varchar(64) 'knr'"`
	Max_redo_times	int			`xorm:"int 'max_redo_times'"`
	Max_op_time		int			`xorm:"int 'max_op_time'"`
	Status			string		`xorm:"varchar(32) 'status'"`
	Result_ids		string		`xorm:"text 'result_ids'"`
	WorkSheet		string		`xorm:"text 'work_sheet'"`
}

type Results struct {
	Result_id		int			`xorm:"int 'result_id'"`
	Controller_sn	string		`xorm:"varchar(64) 'controller_sn'"`
	Workorder_id	int			`xorm:"int 'workorder_id'"`
	Result			string		`xorm:"varchar(32) 'result'"`
	Cur_upload		bool		`xorm:"bool 'cur_upload'"`
	Result_upload	bool		`xorm:"bool 'result_upload'"`
	Need_upload 	bool		`xorm:"bool 'need_upload'"`
	Update_time		time.Time	`xorm:"date 'update_time'"`
	Result_data		string		`xorm:"text 'result_data'"`
	Cur_data		string		`xorm:"text 'cur_data'"`
	Total_count		int			`xorm:"int 'total_count'"`
	Count			int			`xorm:"int 'count'"`
}
