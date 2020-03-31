package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/reader"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"strings"
	"time"
)

// 请求获取工单列表
func (s *Service) OnWSOrderList(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	workOrders, err := s.storageService.Workorders(byteData)
	if err != nil {
		s.diag.Error("Get WorkOrder Error", err)
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, workOrders))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body), s.diag)
}

// 请求获取工单详情
func (s *Service) OnWSOrderDetail(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	w, err := s.storageService.WorkorderOut("", orderReq.ID)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}
	if w == nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, fmt.Sprintf("找不到對應Id=%d的工單", orderReq.ID)), s.diag)
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, w))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body), s.diag)
}

// 更新工单状态
func (s *Service) OnWSOrderUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReq{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	currentOrder, err := s.storageService.GetWorkorderByCode(orderReq.WorkorderCode)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	_, err = s.storageService.UpdateWorkorder(&storage.Workorders{
		Id:     currentOrder.Id,
		Status: orderReq.Status,
	})

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

// 更新工步状态
func (s *Service) OnWSOrderStepUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	stepReq := WSStepReq{}
	err := json.Unmarshal(byteData, &stepReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	currentStep, err := s.storageService.GetStepByCode(stepReq.StepCode)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	_, err = s.storageService.UpdateStep(&storage.Steps{
		Id:     currentStep.Id,
		Status: stepReq.Status,
	})

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

// 开工请求
func (s *Service) OnWSOrderStart(c websocket.Connection, msg *wsnotify.WSMsg) {
	s.doDispatch(dispatcherbus.DispatcherOrderStart, msg.Data)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

// 完工请求
func (s *Service) OnWSOrderFinish(c websocket.Connection, msg *wsnotify.WSMsg) {
	s.doDispatch(dispatcherbus.DispatcherOrderFinish, msg.Data)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

// 更新工步数据
func (s *Service) OnWSOrderStepDataUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	stepReq := WSStepReqData{}
	err := json.Unmarshal(byteData, &stepReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	currentStep, err := s.storageService.GetStepByCode(stepReq.StepCode)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	var step storage.Steps
	step.Id = currentStep.Id
	step.Data = stepReq.Data
	_, err = s.storageService.UpdateStepData(&step)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

// 更新工单数据
func (s *Service) OnWSOrderDataUpdate(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqData{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	currentOrder, err := s.storageService.GetWorkorderByCode(orderReq.WorkorderCode)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	var order storage.Workorders
	order.Id = currentOrder.Id
	order.Data = orderReq.Data
	_, err = s.storageService.UpdateOrderData(&order)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""), s.diag)
}

// 根据CODE获取工单
func (s *Service) OnWSOrderDetailByCode(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	orderReq := WSOrderReqCode{}
	err := json.Unmarshal(byteData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	w, err := s.storageService.WorkorderOut(orderReq.Code, 0)
	if w == nil && err == nil {
		_, _ = s.backendService.GetWorkorder("", orderReq.Workcenter, orderReq.Code)
	}

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	if w == nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, "本地没有工单,正在向后台查询..."), s.diag)
		return
	}

	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, w))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body), s.diag)
}

func (s *Service) OnWSLocalResults(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)
	req := WSLocalResults{}
	err := json.Unmarshal(byteData, &req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	results, err := s.storageService.FindLocalResults(req.HmiSN, req.Limit)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	filters := strings.Join(req.Filters, ",")
	rt := []LocalResults{}
	sr := tightening_device.ResultValue{}
	for _, v := range results {
		stime := v.Results.UpdateTime.Format("2006-01-02 15:04:05")
		dt, _ := time.Parse("2006-01-02 15:04:05", stime)

		_ = json.Unmarshal([]byte(v.ResultValue), &sr)

		lr := LocalResults{
			HmiSN:        filterValue(filters, "hmi_sn", v.HMISN),
			Vin:          filterValue(filters, "vin", v.Vin),
			ControllerSN: filterValue(filters, "controller_sn", v.ControllerSN),
			ToolSN:       filterValue(filters, "tool_sn", v.ToolSN),
			Result:       filterValue(filters, "result", v.Result),
			Torque:       filterValue(filters, "torque", sr.Mi),
			Angle:        filterValue(filters, "angle", sr.Wi),
			Spent:        filterValue(filters, "spent", sr.Ti),
			TimeStamp:    filterValue(filters, "timestamp", dt.Local()),
			Batch:        filterValue(filters, "batch", v.Batch),
			VehicleType:  filterValue(filters, "vehicle_type", v.MO_Model),
			JobID:        filterValue(filters, "job_id", v.JobID),
			PSetID:       filterValue(filters, "pset_id", v.PSet),
		}

		rt = append(rt, lr)
	}

	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, rt))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, string(body), s.diag)
}

func filterValue(filters string, key string, value interface{}) interface{} {
	if filters == "" || strings.Contains(filters, key) {
		return value
	}

	return nil
}

func (s *Service) OnWSAuthUID(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)
	req := reader.ReaderUID{}
	err := json.Unmarshal(byteData, &req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()), s.diag)
		return
	}

	user, err := s.backendService.GetUserByUID(req.UID)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()), s.diag)
		return
	}

	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, user))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, string(body), s.diag)
}
