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

type RoutingOperationDelete struct {
	OperationID int64  `json:"id"`
	ProductType string `json:"product_type"`
}
