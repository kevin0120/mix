package wsnotify

import (
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/websocket"
	"sync"

	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/httpd"
	"sync/atomic"
)

const (
	WS_EVENT_CONTROLLER        = "controller"
	WS_EVENT_RESULT            = "result"
	WS_EVENT_REG               = "regist"
	WS_EVENT_SELECTOR          = "selector"
	WS_EVENT_JOB               = "job"
	WS_EVENT_SCANNER           = "scanner"
	WS_EVENT_IO                = "io"
	WS_EVENT_ODOO              = "odoo"
	WS_EVENT_MAINTENANCE       = "maintenance"
	WS_EVETN_TOOL              = "tool"
	WS_EVENT_READER            = "reader"
	WS_EVENT_TIGHTENING_DEVICE = "tightening_device"
	WS_EVENT_REPLY             = "reply"
	WS_EVENT_DEVICE            = "device"
	WS_EVENT_ORDER             = "order"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	OnMessage(msg string)
	Close()
	Closed()
}

type OnNewClient func(c websocket.Connection)

type WSNotify interface {
	OnWSMsg(c websocket.Connection, data []byte)
}

type Service struct {
	configValue atomic.Value
	diag        Diagnostic

	ws *websocket.Server

	Httpd *httpd.Service

	clientManager *WSClientManager
	msgBuffer     chan string

	OnNewClient OnNewClient

	notifies    []WSNotify
	mtxNotifies sync.Mutex
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) onConnect(c websocket.Connection) {

	c.OnMessage(func(data []byte) {

		s.diag.OnMessage(string(data))
		msg := WSMsg{}
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return
		}

		if msg.Type == WS_REG {
			reg := WSRegist{}
			strData, _ := json.Marshal(msg.Data)
			err := json.Unmarshal([]byte(strData), &reg)
			if err != nil {
				c.Disconnect()
			}

			_, exist := s.clientManager.GetClient(reg.HMI_SN)
			if exist {
				Msg := fmt.Sprintf("client with sn:%s already exists", reg.HMI_SN)
				msgs := map[string]string{"msg": Msg}
				regStrs, err := json.Marshal(msgs)
				if err != nil {
					c.Emit(WS_EVENT_REG, regStrs)
				}

				c.Disconnect()
			} else {
				// 将客户端加入列表
				s.clientManager.AddClient(reg.HMI_SN, c)
				Msg := map[string]string{"msg": "OK"}
				msg, err := json.Marshal(Msg)
				if err != nil {
					c.Emit(WS_EVENT_REG, msg)
				}

				if s.OnNewClient != nil {
					s.OnNewClient(c)
				}
			}
		} else {
			go s.notify(c, data)
		}
	})

	c.OnDisconnect(func() {
		s.clientManager.RemoveClient(c.ID())
		s.diag.Disconnect(c.ID())
	})

	c.OnError(func(err error) {
		s.diag.Error("Connection get error", err)
		c.Disconnect()
	})

}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
		ws: websocket.New(websocket.Config{
			WriteBufferSize: c.WriteBufferSize,
			ReadBufferSize:  c.ReadBufferSize,
			ReadTimeout:     websocket.DefaultWebsocketPongTimeout, //此作为readtimeout, 默认 如果有ping没有发送也成为read time out
		}),
		clientManager: &WSClientManager{},
		msgBuffer:     make(chan string, 1024),
		notifies:      []WSNotify{},
		mtxNotifies:   sync.Mutex{},
	}

	s.clientManager.Init()

	s.configValue.Store(c)

	return s

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
	s.Httpd.Handler[0].AddRoute(r)

	return nil

}

func (s *Service) Close() error {
	s.diag.Close()

	s.clientManager.CloseAll()

	s.diag.Closed()

	return nil
}

func (s *Service) sendProcess() {

}

func (s *Service) AddNotify(n WSNotify) {
	defer s.mtxNotifies.Unlock()
	s.mtxNotifies.Lock()

	s.notifies = append(s.notifies, n)
}

func (s *Service) notify(c websocket.Connection, data []byte) {
	defer s.mtxNotifies.Unlock()
	s.mtxNotifies.Lock()

	for _, v := range s.notifies {
		v.OnWSMsg(c, data)
	}
}

// ws推送结果到指定控制器
func (s *Service) WSSendResult(sn string, payload string) {

	c, exist := s.clientManager.GetClient(sn)
	if exist {
		c.Emit(WS_EVENT_RESULT, payload)
	}
}

func (s *Service) WSSend(evt string, payload string) {
	s.clientManager.NotifyALL(evt, payload)
}

// ws群发控制器状态
func (s *Service) WSSendControllerStatus(payload string) {
	s.clientManager.NotifyALL(WS_EVENT_CONTROLLER, payload)
}

// ws群发控制器套筒状态
func (s *Service) WSSendControllerSelectorStatus(payload string) {
	s.clientManager.NotifyALL(WS_EVENT_SELECTOR, payload)
}

// ws群发job选择信息
func (s *Service) WSSendJob(payload string) {
	s.clientManager.NotifyALL(WS_EVENT_JOB, payload)
}

// ws群发扫码信息
func (s *Service) WSSendScanner(payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(WS_EVENT_SCANNER, payload)
}

func (s *Service) WSSendIOInput(payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(WS_EVENT_IO, payload)
}

func (s *Service) WSSendIO(payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(WS_EVENT_IO, payload)
}

func (s *Service) WSSendReader(payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(WS_EVENT_READER, payload)
}

func (s *Service) WSSendTightening(payload string) {
	if s == nil || s.clientManager == nil {
		return
	}
	s.clientManager.NotifyALL(WS_EVENT_TIGHTENING_DEVICE, payload)
}

func (s *Service) WSSendReply(reply *WSMsg) {
	if s == nil || s.clientManager == nil {
		return
	}

	body, _ := json.Marshal(reply)
	s.clientManager.NotifyALL(WS_EVENT_REPLY, string(body))
}

func (s *Service) WSSendOrder(payload []byte) {
	if s == nil || s.clientManager == nil {
		return
	}

	s.clientManager.NotifyALL(WS_EVENT_ORDER, string(payload))
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

func GenerateResult(sn uint64, wsType string, data interface{}) *WSMsg {
	return &WSMsg{
		SN:   sn,
		Type: wsType,
		Data: data,
	}
}

func (s *Service) WSTestRecv(evt string, payload string) {
	go s.notify(nil, []byte(payload))
}

func WSClientSend(c websocket.Connection, event string, payload interface{}) error {
	if c == nil {
		return errors.New("conn is nil")
	}

	return c.Emit(event, payload)
}
