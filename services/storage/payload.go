package storage

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
