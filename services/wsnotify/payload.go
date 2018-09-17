package wsnotify

type WSRegist struct {
	HMI_SN string `json:"hmi_sn"`
}

type WSStatus struct {
	SN     string `json:"controller_sn"`
	Status string `json:"status"`
}

type WSResult struct {
	Result_id int64   `json:"result_id"`
	//Seq       int     `json:"sequence"`
	GroupSeq       int     `json:"group_sequence"`
	Count     int     `json:"count"`
	Result    string  `json:"result"`
	MI        float64 `json:"mi"`
	WI        float64 `json:"wi"`
	TI        float64 `json:"ti"`
}

type WSSelector struct {
	SN        string `json:"controller_sn"`
	Selectors []int  `json:"selectors"`
}

type WSJobSelect struct {
	JobID int `json:"job_id"`
}

type WSScanner struct {
	Barcode string `json:"barcode"`
}

//type WSIO struct {
//	Barcode string `json:"barcode"`
//}

type WSRegistMsg struct {
	Msg string `json:"msg"`
}
