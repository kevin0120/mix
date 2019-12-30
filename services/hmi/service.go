package hmi

import (
	"encoding/json"
	"fmt"
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

	// 接收设备状态变化
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_DEVICE_STATUS, utils.CreateDispatchHandlerStruct(s.onDeviceStatus))

	// 接收读卡器数据
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_READER_DATA, utils.CreateDispatchHandlerStruct(s.onReaderData))

	// 接收条码数据
	s.DispatcherBus.Register(dispatcherBus.DISPATCHER_SCANNER_DATA, utils.CreateDispatchHandlerStruct(s.onScannerData))

	// 接收IO输入输出状态变化
	s.DispatcherBus.Register(dispatcherBus.DISPATCH_IO, utils.CreateDispatchHandlerStruct(s.onIOContactData))

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

// 收到工具拧紧结果
func (s *Service) onTighteningResult(data interface{}) {
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

// TODO: 外部系统状态推送
// 收到Aiis连接状态变化
//func (s *Service) OnAiisStatus(data interface{}) {
//	if data == nil {
//		return
//	}
//
//	status := data.(string)
//
//	// Aiis状态推送给HMI
//	payload, _ := json.Marshal(&wsnotify.WSMsg{
//		Type: WS_AIIS_STATUS,
//		Data: &aiis.SystemStatus{
//			Name:   "AIIS",
//			Status: status,
//		},
//	})
//
//	s.WS.NotifyAll(wsnotify.WS_EVENT_AIIS, string(payload))
//	s.diag.Debug(fmt.Sprintf("Aiis连接状态推送HMI: %s", string(payload)))
//}
//
//// 收到Odoo连接状态变化
//func (s *Service) OnOdooStatus(data interface{}) {
//	if data == nil {
//		return
//	}
//	status := data.(string)
//
//	// Odoo状态推送给HMI
//	payload, _ := json.Marshal(&wsnotify.WSMsg{
//		Type: WS_ODOO_STATUS,
//		Data: &aiis.SystemStatus{
//			Name:   "ODOO",
//			Status: status,
//		},
//	})
//
//	s.WS.NotifyAll(wsnotify.WS_EVENT_ODOO, string(payload))
//	s.diag.Debug(fmt.Sprintf("ODOO连接状态推送HMI: %s", string(payload)))
//}
//
//// 收到第三方系统状态变化
//func (s *Service) OnExSysStatus(data interface{}) {
//	if data == nil {
//		return
//	}
//	status := data.(*aiis.SystemStatus)
//
//	// 第三方系统状态推送给HMI
//	payload, _ := json.Marshal(&wsnotify.WSMsg{
//		Type: WS_EXSYS_STATUS,
//		Data: status,
//	})
//
//	s.WS.NotifyAll(wsnotify.WS_EVENT_EXSYS, string(payload))
//	s.diag.Debug(fmt.Sprintf("第三方系统连接状态推送HMI: %s", string(payload)))
//}

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

	s.WS.NotifyAll(wsnotify.WS_EVENT_READER, string(data))
}

// websocket发送条码信息
func (s *Service) wsSendBarcode(scannerData *scanner.ScannerRead) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_SCANNER_READ,
		Data: scannerData,
	})

	s.WS.NotifyAll(wsnotify.WS_EVENT_SCANNER, string(data))
}

// websocket发送IO输入输出状态变化
func (s *Service) wsSendIOContact(ioContact *io.IoContact) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_IO_CONTACT,
		Data: ioContact,
	})

	s.WS.NotifyAll(wsnotify.WS_EVENT_IO, string(data))
}

// websocket发送设备状态
func (s *Service) wsSendDeviceStatus(deviceStatus []device.DeviceStatus) {
	data, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_DEVICE_STATUS,
		Data: deviceStatus,
	})

	s.WS.NotifyAll(wsnotify.WS_EVENT_DEVICE, string(data))
}
