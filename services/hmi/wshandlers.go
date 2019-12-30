package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
)

// 请求获取工单列表
func (s *Service) OnWSOrderList(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	workOrders, err := s.DB.Workorders(byteData)
	if err != nil {
		s.diag.Error("Get WorkOrder Error", err)
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, workOrders))
	_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_ORDER, string(body))
}

// 请求获取工单详情
func (s *Service) OnWSOrderDetail(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	w, err := s.DB.WorkorderOut("", orderReq.ID)
	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}
	if w == nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, fmt.Sprintf("找不到對應Id=%d的工單", orderReq.ID)))
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, w))
	_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_ORDER, string(body))
}

// 更新工单状态
func (s *Service) OnWSOrderUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.DB.UpdateWorkorder(&storage.Workorders{
		Id:     orderReq.ID,
		Status: orderReq.Status,
	})

	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 更新工步状态
func (s *Service) OnWSOrderStepUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.DB.UpdateStep(&storage.Steps{
		Id:     orderReq.ID,
		Status: orderReq.Status,
	})

	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 开工请求
func (s *Service) OnWSOrderStart(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqCode{}
	err := json.Unmarshal(byteData, &orderReq)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}
	s.ChStart <- 1
	go s.Aiis.PutMesOpenRequest(msg.SN, msg.Type, orderReq.Code, byteData, s.ChStart)
}

// 完工请求
func (s *Service) OnWSOrderFinish(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqCode{}
	err := json.Unmarshal(byteData, &orderReq)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}
	s.ChFinish <- 1
	go s.Aiis.PutMesFinishRequest(msg.SN, msg.Type, orderReq.Code, byteData, s.ChFinish)

}

// 更新工步数据
func (s *Service) OnWSOrderStepDataUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqData{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.DB.UpdateStepData(&storage.Steps{
		Id:   orderReq.ID,
		Data: orderReq.Data,
	})

	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 更新工单数据
func (s *Service) OnWSOrderDataUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqData{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.DB.UpdateOrderData(&storage.Workorders{
		Id:   orderReq.ID,
		Data: orderReq.Data,
	})

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 根据CODE获取工单
func (s *Service) OnWSOrderDetailByCode(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqCode{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	w, err := s.DB.WorkorderOut(orderReq.Code, 0)
	//todo 判定本地无工单
	if w == nil && err == nil {
		s.ChWorkorder <- 1
		go s.ODOO.GetWorkorder("", "", orderReq.Workcenter, orderReq.Code, s.ChWorkorder)
	}
	if err != nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	if w == nil {
		_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, "本地没有工单,正在向后台查询..."))
		return
	}

	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, w))
	//fmt.Println(string(body))
	_ = s.commonSendWebSocketMsg(c, wsnotify.WS_EVENT_ORDER, string(body))
}
