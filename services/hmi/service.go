package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/reader"
	"github.com/masami10/rush/services/scanner"
	"github.com/masami10/rush/utils"

	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
)

type Service struct {
	diag           Diagnostic
	storageService IStorageService
	httpd          HTTPService
	backendService IBackendService
	notifyService  INotifyService

	dispatcherBus Dispatcher

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(d Diagnostic, dp Dispatcher, ns INotifyService, httpd HTTPService, backend IBackendService, storage IStorageService) *Service {

	s := &Service{
		diag:           d,
		dispatcherBus:  dp,
		notifyService:  ns,
		httpd:          httpd,
		backendService: backend,
		storageService: storage,
	}

	s.setupWSRequestHandlers()
	s.setupTestInterface()

	return s
}

func (s *Service) Open() error {

	s.initDispatcherRegisters()

	return nil
}

func (s *Service) Close() error {
	s.diag.Close()
	s.diag.Closed()

	return nil
}

// 请求测试接口
func (s *Service) wsTest(ctx iris.Context) {
	ws := wsnotify.WSMsg{}
	err := ctx.ReadJSON(&ws)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	payload, _ := json.Marshal(ws)
	s.doDispatch(dispatcherbus.DispatcherWsNotify, &wsnotify.DispatcherNotifyPackage{
		C:    nil,
		Data: payload,
	})
}

func (s *Service) initDispatcherRegisters() {

	// 接收设备状态变化
	s.dispatcherBus.Register(dispatcherbus.DispatcherDeviceStatus, utils.CreateDispatchHandlerStruct(s.onDeviceStatus))

	// 接收读卡器数据
	s.dispatcherBus.Register(dispatcherbus.DispatcherReaderData, utils.CreateDispatchHandlerStruct(s.onReaderData))

	// 接收条码数据
	s.dispatcherBus.Register(dispatcherbus.DispatcherScannerData, utils.CreateDispatchHandlerStruct(s.onScannerData))

	// 接收IO输入输出状态变化
	s.dispatcherBus.Register(dispatcherbus.DispatcherIO, utils.CreateDispatchHandlerStruct(s.onIOContactData))

	// 注册websocket请求
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))

	// 接收拧紧结果
	s.dispatcherBus.Register(dispatcherbus.DispatcherResult, utils.CreateDispatchHandlerStruct(s.onTighteningResult))

	// 接收外部服务状态推送
	s.dispatcherBus.Register(dispatcherbus.DispatcherServiceStatus, utils.CreateDispatchHandlerStruct(s.onServiceStatus))

	// 接收新工单推送
	s.dispatcherBus.Register(dispatcherbus.DispatcherOrderNew, utils.CreateDispatchHandlerStruct(s.onNewOrder))

	// 接收保养信息推送
	s.dispatcherBus.Register(dispatcherbus.DispatcherMaintenanceInfo, utils.CreateDispatchHandlerStruct(s.onNewMaintenance))

	// 接收控制器JOB推送
	s.dispatcherBus.Register(dispatcherbus.DispatcherJob, utils.CreateDispatchHandlerStruct(s.onTighteningJob))
}

func (s *Service) setupTestInterface() {
	var r httpd.Route

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/ws-test",
		HandlerFunc: s.wsTest,
	}
	s.httpd.AddNewHttpHandler(r)
}

func (s *Service) setupWSRequestHandlers() {
	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

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
		wsnotify.WS_LOCAL_RESULTS:          s.OnWSLocalResults,
	})
}

// 收到工具拧紧结果
func (s *Service) onTighteningResult(data interface{}) {
	if data == nil {
		return
	}

	tighteningResult := data.(tightening_device.TighteningResult)

	// 拧紧结果推送HMI
	hmiTighteningResult := []tightening_device.BaseResult{tighteningResult.BaseResult}
	s.wsSendTighteningResult(hmiTighteningResult)

	body, _ := json.Marshal(hmiTighteningResult)
	s.diag.Info(fmt.Sprintf("拧紧结果推送HMI:%s", string(body)))
}

// 收到推送的JOB
func (s *Service) onTighteningJob(data interface{}) {
	if data == nil {
		return
	}

	job := data.(tightening_device.JobInfo)

	// JOB推送HMI
	s.wsSendTighteningJob(&job)

	body, _ := json.Marshal(job)
	s.diag.Info(fmt.Sprintf("JOB推送HMI:%s", string(body)))
}

// 设备连接状态变化(根据设备类型,序列号来区分具体的设备)
func (s *Service) onDeviceStatus(data interface{}) {
	if data == nil {
		return
	}

	deviceStatus := data.([]device.Status)

	// 设备状态推送HMI
	s.wsSendDeviceStatus(deviceStatus)

	body, _ := json.Marshal(deviceStatus)
	s.diag.Info(fmt.Sprintf("设备状态推送HMI:%s", string(body)))
}

// 收到IO状态变化(根据来源区分IO模块/控制器IO等)
func (s *Service) onIOContactData(data interface{}) {
	if data == nil {
		return
	}

	ioContact := data.(io.IoContact)

	// IO状态变化推送HMI
	s.wsSendIOContact(&ioContact)

	body, _ := json.Marshal(ioContact)
	s.diag.Info(fmt.Sprintf("IO状态变化推送HMI:%s", string(body)))
}

// 收到外部系统连接状态变化
func (s *Service) onServiceStatus(data interface{}) {
	if data == nil {
		return
	}

	status := data.(aiis.ServiceStatus)

	// 外部服务状态推送HMI
	s.wsSendServiceStatus(&status)

	body, _ := json.Marshal(status)
	s.diag.Info(fmt.Sprintf("外部服务状态推送HMI: %s", string(body)))
}

// 收到读卡器数据
func (s *Service) onReaderData(data interface{}) {
	if data == nil {
		return
	}

	readerData := data.(string)

	// 读卡器数据推送HMI
	s.wsSendReaderData(readerData)

	body, _ := json.Marshal(readerData)
	s.diag.Info(fmt.Sprintf("读卡器数据推送HMI:%s", string(body)))
}

// 收到条码数据
func (s *Service) onScannerData(data interface{}) {
	if data == nil {
		return
	}

	scannerData := data.(scanner.ScannerRead)

	// 条码数据推送HMI
	s.wsSendBarcode(&scannerData)

	body, _ := json.Marshal(scannerData)
	s.diag.Info(fmt.Sprintf("条码数据推送HMI:%s", string(body)))
}

// 收到新工单
func (s *Service) onNewOrder(data interface{}) {
	if data == nil {
		return
	}

	s.wsSendNewOrder(data)

	body, _ := json.Marshal(data)
	s.diag.Info(fmt.Sprintf("新工单推送HMI:%s", string(body)))
}

// 收到保养通知
func (s *Service) onNewMaintenance(data interface{}) {
	if data == nil {
		return
	}

	s.wsSendMaintenance(data)

	body, _ := json.Marshal(data)
	s.diag.Info(fmt.Sprintf("保养信息推送HMI:%s", string(body)))
}

// websocket发送新工单
func (s *Service) wsSendNewOrder(order interface{}) {
	var orders []interface{}
	orders = append(orders, order)
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(0, wsnotify.WS_ORDER_NEW, orders))
	s.notifyService.NotifyAll(wsnotify.WS_EVENT_ORDER, string(body))
}

// websocket发送保养通知
func (s *Service) wsSendMaintenance(info interface{}) {
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(0, wsnotify.WS_MAINTENANCE, info))
	s.notifyService.NotifyAll(wsnotify.WS_EVENT_MAINTENANCE, string(body))
}

// websocket发送读卡器信息
func (s *Service) wsSendReaderData(uid string) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_READER_UID,
		Data: &reader.ReaderUID{
			UID: uid,
		},
	})

	s.notifyService.NotifyAll(wsnotify.WS_EVENT_READER, string(data))
}

// websocket发送条码信息
func (s *Service) wsSendBarcode(scannerData *scanner.ScannerRead) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_SCANNER_READ,
		Data: scannerData,
	})

	s.notifyService.NotifyAll(wsnotify.WS_EVENT_SCANNER, string(data))
}

// websocket发送IO输入输出状态变化
func (s *Service) wsSendIOContact(ioContact *io.IoContact) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_IO_CONTACT,
		Data: ioContact,
	})

	s.notifyService.NotifyAll(wsnotify.WS_EVENT_IO, string(data))
}

// websocket发送设备状态
func (s *Service) wsSendDeviceStatus(deviceStatus []device.Status) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_DEVICE_STATUS,
		Data: deviceStatus,
	})

	s.notifyService.NotifyAll(wsnotify.WS_EVENT_DEVICE, string(data))
}

// websocket发送拧紧结果
func (s *Service) wsSendTighteningResult(results []tightening_device.BaseResult) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_TOOL_RESULT,
		Data: results,
	})

	s.notifyService.NotifyAll(wsnotify.WS_EVENT_TIGHTENING, string(data))
}

// websocket发送JOB
func (s *Service) wsSendTighteningJob(job *tightening_device.JobInfo) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_TOOL_JOB,
		Data: job,
	})

	s.notifyService.NotifyAll(wsnotify.WS_EVENT_TIGHTENING, string(data))
}

// websocket发送外部系统状态
func (s *Service) wsSendServiceStatus(status *aiis.ServiceStatus) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_READER_UID,
		Data: status,
	})

	s.notifyService.NotifyAll(wsnotify.WS_EVENT_SERVICE, string(data))
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error("Dispatch Failed", err)
	}
}
