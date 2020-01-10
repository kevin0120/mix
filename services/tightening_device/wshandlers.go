package tightening_device

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/wsnotify"
)

func (s *Service) OnWS_TOOL_MODE_SELECT(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolModeSelect{}
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolModeSelect(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_ENABLE(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolControl{}
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolControl(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_JOB(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req JobSet
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolJobSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))

}

func (s *Service) OnWS_TOOL_PSET(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetSet
	_ = json.Unmarshal(byteData, &req)

	err := s.ToolPSetSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_PSET_BATCH(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetBatchSet
	_ = json.Unmarshal(byteData, &req)

	err := s.ToolPSetBatchSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_PSET_LIST(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	psetList, err := s.GetToolPSetList(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, psetList))
}

func (s *Service) OnWS_TOOL_PSET_DETAIL(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolPSet
	_ = json.Unmarshal(byteData, &req)

	psetDetail, err := s.GetToolPSetDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, psetDetail))
}

func (s *Service) OnWS_TOOL_JOB_LIST(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	jobList, err := s.GetToolJobList(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, jobList))
}

func (s *Service) OnWS_TOOL_JOB_DETAIL(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolJob
	_ = json.Unmarshal(byteData, &req)

	jobDetail, err := s.GetToolJobDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, jobDetail))
}
