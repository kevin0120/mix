package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/scanner"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Disconnect(id string)
	Close()
	Closed()
}

type Service struct {
	diag              Diagnostic
	methods           Methods
	DB                *storage.Service
	Httpd             *httpd.Service
	ODOO              *odoo.Service
	AudiVw            *audi_vw.Service
	OpenProtocol      *openprotocol.Service
	ControllerService *controller.Service

	SN                 string
	WS                 *wsnotify.Service
	TighteningService  *tightening_device.Service
	requestDispatchers map[string]*utils.Dispatcher
}

func NewService(d Diagnostic) *Service {

	s := &Service{
		diag:    d,
		methods: Methods{},
	}

	s.methods.service = s

	return s
}

func (s *Service) Open() error {

	s.ControllerService.WS.OnNewClient = s.OnNewHmiClient
	s.TighteningService.GetDispatcher(tightening_device.DISPATCH_RESULT).Register(s.OnTighteningResult)
	s.TighteningService.GetDispatcher(tightening_device.DISPATCH_IO).Register(s.OnTighteningControllerIO)
	s.TighteningService.GetDispatcher(tightening_device.DISPATCH_CONTROLLER_STATUS).Register(s.OnTighteningControllerStatus)
	s.TighteningService.GetDispatcher(tightening_device.DISPATCH_TOOL_STATUS).Register(s.OnTighteningToolStatus)
	s.TighteningService.GetDispatcher(tightening_device.DISPATCH_CONTROLLER_ID).Register(s.OnTighteningControllereID)
	s.ODOO.Aiis.OdooStatusDispatcher.Register(s.OnOdooStatus)
	s.ODOO.Aiis.AiisStatusDispatcher.Register(s.OnAiisStatus)

	s.initRequestDispatchers()

	var r httpd.Route

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/tool-enable",
		HandlerFunc: s.methods.putToolControl,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/psets",
		HandlerFunc: s.methods.putPSets,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/psets-manual",
		HandlerFunc: s.methods.putManualPSets,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/psets",
		HandlerFunc: s.methods.getPSetList,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/pset-detail",
		HandlerFunc: s.methods.getPSetDetail,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/jobs",
		HandlerFunc: s.methods.getJobList,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/job-detail",
		HandlerFunc: s.methods.getJobDetail,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/jobs",
		HandlerFunc: s.methods.putJobs,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/jobs-manual",
		HandlerFunc: s.methods.putManualJobs,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/controller-mode",
		HandlerFunc: s.methods.enableJobMode,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/workorder",
		HandlerFunc: s.methods.getWorkorder,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/next-workorder",
		HandlerFunc: s.methods.getNextWorkorder,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/controller-status",
		HandlerFunc: s.methods.getStatus,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/healthz",
		HandlerFunc: s.methods.getHealthz,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/hmi-results",
		HandlerFunc: s.methods.getHmiResults,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/io-sets",
		HandlerFunc: s.methods.putIOSet,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/scanner",
		HandlerFunc: s.methods.putBarcodeTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/input-test",
		HandlerFunc: s.methods.putIOInputTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/result-test",
		HandlerFunc: s.methods.putResultTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/job-control",
		HandlerFunc: s.methods.putJobAbort,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/operation/{code:string}",
		HandlerFunc: s.methods.getRoutingOpertions,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/local-results",
		HandlerFunc: s.methods.getLocalResults,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/workorders",
		HandlerFunc: s.methods.listWorkorders,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/ak2",
		HandlerFunc: s.methods.postAK2,
	}
	s.Httpd.Handler[0].AddRoute(r)

	//r = httpd.Route{
	//	RouteType:   httpd.ROUTE_TYPE_HTTP,
	//	Method:      "PUT",
	//	Pattern:     "/test-protocol",
	//	HandlerFunc: s.methods.testProtocol,
	//}
	//s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/WS-test",
		HandlerFunc: s.methods.wsTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/storage",
		HandlerFunc: s.methods.getStorage,
	}
	s.Httpd.Handler[0].AddRoute(r)

	s.ODOO.WS.AddNotify(s)

	return nil

}

func (s *Service) Close() error {
	s.diag.Close()
	s.diag.Closed()

	for _, v := range s.requestDispatchers {
		v.Release()
	}

	return nil
}

func (s *Service) initRequestDispatchers() {
	s.requestDispatchers = map[string]*utils.Dispatcher{
		WS_ORDER_LIST:             utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_ORDER_DETAIL:           utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_ORDER_UPDATE:           utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_ORDER_STEP_UPDATE:      utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_ORDER_START_REQUEST:    utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_ORDER_FINISH_REQUEST:   utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_ORDER_STEP_DATA_UPDATE: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_ORDER_DETAIL_BY_CODE:   utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
	}

	s.requestDispatchers[WS_ORDER_LIST].Register(s.OnWS_ORDER_LIST)
	s.requestDispatchers[WS_ORDER_DETAIL].Register(s.OnWS_ORDER_DETAIL)
	s.requestDispatchers[WS_ORDER_UPDATE].Register(s.OnWS_ORDER_UPDATE)
	s.requestDispatchers[WS_ORDER_STEP_UPDATE].Register(s.OnWS_ORDER_STEP_UPDATE)
	s.requestDispatchers[WS_ORDER_START_REQUEST].Register(s.OnWS_ORDER_START_REQUEST)
	s.requestDispatchers[WS_ORDER_FINISH_REQUEST].Register(s.OnWS_ORDER_FINISH_REQUEST)
	s.requestDispatchers[WS_ORDER_STEP_DATA_UPDATE].Register(s.OnWS_ORDER_STEP_DATA_UPDATE)
	s.requestDispatchers[WS_ORDER_DETAIL_BY_CODE].Register(s.OnWS_ORDER_DETAIL_BY_CODE)

	for _, v := range s.requestDispatchers {
		v.Start()
	}
}

func (s *Service) dispatchRequest(req *wsnotify.WSRequest) {
	d, exist := s.requestDispatchers[req.WSMsg.Type]
	if exist {
		d.Dispatch(req)
	}
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

	return nil
}

func (s *Service) OnNewHmiClient(conn websocket.Connection) {

	//主动推送工位号
	msg := WSWorkcenter{
		WorkCenter: s.WS.Config().Workcenter,
	}
	_ = conn.Emit(wsnotify.WS_EVENT_REG, wsnotify.GenerateMessage(0, wsnotify.WS_RUSH_DATA, msg))
}

// 请求获取工单列表
func (s *Service) OnWS_ORDER_LIST(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	workorders, err := s.DB.Workorders(sData)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, workorders))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))
}

// 请求获取工单详情
func (s *Service) OnWS_ORDER_DETAIL(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	orderReq := WSOrderReq{}
	err := json.Unmarshal(sData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	w, err := s.DB.WorkorderOut("", orderReq.ID)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}
	if w == nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -3, fmt.Sprintf("找不到對應Id=%d的工單", orderReq.ID)))
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, w))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))
}

// 更新工单状态
func (s *Service) OnWS_ORDER_UPDATE(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	orderReq := WSOrderReq{}
	err := json.Unmarshal(sData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.DB.UpdateWorkorder(&storage.Workorders{
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
func (s *Service) OnWS_ORDER_STEP_UPDATE(data interface{}) {

	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	orderReq := WSOrderReq{}
	err := json.Unmarshal(sData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.DB.UpdateStep(&storage.Steps{
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
func (s *Service) OnWS_ORDER_START_REQUEST(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	orderReq := WSOrderReqCode{}
	err := json.Unmarshal(sData, &orderReq)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}
	resp, err := s.Aiis.PutMesOpenRequest(orderReq.Code, sData)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}
	s.diag.Debug(fmt.Sprintf("Mes处理开工请求的结果推送HMI: %s", resp.(string)))
	//_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, resp.(string)))
	body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, resp))
	s.WS.WSSend(wsnotify.WS_EVENT_ORDER, string(body))
}

// 完工请求
func (s *Service) OnWS_ORDER_FINISH_REQUEST(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	orderReq := WSOrderReqCode{}
	err := json.Unmarshal(sData, &orderReq)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}
	resp, err := s.Aiis.PutMesFinishRequest(orderReq.Code, sData)

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}
	s.diag.Debug(fmt.Sprintf("Mes处理完工请求的结果推送HMI: %s", resp.(string)))
	//_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
	body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, resp))
	s.WS.WSSend(wsnotify.WS_EVENT_ORDER, string(body))

}

// 更新工步数据
func (s *Service) OnWS_ORDER_STEP_DATA_UPDATE(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	orderReq := WSOrderReqData{}
	err := json.Unmarshal(sData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_, err = s.DB.UpdateStepData(&storage.Steps{
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
func (s *Service) OnWS_ORDER_DETAIL_BY_CODE(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	sData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	orderReq := WSOrderReqCode{}
	err := json.Unmarshal(sData, &orderReq)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	w, err := s.DB.WorkorderOut(orderReq.Code, 0)
	//todo 判定本地无工单
	if w == nil && err == nil {
		fmt.Println("如果RUSH收到HMI请求后找不到新工单,可通过调用ODOO api获取对应工单并推送HMI")
		orderData, err := s.ODOO.GetWorkorder("", "", orderReq.Workcenter, orderReq.Code)

		if err == nil {
			s.ODOO.HandleWorkorder(orderData)
		}

	}
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}
	body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, w))
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))
}

func (s *Service) OnWSMsg(c websocket.Connection, data []byte) {
	msg := wsnotify.WSMsg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(string(data), err)
		return
	}

	s.dispatchRequest(&wsnotify.WSRequest{
		C:     c,
		WSMsg: &msg,
	})
}

// 收到控制器结果
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
		Type: tightening_device.WS_TIGHTENING_RESULT,
		Data: result,
	}

	payload, _ := json.Marshal(msg)

	s.WS.WSSend(wsnotify.WS_EVENT_RESULT, string(payload))
	s.diag.Debug(fmt.Sprintf("拧紧结果推送HMI: %s", string(payload)))
}

// 控制器状态变化
func (s *Service) OnTighteningControllerStatus(data interface{}) {
	if data == nil {
		return
	}

	controllerStatus := data.(*[]device.DeviceStatus)

	msg := wsnotify.WSMsg{
		Type: device.WS_DEVICE_STATUS,
		Data: controllerStatus,
	}
	payload, _ := json.Marshal(msg)
	s.WS.WSSend(wsnotify.WS_EVENT_DEVICE, string(payload))
	s.diag.Debug(fmt.Sprintf("控制器状态推送HMI: %s", string(payload)))
}

// 工具状态变化
func (s *Service) OnTighteningToolStatus(data interface{}) {
	if data == nil {
		return
	}

	toolStatus := data.(*[]device.DeviceStatus)

	msg := wsnotify.WSMsg{
		Type: device.WS_DEVICE_STATUS,
		Data: toolStatus,
	}
	payload, _ := json.Marshal(msg)
	s.WS.WSSend(wsnotify.WS_EVENT_DEVICE, string(payload))
	s.diag.Debug(fmt.Sprintf("工具状态推送HMI: %s", string(payload)))
}

// 控制器IO变化
func (s *Service) OnTighteningControllerIO(data interface{}) {
	if data == nil {
		return
	}

	inputStatus := data.(*tightening_device.TighteningControllerInput)

	msg := wsnotify.WSMsg{
		Type: io.WS_IO_CONTACT,
		Data: inputStatus,
	}
	payload, _ := json.Marshal(msg)
	s.WS.WSSend(wsnotify.WS_EVENT_IO, string(payload))
	s.diag.Debug(fmt.Sprintf("控制器IO推送HMI: %s", string(payload)))
}

// 收到控制器条码
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
	s.WS.WSSend(wsnotify.WS_EVENT_SCANNER, string(payload))
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

	s.WS.WSSend(wsnotify.WS_EVENT_AIIS, string(payload))
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

	s.WS.WSSend(wsnotify.WS_EVENT_ODOO, string(payload))
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

	s.WS.WSSend(wsnotify.WS_EVENT_EXSYS, string(payload))
	s.diag.Debug(fmt.Sprintf("第三方系统连接状态推送HMI: %s", string(payload)))
}
