package aiis

const (
	TYPE_RESULT = "result_patch"
	TYPE_ODOO_STATUS = "odoo_status"
)


type RPCPayload struct {
	Type string `json:"type"`
	Data interface{}`json:"data"`
}