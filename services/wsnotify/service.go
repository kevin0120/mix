package wsnotify

import (
	"github.com/kataras/iris/websocket"

	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/httpd"
	"sync/atomic"
)

const (
	WS_EVENT_CONTROLLER = "controller"
	WS_EVENT_RESULT     = "result"
	WS_EVENT_REG        = "regist"
	WS_EVENT_SELECTOR   = "selector"
	WS_EVENT_JOB        = "job"
	WS_EVENT_SCANNER    = "scanner"
	WS_EVENT_IO         = "io"
	WS_EVENT_ODOO       = "odoo"
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

	Httpd *httpd.Service

	clientManager WSClientManager
	msgBuffer     chan string

	OnNewClient OnNewClient
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) onConnect(c websocket.Connection) {

	c.OnMessage(func(data []byte) {
		s.diag.OnMessage(string(data))
		reg := WSRegist{}
		err := json.Unmarshal(data, &reg)
		if err != nil {
			Msg := map[string]string{"msg": "regist msg error"}
			msg, err := json.Marshal(Msg)
			if err != nil {
				c.Emit(WS_EVENT_REG, msg)
			}

			c.Disconnect()
			return
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
		clientManager: WSClientManager{},
		msgBuffer:     make(chan string, 1024),
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
	s.clientManager.NotifyALL(WS_EVENT_SCANNER, payload)
}

func (s *Service) WSSendIOInput(payload string) {
	s.clientManager.NotifyALL(WS_EVENT_IO, payload)
}
