package hmi

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/reader"
	"github.com/masami10/rush/services/scanner"

	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/device"
	dispatcherBus "github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/odoo"
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
	diag          Diagnostic
	DB            *storage.Service
	Httpd         *httpd.Service
	ODOO          *odoo.Service
	Aiis          *aiis.Service
	ChWorkorder   chan int
	SN            string
	NotifyService INotifyService

	DispatcherBus Dispatcher

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(d Diagnostic, dp Dispatcher, ns INotifyService) *Service {

	s := &Service{
		diag:          d,
		DispatcherBus: dp,
		ChWorkorder:   make(chan int, CH_LENGTH),
		NotifyService: ns,
	}

	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.initWSRequestHandlers()

	return s
}

func (s *Service) SendScannerInfo(identification string) error {
	if s.NotifyService == nil {
		return errors.New("Please Inject Notify Service First")
	}
	s.NotifyService.NotifyAll(wsnotify.WS_EVENT_SCANNER, identification)
	return nil
}

func (s *Service) Open() error {
	//fixme: 嵌套那么深
	//s.ControllerService.notifyService.OnNewClient = s.OnNewHmiConnect

	// 接收设备状态变化
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_DEVICE_STATUS, utils.CreateDispatchHandlerStruct(s.onDeviceStatus))

	// 接收读卡器数据
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_READER_DATA, utils.CreateDispatchHandlerStruct(s.onReaderData))

	// 接收条码数据
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_SCANNER_DATA, utils.CreateDispatchHandlerStruct(s.onScannerData))

	// 接收IO输入输出状态变化
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_IO, utils.CreateDispatchHandlerStruct(s.onIOContactData))

	// 注册websocket请求
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_WS_NOTIFY, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))

	// 接收拧紧结果
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_RESULT, utils.CreateDispatchHandlerStruct(s.onTighteningResult))

	// 接收外部服务状态推送
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_SERVICE_STATUS, utils.CreateDispatchHandlerStruct(s.onServiceStatus))

	var r httpd.Route

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/ws-test",
		HandlerFunc: s.wsTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	return nil
}

func (s *Service) Close() error {
	s.diag.Close()
	s.diag.Closed()

	return nil
}

func (s *Service) wsTest(ctx iris.Context) {
	ws := wsnotify.WSMsg{}
	err := ctx.ReadJSON(&ws)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	payload, _ := json.Marshal(ws)
	s.DispatcherBus.Dispatch(dispatcherBus.DISPATCHER_WS_NOTIFY, &wsnotify.DispatcherNotifyPackage{
		C:    nil,
		Data: payload,
	})
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
		WorkCenter: s.NotifyService.GetWorkCenter(),
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

// 收到工具拧紧结果
func (s *Service) onTighteningResult(data interface{}) {
	if data == nil {
		return
	}

	tighteningResult := data.(tightening_device.TighteningResult)

	// 拧紧结果推送HMI
	s.wsSendTighteningResult([]tightening_device.BaseResult{tighteningResult.BaseResult})

	body, _ := json.Marshal(tighteningResult)
	s.diag.Info(fmt.Sprintf("拧紧结果推送HMI:%s", string(body)))
}

// 设备连接状态变化(根据设备类型,序列号来区分具体的设备)
func (s *Service) onDeviceStatus(data interface{}) {
	if data == nil {
		return
	}

	deviceStatus := data.([]device.DeviceStatus)

	// 设备状态推送HMI
	s.wsSendDeviceStatus(deviceStatus)
}

// 收到IO状态变化(根据来源区分IO模块/控制器IO等)
func (s *Service) onIOContactData(data interface{}) {
	if data == nil {
		return
	}

	ioContact := data.(io.IoContact)

	// IO状态变化推送HMI
	s.wsSendIOContact(&ioContact)
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
}

// 收到条码数据
func (s *Service) onScannerData(data interface{}) {
	if data == nil {
		return
	}

	scannerData := data.(scanner.ScannerRead)

	// 条码数据推送HMI
	s.wsSendBarcode(&scannerData)
}

// websocket发送读卡器信息
func (s *Service) wsSendReaderData(uid string) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_READER_UID,
		Data: &reader.ReaderUID{
			UID: uid,
		},
	})

	s.NotifyService.NotifyAll(wsnotify.WS_EVENT_READER, string(data))
}

// websocket发送条码信息
func (s *Service) wsSendBarcode(scannerData *scanner.ScannerRead) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_SCANNER_READ,
		Data: scannerData,
	})

	s.NotifyService.NotifyAll(wsnotify.WS_EVENT_SCANNER, string(data))
}

// websocket发送IO输入输出状态变化
func (s *Service) wsSendIOContact(ioContact *io.IoContact) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_IO_CONTACT,
		Data: ioContact,
	})

	s.NotifyService.NotifyAll(wsnotify.WS_EVENT_IO, string(data))
}

// websocket发送设备状态
func (s *Service) wsSendDeviceStatus(deviceStatus []device.DeviceStatus) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_DEVICE_STATUS,
		Data: deviceStatus,
	})

	s.NotifyService.NotifyAll(wsnotify.WS_EVENT_DEVICE, string(data))
}

// websocket发送拧紧结果
func (s *Service) wsSendTighteningResult(results []tightening_device.BaseResult) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_TIGHTENING_RESULT,
		Data: results,
	})

	s.NotifyService.NotifyAll(wsnotify.WS_EVENT_RESULT, string(data))
}

// websocket发送外部系统状态
func (s *Service) wsSendServiceStatus(status *aiis.ServiceStatus) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_READER_UID,
		Data: status,
	})

	s.NotifyService.NotifyAll(wsnotify.WS_EVENT_SERVICE, string(data))
}
