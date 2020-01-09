package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
)

// 请求获取工单列表
func (s *Service) OnWSOrderList(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	workOrders, err := s.storageService.Workorders(byteData)
	if err != nil {
		s.diag.Error("Get WorkOrder Error", err)
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, workOrders))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))
}

// 请求获取工单详情
func (s *Service) OnWSOrderDetail(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	w, err := s.storageService.WorkorderOut("", orderReq.ID)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}
	if w == nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, fmt.Sprintf("找不到對應Id=%d的工單", orderReq.ID)))
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, w))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))
}

// 更新工单状态
func (s *Service) OnWSOrderUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.storageService.UpdateWorkorder(&storage.Workorders{
		Id:     orderReq.ID,
		Status: orderReq.Status,
	})

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 更新工步状态
func (s *Service) OnWSOrderStepUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.storageService.UpdateStep(&storage.Steps{
		Id:     orderReq.ID,
		Status: orderReq.Status,
	})

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 开工请求
func (s *Service) OnWSOrderStart(c websocket.Connection, msg *wsnotify.WSMsg) {
	s.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_ORDER_START, msg.Data)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 完工请求
func (s *Service) OnWSOrderFinish(c websocket.Connection, msg *wsnotify.WSMsg) {
	s.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_ORDER_FINISH, msg.Data)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

// 更新工步数据
func (s *Service) OnWSOrderStepDataUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqData{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.storageService.UpdateStepData(&storage.Steps{
		Id:   orderReq.ID,
		Data: orderReq.Data,
	})

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
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

	_, err = s.storageService.UpdateOrderData(&storage.Workorders{
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
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	w, err := s.storageService.WorkorderOut(orderReq.Code, 0)
	if w == nil && err == nil {
		s.backendService.GetWorkorder("", "", orderReq.Workcenter, orderReq.Code)
	}
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	if w == nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, "本地没有工单,正在向后台查询..."))
		return
	}

	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, w))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))
}
