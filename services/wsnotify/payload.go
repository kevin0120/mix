package wsnotify

const (
	WsReg = "WS_REG"
)

type WSRegist struct {
	HMISn string `json:"hmi_sn"`
}

type WSMsg struct {
	SN   uint64      `json:"sn"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WSReply struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}
