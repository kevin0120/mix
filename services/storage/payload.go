package storage

import "time"

const (
	RESULT_NONE = "NONE"
	RESULT_OK   = "OK"
	RESULT_NOK  = "NOK"
	RESULT_LSN  = "LSN"
	RESULT_AK2  = "AK2"
)

const (
	RESULT_STAGE_INIT  = "init"
	RESULT_STAGE_FINAL = "final"
)

const (
	WORKORDER_STATUS_TODO    = "todo"
	WORKORDER_STATUS_WIP     = "wip"
	WORKORDER_STATUS_DONE    = "done"
	WORKORDER_STATUS_CANCEL  = "cancel"
	WORKORDER_STATUS_PENDING = "pending"
)

const (
	STEP_STATUS_READY  = "ready"
	STEP_STATUS_DOING  = "doing"
	STEP_STATUS_FAIL   = "fail"
	STEP_STATUS_FINISH = "finish"
)

type RoutingOperationDelete struct {
	OperationID int64  `json:"id"`
	ProductType string `json:"product_type"`
}

type WorkorderStep struct {
	Workorders
	Steps []Steps `json:"steps"`
}

type WorkorderList struct {
	Id                    int64     `json:"id"`
	Code                  string    `json:"code"`
	Track_code            string    `json:"track_code"`
	Product_code          string    `json:"product_code"`
	WorkCenter            string    `json:"-"`
	Date_planned_start    time.Time `json:"date_planned_start"`
	Date_planned_complete time.Time `json:"date_planned_complete"`
	Status                string    `json:"status"`
	Product_type_image    string    `json:"product_type_image"`
}

type WorkorderListPar struct {
	Time_from string `json:"time_from"`
	Time_to   string `json:"time_to"`
	Status    string `json:"status"`
	Page_size int    `json:"page_size"`
	Page_no   int    `json:"page_no"`
}

type WorkorderPayload struct {
	Products     interface{}   `json:"products"`
	Worksheet    interface{}   `json:"worksheet"`
	Environments []interface{} `json:"environments"`
	Components   []interface{} `json:"components"`
	Operation    interface{}   `json:"operation"`
	Workcenter   interface{}   `json:"workcenter"`

	SYSTEMTYPE         interface{} `json:"SYSTEMTYPE"`
	WIPORDERTYPE       interface{} `json:"WIPORDERTYPE"`
	MOMDISPOSITIONS    interface{} `json:"MOMDISPOSITIONS"`
	MOMCONFIG          interface{} `json:"MOMCONFIG"`
	RESOURCEGROUP      interface{} `json:"RESOURCEGROUP"`
	STARTEMPLOYEE      interface{} `json:"STARTEMPLOYEE"`
	RESOURCENAMES      interface{} `json:"RESOURCENAMES"`
	PARENTWIPORDERNO   interface{} `json:"PARENTWIPORDERNO"`
	PARENTWIPORDERTYPE interface{} `json:"PARENTWIPORDERTYPE"`
}

type WorkCenterPayload struct {
	Code string `json:"code"`
}

type StepTextPayload struct {
}

type StepTighteningPayload struct {
	TighteningTotal int           `json:"tightening_total"`
	TighteningPoint []interface{} `json:"tightening_points"`
}
