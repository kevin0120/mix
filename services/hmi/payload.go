package hmi

const (
	MODE_PSET = "pset"
	MODE_JOB = "job"
)

type PSet struct {
	Controller_SN string `json:"controller_sn"`
	GunSN         string `json:"gun_sn"`
	PSet          int    `json:"pset"`
	Result_id     int64  `json:"result_id"`
	Count         int    `json:"count"`
	UserID        int64  `json:"user_id"`
}

type Job struct {
	Controller_SN string `json:"controller_sn"`
	Job           int    `json:"job_id"`
	WorkorderiD   int64  `json:"workorder_id"`
	UserID        int64  `json:"user_id"`
}

type ControllerMode struct {
	Controller_SN string `json:"controller_sn"`
	Mode        string   `json:"mode"`
}

type PSetDetailRequest struct {
	Controller_SN string `json:"controller_sn"`
	PSet          int     `json:"pset"`
}

type Result struct {
	ID            int64   `json:"id"`
	Controller_SN string  `json:"controller_sn"`
	GunSN         string  `json:"gun_sn"`
	PSet          int     `json:"pset"`
	MaxRedoTimes  int     `json:"max_redo_times"`
	X             float64 `json:"offset_x"`
	Y             float64 `json:"offset_Y"`
}

type Workorder struct {
	Workorder_id int64  `json:"workorder_id"`
	HMI_sn       string `json:"hmi_sn"`
	//PSet         int       `json:"pset"`
	//Nut_total    float64   `json:"nut_total"`
	Vin     string `json:"vin"`
	Knr     string `json:"knr"`
	LongPin string `json:"long_pin"`
	Status  string `json:"status"`
	//Result_ids   []int64   `json:"result_ids"`
	WorkSheet string `json:"work_sheet"`

	MaxOpTime int      `json:"max_op_time"`
	Results   []Result `json:"results"`
}
