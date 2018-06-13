package hmi

type PSet struct {
	Controller_SN string `json:"controller_sn"`
	PSet          int    `json:"pset"`
	Result_id     int64  `json:"result_id"`
	Count         int    `json:"count"`
	UserID		  int64	 `json:"user_id"`
}

type Worksheet struct {
	Content string `json:"content"`
	Points  []struct {
		X int32 `json:"x_offset"`
		Y int32 `json:"y_offset"`
	} `json:"points"`
}

type Workorder struct {
	Workorder_id int64     `json:"workorder_id"`
	HMI_sn       string    `json:"hmi_sn"`
	PSet         int       `json:"pset"`
	Nut_total    float64   `json:"nut_total"`
	Vin          string    `json:"vin"`
	Knr          string    `json:"knr"`
	LongPin		 string		`json:"long_pin"`
	Status       string    `json:"status"`
	Result_ids   []int64   `json:"result_ids"`
	WorkSheet    Worksheet `json:"work_sheet"`
	MaxRedoTimes int       `json:"max_redo_times"`
	MaxOpTime    int       `json:"max_op_time"`
}
