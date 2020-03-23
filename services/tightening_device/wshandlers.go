package tightening_device

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/wsnotify"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (s *Service) OnWS_TOOL_MODE_SELECT(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolModeSelect{}
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolModeSelect(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_ENABLE(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolControl{}
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolControl(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_JOB(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req JobSet
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolJobSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)

}

func (s *Service) OnWS_TOOL_PSET(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetSet
	_ = json.Unmarshal(byteData, &req)

	err := s.ToolPSetSet(&req)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_PSET_BATCH(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetBatchSet
	_ = json.Unmarshal(byteData, &req)

	err := s.ToolPSetBatchSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_PSET_LIST(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	psetList, err := s.GetToolPSetList(&req)
	//psetList= append(append(psetList, 1), 2)
	//err= nil
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, psetList), s.diag)
}

func (s *Service) OnWS_TOOL_PSET_DETAIL(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolPSet
	_ = json.Unmarshal(byteData, &req)

	psetDetail, err := s.GetToolPSetDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, psetDetail), s.diag)
}

func (s *Service) OnWS_TOOL_JOB_LIST(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	jobList, err := s.GetToolJobList(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, jobList), s.diag)
}

func (s *Service) OnWS_TOOL_JOB_DETAIL(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolJob
	_ = json.Unmarshal(byteData, &req)

	jobDetail, err := s.GetToolJobDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SN, msg.Type, jobDetail), s.diag)
}

func (s *Service) OnWS_TOOL_RESULT_SET(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var result TighteningResult
	_ = json.Unmarshal(byteData, &result)

	if err := result.ValidateSet(); err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	tool, err := s.getTool(result.ControllerSN, result.ToolSN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	dbTool, err := s.storageService.GetTool(result.ToolSN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, err.Error()), s.diag)
		return
	}

	dbTool.Count = result.Count
	if err := s.storageService.UpdateTool(&dbTool); err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -4, err.Error()), s.diag)
		return
	}

	result.TighteningID = uuid.NewV4().String()
	result.UpdateTime = time.Now()
	s.doDispatch(tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherResult), result)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}
