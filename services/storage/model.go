package storage

import (
	"time"
)


//odoo中的operation_result数据库 模型
type OperationResultModel struct {
	PsetMThreshold  float32				`xorm:"numeric 'pset_m_threshold'"`
	PsetMMax 		float32   			`xorm:"numeric 'pset_m_max'"`
	ControlDate	    time.Time			`xorm:"TIMESTAMPZ 'control_date'"`
	PsetWMax 		float32   			`xorm:"numeric 'pset_w_max'"`
	UserId 			int64				`xorm:" INTEGER 'user_id'"`
	//Id   			int64				`xorm:"pk autoincr notnull 'id'"`


	FinalPass       string    			`xorm:"varchar(32) 'final_pass'"`
	OneTimePass    	string      		`xorm:"varchar(32) 'one_time_pass'"`
	Stage        	string    `xorm:"varchar(32) 'stage'"`
	UpdateTime   	time.Time `xorm:"datetime 'update_time'"`
	PSetDefine   	string    `xorm:"text 'pset_define'"`
	ResultValue  	string    `xorm:"text 'result_value'"`
	Count        	int       `xorm:"int 'count'"`
}
