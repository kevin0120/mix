package rush

import "time"

type CurObj struct {
	File     			string 			`json:"file"`
	Operation 		 	int  			`json:"op"`
}

type OperationResult struct {
	PsetMThreshold  float32				`json:"pset_m_threshold"`
	PsetMMax 		float32   			`json:"pset_m_max"`
	ControlDate	    time.Time			`json:"control_date"`
	PsetWMax 		float32   			`json:"pset_w_max"`
	UserId 			int64				`json:"user_id"`
	OneTimePass    	string      		`json:"one_time_pass"`
	PsetStrategy	string				`json:"pset_strategy"`
	PsetWThreshold  float32				`json:"pset_w_threshold"`
	CurObjects 		[]*CurObj			`json:"cur_objects"`
	PsetMTarget		float32				`json:"pset_m_target"`
	PsetMMin		float32				`json:"pset_m_min"`
	FinalPass       string    			`json:"final_pass"`
	MeasureDegree	float32				`json:"measure_degree"`
	MeasureTDone	float32				`json:"measure_t_don"` //操作所用的时间
	MeasureTorque	float32				`json:"measure_torque"`
	Id 				int64				`json:"id"`
	MeasureResult   string 				`json:"measure_result"`
	OPTime			int					`json:"op_time"`
	PsetWMin		float32				`json:"pset_w_min"`
	PsetWTarget		float32				`json:"pset_w_target"`
}

