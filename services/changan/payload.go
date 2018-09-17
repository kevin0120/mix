package changan

const (
	MSG_TASK      = 1001
	MSG_TASK_ACK  = 2001
	MSG_HEART     = 1999
	MSG_HEART_ACK = 2999
)

type AndonMsg struct {
	MsgType int         `json:"msg_type"`
	Seq     int         `json:"seq_no"`
	Data    interface{} `json:"data"`
}

type AndonTask struct {
	Workcenter string `json:"workcenter_code"`
	Vin        string `json:"vin_code"`
	Cartype    string `json:"cartype_code"`
}
