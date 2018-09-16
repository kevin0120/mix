package changan

import "encoding/json"

const (
	MSG_TASK      = 1001
	MSG_TASK_ACK  = 2001
	MSG_HEART     = 1999
	MSG_HEART_ACK = 2999
	MSG_GET_TASK  = 1002
	MSG_GET_TASK_ACK = 2002
	MSG_REGIST = 1003
	MSG_REGIST_ACK = 2003
	MSG_GUID_REQ = 1004
	MSG_GUID_REQ_ACK = 2004

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

type AndonKeepAlive struct {
}

type AndonGUID struct {
	GUID string `json:"topServerGUID"`
}

type AndonWorkCenter struct {
	Workcenter string `json:"workcenter_code"`
}

func PakcageMsg(msgType int, seq int, info interface{}) []byte {
	switch msgType {
	case MSG_HEART:
		d := AndonMsg{
			MsgType: MSG_HEART,
			Seq:     seq,
			Data:    nil,
		}
		ret, _ := json.Marshal(d)
		return ret
	case MSG_GET_TASK:
		d := AndonMsg{
			MsgType: MSG_GET_TASK,
			Seq:     seq,
			Data:    info,
		}
		ret, _ := json.Marshal(d)
		return ret
	case MSG_REGIST:
		d := AndonMsg{
			MsgType: MSG_REGIST,
			Seq:     seq,
			Data:    nil,
		}
		ret, _ := json.Marshal(d)
		return ret
	case MSG_GUID_REQ_ACK:
		d := AndonMsg{
			MsgType: MSG_GUID_REQ_ACK,
			Seq:     seq,
			Data:    info,
		}
		ret, _ := json.Marshal(d)
		return ret
	}
	return nil
}
