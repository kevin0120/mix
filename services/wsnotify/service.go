package wsnotify

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/utils"
	"sync/atomic"
)

const (
	WS_EVENT_CONTROLLER        = "controller"
	WS_EVENT_RESULT            = "result"
	WS_EVENT_REG               = "register"
	WS_EVENT_SELECTOR          = "selector"
	WS_EVENT_JOB               = "job"
	WS_EVENT_SCANNER           = "scanner"
	WS_EVENT_IO                = "io"
	WS_EVENT_ODOO              = "odoo"
	WS_EVENT_AIIS              = "aiis"
	WS_EVENT_EXSYS             = "exsys"
	WS_EVENT_MAINTENANCE       = "maintenance"
	WS_EVENT_TOOL              = "tool"
	WS_EVENT_READER            = "reader"
	WS_EVENT_TIGHTENING_DEVICE = "tightening_device"
	WS_EVENT_REPLY             = "reply"
	WS_EVENT_DEVICE            = "device"
	WS_EVENT_ORDER             = "order"
	WS_EVENT_MES               = "mes"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	OnMessage(msg string)
	Close()
	Closed()
}
type OnNewClient func(c websocket.Connection)

//type WSNotify interface {
//	OnWSMsg(c websocket.Connection, data []byte)
//}

type Service struct {
	configValue atomic.Value
	diag        Diagnostic

	ws *websocket.Server

	Httpd HTTPService

	clientManager *WSClientManager

	OnNewClient OnNewClient

	dispatcherBus Dispatcher

	dispatcherArray []*utils.DispatchHandlerStruct
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) NewWebSocketRecvHandler(handler func(interface{})) {
	fn := utils.CreateDispatchHandlerStruct(handler)
	s.dispatcherBus.Register(dispatcherbus.DISPATCH_WS_NOTIFY, fn)
}

func (s *Service) onConnect(c websocket.Connection) {

	c.OnMessage(func(data []byte) {
		msg := WSMsg{}
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return
		}

		if msg.Type == WS_REG {
			var reg WSRegist
			strData, _ := json.Marshal(msg.Data)
			err := json.Unmarshal([]byte(strData), &reg)
			if err != nil {
				_ = c.Disconnect()
				s.clientManager.RemoveClientBySN(reg.HMI_SN)

				_ = c.Emit(WS_EVENT_REPLY, GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			}

			_, exist := s.clientManager.GetClient(reg.HMI_SN)
			if exist {
				Msg := fmt.Sprintf("client with sn:%s already exists", reg.HMI_SN)
				_ = c.Disconnect()
				s.clientManager.RemoveClientBySN(reg.HMI_SN)

				_ = c.Emit(WS_EVENT_REPLY, GenerateReply(msg.SN, msg.Type, -2, Msg))

			} else {
				// 将客户端加入列表
				s.clientManager.AddClient(reg.HMI_SN, c)

				if s.OnNewClient != nil {
					s.OnNewClient(c)
				}

				// 注册成功
				_ = c.Emit(WS_EVENT_REPLY, GenerateReply(msg.SN, msg.Type, 0, ""))

			}
		} else {
			s.postNotify(&DispatcherNotifyPackage{
				C:    c,
				Data: data,
			})
		}
	})

	c.OnDisconnect(func() {
		s.clientManager.RemoveClient(c.ID())
		s.diag.Disconnect(c.ID())
	})

	c.OnError(func(err error) {
		s.diag.Error("Connection get error", err)
		s.clientManager.RemoveClient(c.ID())
		//c.Disconnect()
	})

}

func NewService(c Config, d Diagnostic, dp Dispatcher) *Service {
	s := &Service{
		diag:          d,
		dispatcherBus: dp,
		ws: websocket.New(websocket.Config{
			WriteBufferSize: c.WriteBufferSize,
			ReadBufferSize:  c.ReadBufferSize,
			MaxMessageSize:  int64(c.WriteBufferSize),
			ReadTimeout:     websocket.DefaultWebsocketPongTimeout, //此作为readtimeout, 默认 如果有ping没有发送也成为read time out
		}),
		clientManager:   &WSClientManager{},
		dispatcherArray: []*utils.DispatchHandlerStruct{},
	}

	s.clientManager.Init()

	s.configValue.Store(c)

	return s

}

func (s *Service) createAndStartWebSocketNotifyDispatcher() error {
	if err := s.dispatcherBus.Create(dispatcherbus.DISPATCH_WS_NOTIFY, utils.DefaultDispatcherBufLen); err != nil {
		return err
	} else {
		return s.dispatcherBus.Start(dispatcherbus.DISPATCH_WS_NOTIFY)
	}
}

func (s *Service) postNotify(msg *DispatcherNotifyPackage) {
	if err := s.dispatcherBus.Dispatch(dispatcherbus.DISPATCH_WS_NOTIFY, msg); err != nil {
		s.diag.Error("notify", err)
	}
}

func (s *Service) addNewHttpHandler(r httpd.Route) {
	if s.Httpd == nil {
		return
	}
	h, err := s.Httpd.GetHandlerByName(httpd.BasePath)
	if err != nil {
		return
	}
	err = h.AddRoute(r)
	if err != nil {
		return
	}
}

func (s *Service) Open() error {

	c := s.Config()

	s.ws.OnConnection(s.onConnect) // 注册连接回调函数

	//s.HTTPD.Server.Get(c.Route, s.ws.Handler()) //将websocket 服务注册到get服务中

	r := httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_WS,
		Method:      "GET",
		Pattern:     c.Route,
		HandlerFunc: s.ws.Handler(),
	}
	s.addNewHttpHandler(r)

	if err := s.createAndStartWebSocketNotifyDispatcher(); err != nil {
		s.diag.Error("createAndStartWebSocketNotifyDispatcher Error", err)
		return nil
	}

	return nil
}

func (s *Service) Close() error {
	s.diag.Close()

	s.clientManager.CloseAll()

	s.diag.Closed()

	return nil
}

// ws推送结果到指定控制器
//func (s *Service) WSSendResult(sn string, payload string) {
//
//	c, exist := s.clientManager.GetClient(sn)
//	if exist {
//		c.Emit(WS_EVENT_RESULT, payload)
//	}
//}

// ws推送MES指示到显示ping
//func (s *Service) WSSendMes(event string,sn string, payload string) {
//
//	c, exist := s.clientManager.GetClient(sn)
//	if exist {
//		c.Emit(event, payload)
//	}
//}

func (s *Service) NotifyAll(evt string, payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(evt, payload)
}

// ws群发控制器状态
//func (s *Service) WSSendControllerStatus(payload string) {
//	s.clientManager.NotifyALL(WS_EVENT_CONTROLLER, payload)
//}

// ws群发控制器套筒状态
//func (s *Service) WSSendControllerSelectorStatus(payload string) {
//	s.clientManager.NotifyALL(WS_EVENT_SELECTOR, payload)
//}

// ws群发job选择信息
//func (s *Service) WSSendJob(payload string) {
//	s.clientManager.NotifyALL(WS_EVENT_JOB, payload)
//}

//func (s *Service) WSSendIOInput(payload string) {
//	if s == nil || s.clientManager == nil {
//		return
//	}
//	s.clientManager.NotifyALL(WS_EVENT_IO, payload)
//}

//func (s *Service) WSSendIOContact(sn string, ioType string, contactStatus string) {
//	wsMsg := GenerateWSMsg(0, WS_IO_CONTACT, WSReaderUID{
//		UID: uid,
//	})
//
//	s.clientManager.NotifyALL(WS_EVENT_IO, payload)
//}

//func (s *Service) WSSendReader(uid string) {
//	wsMsg := GenerateWSMsg(0, WS_READER_UID, WSReaderUID{
//		UID: uid,
//	})
//
//	body, _ := json.Marshal(wsMsg)
//	s.NotifyAll(WS_EVENT_READER, string(body))
//}

//func (s *Service) WSSendBarcode(src string, sn string, barcode string) {
//	wsMsg := GenerateWSMsg(0, WS_READER_UID, WSScannerRead{
//		Src: src,
//		SN: sn,
//		Barcode: barcode,
//	})
//
//	body, _ := json.Marshal(wsMsg)
//	s.NotifyAll(WS_EVENT_SCANNER, string(body))
//}

//func (s *Service) WSSendTightening(payload string) {
//	if s == nil || s.clientManager == nil {
//		return
//	}
//	s.clientManager.NotifyALL(WS_EVENT_TIGHTENING_DEVICE, payload)
//}

//func (s *Service) WSSendReply(reply *WSMsg) {
//	if s == nil || s.clientManager == nil {
//		return
//	}
//
//	body, _ := json.Marshal(reply)
//	s.clientManager.NotifyALL(WS_EVENT_REPLY, string(body))
//}

//func (s *Service) WSSendOrder(payload []byte) {
//	if s == nil || s.clientManager == nil {
//		return
//	}
//
//	s.clientManager.NotifyALL(WS_EVENT_ORDER, string(payload))
//}

func (s *Service) WSTestRecv(evt string, payload string) {
	go s.postNotify(&DispatcherNotifyPackage{nil, []byte(payload)})
}

func GenerateReply(sn uint64, wsType string, result int, msg string) *WSMsg {
	return &WSMsg{
		SN:   sn,
		Type: wsType,
		Data: WSReply{
			Result: result,
			Msg:    msg,
		},
	}
}

func GenerateWSMsg(sn uint64, wsType string, data interface{}) *WSMsg {
	return &WSMsg{
		SN:   sn,
		Type: wsType,
		Data: data,
	}
}

func WSClientSend(c websocket.Connection, event string, payload interface{}) error {
	if c == nil {
		return errors.New("conn is nil")
	}

	//str, _ := json.Marshal(payload)
	//fmt.Println(string(str))
	return c.Emit(event, payload)
}
