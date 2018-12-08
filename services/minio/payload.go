package minio

type ControllerCurve struct {
	ResultID     int64
	CurveFile    string
	CurveContent ControllerCurveFile
	Count        int
	UpdateTime   string
}

type ControllerCurveFile struct {
	Result string    `json:"result"`
	CUR_M  []float64 `json:"cur_m"`
	CUR_W  []float64 `json:"cur_w"`
	CUR_T  []float64 `json:"cur_t"`
}
