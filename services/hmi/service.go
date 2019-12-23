package hmi

import (
	"encoding/json"
	"fmt"

	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/device"
	dispatcherBus "github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/scanner"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
)

const (
	CH_LENGTH = 1024
)

type Service struct {
	diag        Diagnostic
	DB          *storage.Service
	Httpd       *httpd.Service
	ODOO        *odoo.Service
	Aiis        *aiis.Service
	ChStart     chan int
	ChFinish    chan int
	ChWorkorder chan int
	SN          string
	WS          *wsnotify.Service

	DispatcherBus Dispatcher

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(d Diagnostic, dp Dispatcher) *Service {

	s := &Service{
		diag:          d,
		DispatcherBus: dp,
		ChStart:       make(chan int, CH_LENGTH),
		ChFinish:      make(chan int, CH_LENGTH),
		ChWorkorder:   make(chan int, CH_LENGTH),
	}

	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.initWSRequestHandlers()

	return s
}

func (s *Service) SendScannerInfo(identification string) error {
	if s.WS == nil {
		return errors.New("Please Inject Notify Service First")
	}
	s.WS.NotifyAll(wsnotify.WS_EVENT_SCANNER, identification)
	return nil
}

func (s *Service) Open() error {
	//fixme: 嵌套那么深
	//s.ControllerService.WS.OnNewClient = s.OnNewHmiConnect

	// 注册websocket请求
	s.DispatcherBus.Register(dispatcherBus.DISPATCH_WS_NOTIFY, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))
	return nil
}

func (s *Service) Close() error {
	s.diag.Close()
	s.diag.Closed()

	return nil
}

func (s *Service) initWSRequestHandlers() {
	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_ORDER_LIST:             s.OnWSOrderList,
		wsnotify.WS_ORDER_DETAIL:           s.OnWSOrderDetail,
		wsnotify.WS_ORDER_UPDATE:           s.OnWSOrderUpdate,
		wsnotify.WS_ORDER_DATA_UPDATE:      s.OnWSOrderDataUpdate,
		wsnotify.WS_ORDER_START:            s.OnWSOrderStart,
		wsnotify.WS_ORDER_FINISH:           s.OnWSOrderFinish,
		wsnotify.WS_ORDER_STEP_UPDATE:      s.OnWSOrderStepUpdate,
		wsnotify.WS_ORDER_STEP_DATA_UPDATE: s.OnWSOrderStepDataUpdate,
		wsnotify.WS_ORDER_DETAIL_BY_CODE:   s.OnWSOrderDetailByCode,
	})
}

func (s *Service) resetResult(id int64) error {
	// 重置结果
	err := s.DB.ResetResult(id)
	if err != nil {
		return err
	}

	// 删除对应曲线
	err = s.DB.DeleteCurvesByResult(id)
	if err != nil {
		return err
	}
	s.diag.Debug("reset Result Success")
	return nil
}

func (s *Service) OnNewHmiConnect(conn websocket.Connection) {

	//主动推送工位号
	msg := WSWorkcenter{
		WorkCenter: s.WS.Config().Workcenter,
	}
	_ = s.commonSendWebSocketMsg(conn, wsnotify.WS_EVENT_REG, wsnotify.GenerateWSMsg(0, wsnotify.WS_RUSH_DATA, msg))
}

func (s *Service) commonSendWebSocketMsg(c websocket.Connection, subject string, msg interface{}) error {
	err := wsnotify.WSClientSend(c, subject, msg)
	if err != nil {
		s.diag.Error("commonSendWebSocketMsg Error", err)
	}
	return err
}

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

// 收到工具拧紧结果
func (s *Service) OnTighteningResult(data interface{}) {
	if data == nil {
		return
	}

	tighteningResult := data.(*tightening_device.TighteningResult)
	result := []wsnotify.WSResult{
		{
			Result: tighteningResult.MeasureResult,
			MI:     tighteningResult.MeasureTorque,
			WI:     tighteningResult.MeasureAngle,
			TI:     tighteningResult.MeasureTime,
			Seq:    tighteningResult.Seq,
			Batch:  tighteningResult.Batch,
			ToolSN: tighteningResult.ToolSN,
		},
	}

	msg := wsnotify.WSMsg{
		Type: wsnotify.WS_TIGHTENING_RESULT,
		Data: result,
	}

	payload, _ := json.Marshal(msg)

	s.WS.NotifyAll(wsnotify.WS_EVENT_RESULT, string(payload))
	s.diag.Debug(fmt.Sprintf("拧紧结果推送HMI: %s", string(payload)))
}

//// 控制器状态变化
//func (s *Service) OnTighteningControllerStatus(data interface{}) {
//	if data == nil {
//		return
//	}
//
//	controllerStatus := data.(*[]device.DeviceStatus)
//
//	msg := wsnotify.WSMsg{
//		Type: wsnotify.NotifywsDeviceStatus,
//		Data: controllerStatus,
//	}
//	payload, _ := json.Marshal(msg)
//	s.WS.NotifyAll(wsnotify.WS_EVENT_DEVICE, string(payload))
//	s.diag.Debug(fmt.Sprintf("控制器状态推送HMI: %s", string(payload)))
//}

// 设备连接状态变化(根据设备类型,序列号来区分具体的设备)
func (s *Service) OnDeviceStatus(data interface{}) {
	if data == nil {
		return
	}

	toolStatus := data.(*[]device.DeviceStatus)

	msg := wsnotify.WSMsg{
		Type: wsnotify.WS_DEVICE_STATUS,
		Data: toolStatus,
	}
	payload, _ := json.Marshal(msg)
	s.WS.NotifyAll(wsnotify.WS_EVENT_DEVICE, string(payload))
	s.diag.Debug(fmt.Sprintf("工具状态推送HMI: %s", string(payload)))
}

// 收到IO状态变化(根据来源区分IO模块/控制器IO等)
func (s *Service) OnTighteningControllerIO(data interface{}) {
	if data == nil {
		return
	}

	ioContact := data.(*io.IoContact)

	msg := wsnotify.WSMsg{
		Type: wsnotify.WS_IO_CONTACT,
		Data: ioContact,
	}
	payload, _ := json.Marshal(msg)
	s.WS.NotifyAll(wsnotify.WS_EVENT_IO, string(payload))
	s.diag.Debug(fmt.Sprintf("控制器IO推送HMI: %s", string(payload)))
}

// 收到条码(根据来源区分扫码枪/控制器条码等)
func (s *Service) OnTighteningControllereID(data interface{}) {
	if data == nil {
		return
	}

	controllerBarcode := data.(*scanner.ScannerRead)

	msg := wsnotify.WSMsg{
		Type: scanner.WS_SCANNER_READ,
		Data: controllerBarcode,
	}

	payload, _ := json.Marshal(msg)
	s.WS.NotifyAll(wsnotify.WS_EVENT_SCANNER, string(payload))
	s.diag.Debug(fmt.Sprintf("控制器条码推送HMI: %s", string(payload)))
}

// 收到Aiis连接状态变化
func (s *Service) OnAiisStatus(data interface{}) {
	if data == nil {
		return
	}

	status := data.(string)

	// Aiis状态推送给HMI
	payload, _ := json.Marshal(&wsnotify.WSMsg{
		Type: WS_AIIS_STATUS,
		Data: &aiis.SystemStatus{
			Name:   "AIIS",
			Status: status,
		},
	})

	s.WS.NotifyAll(wsnotify.WS_EVENT_AIIS, string(payload))
	s.diag.Debug(fmt.Sprintf("Aiis连接状态推送HMI: %s", string(payload)))
}

// 收到Odoo连接状态变化
func (s *Service) OnOdooStatus(data interface{}) {
	if data == nil {
		return
	}
	status := data.(string)

	// Odoo状态推送给HMI
	payload, _ := json.Marshal(&wsnotify.WSMsg{
		Type: WS_ODOO_STATUS,
		Data: &aiis.SystemStatus{
			Name:   "ODOO",
			Status: status,
		},
	})

	s.WS.NotifyAll(wsnotify.WS_EVENT_ODOO, string(payload))
	s.diag.Debug(fmt.Sprintf("ODOO连接状态推送HMI: %s", string(payload)))
}

// 收到第三方系统状态变化
func (s *Service) OnExSysStatus(data interface{}) {
	if data == nil {
		return
	}
	status := data.(*aiis.SystemStatus)

	// 第三方系统状态推送给HMI
	payload, _ := json.Marshal(&wsnotify.WSMsg{
		Type: WS_EXSYS_STATUS,
		Data: status,
	})

	s.WS.NotifyAll(wsnotify.WS_EVENT_EXSYS, string(payload))
	s.diag.Debug(fmt.Sprintf("第三方系统连接状态推送HMI: %s", string(payload)))
}
