package hmi

import (
	"github.com/masami10/rush/services/openprotocol"
)

const (
	WS_ORDER_LIST             = "WS_ORDER_LIST"
	WS_ORDER_DETAIL           = "WS_ORDER_DETAIL"
	WS_ORDER_UPDATE           = "WS_ORDER_UPDATE"
	WS_ORDER_STEP_UPDATE      = "WS_ORDER_STEP_UPDATE"
	WS_ORDER_START_REQUEST    = "WS_ORDER_START_REQUEST"
	WS_ORDER_FINISH_REQUEST   = "WS_ORDER_FINISH_REQUEST"
	WS_ORDER_STEP_DATA_UPDATE = "WS_ORDER_STEP_DATA_UPDATE"
	WS_ORDER_DETAIL_BY_CODE   = "WS_ORDER_DETAIL_BY_CODE"
	WS_WORKORDER_DATA_UPDATE  = "WS_WORKORDER_DATA_UPDATE"
	WS_ORDER_NEW_ORDER        = "WS_NEW_ORDER"

	WS_AIIS_STATUS  = "WS_AIIS_STATUS"
	WS_ODOO_STATUS  = "WS_ODOO_STATUS"
	WS_EXSYS_STATUS = "WS_EXSYS_STATUS"
)

type PSet struct {
	ControllerSN string `json:"controller_sn" validate:"required"`
	GunSN        string `json:"gun_sn" validate:"required"`
	PSet         int    `json:"pset" validate:"required"`
	Result_id    int64  `json:"result_id" validate:"-"`
	Count        int    `json:"count" validate:"gte=1"` //必须大于等于1
	UserID       int64  `json:"user_id" validate:"-"`
	GroupSeq     int    `json:"group_sequence" validate:"gt=0"`
	WorkorderID  int64  `json:"workorder_id" validate:"gt=0"` //必须大于0
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
	ToolSN        string `json:"gun_sn"`
	Enable        bool   `json:"enable"`
	Reason        string `json:"reason"`
}

type Job struct {
	Controller_SN string `json:"controller_sn"`
	ToolSN        string `json:"tool_sn"`
	Job           int    `json:"job_id"`
	WorkorderID   int64  `json:"workorder_id"`
	UserID        int64  `json:"user_id"`
}

type JobManual struct {
	Controller_SN string `json:"controller_sn"`
	GunSN         string `json:"gun_sn"`
	Job           int    `json:"job_id"`
	UserID        int64  `json:"user_id"`
	CarType       string `json:"car_type"`
	Vin           string `json:"vin"`
	HmiSN         string `json:"hmi_sn"`
	OperationID   int64  `json:"operation_id"`
	Skip          bool   `json:"skip"`
	HasSet        bool   `json:"has_set"`
	Mode          string `json:"mode"`
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
	Y             float64 `json:"offset_y"`
	Seq           int     `json:"sequence"`
	GroupSeq      int     `json:"group_sequence"`
}

type Workorder struct {
	Workorder_id   int64  `json:"workorder_id"`
	HMI_sn         string `json:"hmi_sn"`
	Vin            string `json:"vin"`
	Knr            string `json:"knr"`
	LongPin        string `json:"long_pin"`
	Status         string `json:"status"`
	WorkSheet      string `json:"work_sheet"`
	VehicleTypeImg string `json:"vehicleTypeImg"`

	MaxOpTime int      `json:"max_op_time"`
	Job       int      `json:"job_id"`
	Lnr       string   `json:"lnr"`
	Model     string   `json:"model"`
	Results   []Result `json:"results"`

	Reasons []string `json:"reasons"`
}

type NextWorkorder struct {
	Vin     string `json:"vin"`
	Model   string `json:"model"`
	LongPin string `json:"long_pin"`
	Knr     string `json:"knr"`
	Lnr     string `json:"lnr"`
}

type JobControl struct {
	Controller_SN string `json:"controller_sn"`
	Action        string `json:"action"`
}

type LocalResults struct {
	HmiSN        interface{} `json:"hmi_sn,omitempty"`
	Vin          interface{} `json:"vin,omitempty"`
	VehicleType  interface{} `json:"vehicle_type,omitempty"`
	JobID        interface{} `json:"job_id,omitempty"`
	PSetID       interface{} `json:"pset_id,omitempty"`
	ControllerSN interface{} `json:"controller_sn,omitempty"`
	GunSN        interface{} `json:"gun_sn,omitempty"`
	Result       interface{} `json:"result,omitempty"`
	Torque       interface{} `json:"torque,omitempty"`
	Angle        interface{} `json:"angle,omitempty"`
	Spent        interface{} `json:"spent,omitempty"`
	TimeStamp    interface{} `json:"timestamp,omitempty"`
	Batch        interface{} `json:"batch,omitempty"`
}

type TestProtocol struct {
	ProtocolType string `json:"type"`
	Payload      string `json:"payload"`
}

type NewWorkorder struct {
	WorkorderID int64 `json:"workorder_id"`
}

type RoutingOperationDelete struct {
	OperationID int64  `json:"id"`
	ProductType string `json:"product_type"`
}

type WSTest struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WSOrderReq struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type WSOrderReqData struct {
	ID   int64  `json:"id"`
	Data string `json:"data"`
}

type WSOrderReqCode struct {
	Code       string `json:"code"`
	Workcenter string `json:"workcenter"`
}

type WSWorkcenter struct {
	WorkCenter string `json:"workcenter"`
}
