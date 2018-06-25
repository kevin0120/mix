package wsnotify

import (
	"github.com/kataras/iris/websocket"

	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/httpd"
	"sync/atomic"
)

const (
	WS_EVENT_STATUS = "status"
	WS_EVENT_RESULT = "result"
	WS_EVENT_REG    = "regist"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	Close()
	Closed()
}

type Service struct {
	configValue atomic.Value
	diag        Diagnostic

	ws *websocket.Server

	Httpd *httpd.Service

	clientManager WSClientManager
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) onConnect(c websocket.Connection) {

	c.OnMessage(func(data []byte) {
		reg := WSRegist{}
		err := json.Unmarshal(data, &reg)
		if err != nil {
			Msg := map[string]string{"msg": "regist msg error"}
			msg, err := json.Marshal(Msg)
			if err != nil {
				c.Emit(WS_EVENT_REG, msg)
			}
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
		} else {
			// 将客户端加入列表
			s.clientManager.AddClient(reg.HMI_SN, c)
			Msg := map[string]string{"msg": "OK"}
			msg, err := json.Marshal(Msg)
			if err != nil {
				c.Emit(WS_EVENT_REG, msg)
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

// ws推送结果到指定控制器
func (s *Service) WSSendResult(sn string, payload string) {
	c, exist := s.clientManager.GetClient(sn)
	if exist {
		c.Emit(WS_EVENT_RESULT, payload)
	}
}

// ws群发控制器状态
func (s *Service) WSSendControllerStatus(payload string) {
	s.clientManager.NotifyALL(WS_EVENT_STATUS, payload)
}
