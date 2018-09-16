package hmi

import (
	"github.com/masami10/rush/services/openprotocol"
)

type PSet struct {
	Controller_SN string `json:"controller_sn"`
	GunSN         string `json:"gun_sn"`
	PSet          int    `json:"pset"`
	Result_id     int64  `json:"result_id"`
	Count         int    `json:"count"`
	UserID        int64  `json:"user_id"`
}

type PSetManual struct {
	Controller_SN string `json:"controller_sn"`
	GunSN         string `json:"gun_sn"`
	PSet          int    `json:"pset"`
	UserID        int64  `json:"user_id"`
	CarType       string `json:"car_type"`
	Vin           string `json:"vin"`
	HmiSN         string `json:"hmi_sn"`
	ProductID     int64  `json:"product_id"`
	Count         int    `json:"count"`
	WorkcenterID  int64  `json:"workcenter_id"`
}

type IOSet struct {
	Controller_SN string                  `json:"controller_sn"`
	IOStatus      []openprotocol.IOStatus `json:"io_status"`
}

type ToolEnable struct {
	Controller_SN string `json:"controller_sn"`
	GunSN         string `json:"gun_sn"`
	Enable        bool   `json:"enable"`
}

type Job struct {
	Controller_SN string `json:"controller_sn"`
	Job           int    `json:"job_id"`
	WorkorderiD   int64  `json:"workorder_id"`
	UserID        int64  `json:"user_id"`
}

type JobManual struct {
	Controller_SN string     `json:"controller_sn"`
	Job           int        `json:"job_id"`
	UserID        int64      `json:"user_id"`
	CarType       string     `json:"car_type"`
	Vin           string     `json:"vin"`
	HmiSN         string     `json:"hmi_sn"`
	ProductID     int64      `json:"product_id"`
	Points        []JobPoint `json:"points"`
	WorkcenterID  int64      `json:"workcenter_id"`
	Skip          bool       `json:"skip"`
}

type JobPoint struct {
	Seq       int     `json:"sequence"`
	PSet      int     `json:"pset"`
	X         float64 `json:"offset_x"`
	Y         float64 `json:"offset_Y"`
	MaxOpTime int     `json:"max_op_time"`
	NutID     int64   `json:"nut_id"`
}

type ControllerMode struct {
	Controller_SN string `json:"controller_sn"`
	Mode          string `json:"mode"`
}

type PSetDetailRequest struct {
	Controller_SN string `json:"controller_sn"`
	PSet          int    `json:"pset"`
}

type Result struct {
	ID            int64   `json:"id"`
	Controller_SN string  `json:"controller_sn"`
	GunSN         string  `json:"gun_sn"`
	PSet          int     `json:"pset"`
	MaxRedoTimes  int     `json:"max_redo_times"`
	X             float64 `json:"offset_x"`
	Y             float64 `json:"offset_Y"`
	Seq           int     `json:"sequence"`
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
	WorkSheet      string `json:"work_sheet"`
	VehicleTypeImg string `json:"vehicleTypeImg"`

	MaxOpTime int      `json:"max_op_time"`
	Job       int      `json:"job_id"`
	Results   []Result `json:"results"`
}

type JobControl struct {
	Controller_SN string `json:"controller_sn"`
	Action        string `json:"action"`
}
