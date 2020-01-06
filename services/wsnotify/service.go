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
	WS_EVENT_SERVICE           = "service"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	OnMessage(msg string)
	Close()
	Closed()
}

type OnNewClient func(c websocket.Connection)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic

	ws *websocket.Server

	Httpd HTTPService

	clientManager *WSClientManager

	OnNewClient OnNewClient

	dispatcherBus Dispatcher
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) GetWorkCenter() string {
	c := s.Config()
	return c.Workcenter
}

func (s *Service) NewWebSocketRecvHandler(handler func(interface{})) {
	fn := utils.CreateDispatchHandlerStruct(handler)
	s.dispatcherBus.Register(dispatcherbus.DISPATCHER_WS_NOTIFY, fn)
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
		clientManager: &WSClientManager{},
	}

	s.clientManager.Init()

	s.configValue.Store(c)

	return s

}

func (s *Service) createAndStartWebSocketNotifyDispatcher() error {
	if err := s.dispatcherBus.Create(dispatcherbus.DISPATCHER_WS_NOTIFY, utils.DefaultDispatcherBufLen); err != nil {
		return err
	} else {
		return s.dispatcherBus.Start(dispatcherbus.DISPATCHER_WS_NOTIFY)
	}
}

func (s *Service) postNotify(msg *DispatcherNotifyPackage) {
	if err := s.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_WS_NOTIFY, msg); err != nil {
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

func (s *Service) NotifyAll(evt string, payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(evt, payload)
}

func (s *Service) WSTestRecv(evt string, payload string) {
	go s.postNotify(&DispatcherNotifyPackage{nil, []byte(payload)})
}

func GenerateReply(sn uint64, wsType string, result int, msg string) *WSMsg {
	return GenerateWSMsg(sn, wsType, WSReply{
		Result: result,
		Msg:    msg,
	})
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

	return c.Emit(event, payload)
}
