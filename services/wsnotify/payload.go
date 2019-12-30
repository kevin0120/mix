package wsnotify

import "github.com/kataras/iris/websocket"

const (
	WS_REG       = "WS_REG"
	WS_RUSH_DATA = "WS_RUSH_DATA"
)

type WSRegist struct {
	HMI_SN string `json:"hmi_sn"`
}

type WSStatus struct {
	SN     string `json:"controller_sn"`
	Status string `json:"status"`
}

type WSMsgPackage struct {
	SN      string
	Event   string
	Payload string
}

type WSResult struct {
	//Result_id int64 `json:"result_id"`
	Seq      int     `json:"sequence"`
	GroupSeq int     `json:"group_sequence"`
	Count    int     `json:"count"`
	Result   string  `json:"result"`
	MI       float64 `json:"mi"`
	WI       float64 `json:"wi"`
	TI       float64 `json:"ti"`
	Batch    string  `json:"batch"`
	ToolSN   string  `json:"tool_sn"`
}

type WSSelector struct {
	SN        string `json:"controller_sn"`
	Selectors []int  `json:"selectors"`
}

type WSJobSelect struct {
	JobID int `json:"job_id"`
}

type WSRegistMsg struct {
	Msg string `json:"msg"`
}

type WSOdooStatus struct {
	Status string `json:"status"`
}

type WSMsg struct {
	SN   uint64      `json:"sn"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WSRequest struct {
	C     websocket.Connection
	WSMsg *WSMsg
}

type WSReply struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}
