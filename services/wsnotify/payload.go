package wsnotify

const (
	WS_REG       = "WS_REG"
	WS_RUSH_DATA = "WS_RUSH_DATA"
)

type WSRegist struct {
	HMI_SN string `json:"hmi_sn"`
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
