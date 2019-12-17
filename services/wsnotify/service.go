package wsnotify

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/utils"
	"go.uber.org/atomic"
)

const (
	WS_EVENT_TIGHTENING  = "tightening"
	WS_EVENT_SCANNER     = "scanner"
	WS_EVENT_IO          = "io"
	WS_EVENT_MAINTENANCE = "maintenance"
	WS_EVENT_READER      = "reader"
	WS_EVENT_REPLY       = "reply"
	WS_EVENT_DEVICE      = "device"
	WS_EVENT_ORDER       = "order"
	WS_EVENT_SERVICE     = "service"

	WS_EVENT_ERROR = "err"
)

const (
	WS_TYPE_ERROR = "new_error"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	OnMessage(msg string)
	Close()
	Closed()
}

type Service struct {
	configValue   atomic.Value
	diag          Diagnostic
	ws            *websocket.Server
	httpd         HTTPService
	clientManager *WSClientManager
	dispatcherBus Dispatcher
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) NewWebSocketRecvHandler(handler func(interface{})) {
	fn := utils.CreateDispatchHandlerStruct(handler)
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, fn)
}

func (s *Service) onConnect(c websocket.Connection) {

	c.OnMessage(func(data []byte) {
		msg := WSMsg{}
		err := json.Unmarshal(data, &msg)
		if err != nil {
			s.diag.Error("WSMsg Payload Error", err)
			return
		}

		if msg.Type == WsReg {
			s.handleRegister(&msg, c)
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

func NewService(c Config, d Diagnostic, dp Dispatcher, httpd HTTPService) *Service {
	s := &Service{
		diag:          d,
		dispatcherBus: dp,
		ws: websocket.New(websocket.Config{
			WriteBufferSize: c.WriteBufferSize,
			ReadBufferSize:  c.ReadBufferSize,
			MaxMessageSize:  int64(c.WriteBufferSize),
			ReadTimeout:     websocket.DefaultWebsocketPongTimeout, //此作为readtimeout, 默认 如果有ping没有发送也成为read time out
		}),
		clientManager: &WSClientManager{
			diag: d,
		},
		httpd: httpd,
	}

	s.clientManager.Init()

	s.configValue.Store(c)

	return s

}

func (s *Service) handleRegister(msg *WSMsg, c websocket.Connection) {
	var reg WSRegist
	strData, _ := json.Marshal(msg.Data)
	err := json.Unmarshal(strData, &reg)
	if err != nil {
		_ = c.Disconnect()
		s.clientManager.RemoveClientBySN(reg.HMISn)

		_ = WSClientSend(c, WS_EVENT_REPLY, GenerateReply(msg.SN, msg.Type, -1, err.Error()))
	}

	_, exist := s.clientManager.GetClient(reg.HMISn)
	if exist {
		Msg := fmt.Sprintf("Client With SN:%s Already Exists", reg.HMISn)
		_ = c.Disconnect()
		s.clientManager.RemoveClientBySN(reg.HMISn)

		_ = WSClientSend(c, WS_EVENT_REPLY, GenerateReply(msg.SN, msg.Type, -2, Msg))
	} else {
		// 将客户端加入列表
		s.clientManager.AddClient(reg.HMISn, c)

		// 注册成功
		_ = WSClientSend(c, WS_EVENT_REPLY, GenerateReply(msg.SN, msg.Type, 0, ""))
	}
}

func (s *Service) createAndStartWebSocketNotifyDispatcher() error {
	if err := s.dispatcherBus.Create(dispatcherbus.DispatcherWsNotify, utils.DefaultDispatcherBufLen); err != nil {
		return err
	} else {
		return s.dispatcherBus.Start(dispatcherbus.DispatcherWsNotify)
	}
}

func (s *Service) postNotify(msg *DispatcherNotifyPackage) {
	if err := s.dispatcherBus.Dispatch(dispatcherbus.DispatcherWsNotify, msg); err != nil {
		s.diag.Error("notify", err)
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
	s.httpd.AddNewHttpHandler(r)

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
