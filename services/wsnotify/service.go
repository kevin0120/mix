package wsnotify

import (
	"github.com/kataras/iris/websocket"
	"sync"
	"github.com/masami10/rush/services/httpd"
	"sync/atomic"
	"encoding/json"
	"fmt"
)


const (
	WS_EVENT_STATUS    = "status"
	WS_EVENT_RESULT    = "result"
	WS_EVENT_WORKORDER = "workorder"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	Close()
	Closed()
}

type Service struct {
	configValue atomic.Value
	diag Diagnostic

	ws *websocket.Server

	conn map[websocket.Connection] string

	mutex *sync.Mutex //为了连接进行锁操作

	Httpd  *httpd.Service

}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}


type WSRegist struct {
	HMI_SN string `json:"hmi_sn"`
}


func (s *Service) onConnect(c websocket.Connection) {

	c.OnMessage(func(data []byte) {
		reg := WSRegist{}
		err := json.Unmarshal(data, &reg)
		if err != nil {
			Msg := map[string]string{"msg":"regist msg error"}
			msg, _ := json.Marshal(Msg)
			c.EmitMessage(msg)
		}

		s.mutex.Lock()
		conn := s.conn
		s.mutex.Unlock()

		v := make([]string, 0, len(conn))

		for  _, value := range conn {
			v = append(v, value)
		}

		for _, value := range v {
			if value == reg.HMI_SN {
				Msg := fmt.Sprintf("client with sn:%s already exists", reg.HMI_SN)
				_Msg := map[string]string{"msg":Msg}
				reg_str, _ := json.Marshal(_Msg)
				c.EmitMessage(reg_str)
				return
			}
		}

		//加入连接队列
		s.mutex.Lock()
		s.conn[c] = reg.HMI_SN
		s.mutex.Unlock()

	})

	c.OnDisconnect(func() {
		s.mutex.Lock()
		delete(s.conn, c)
		s.mutex.Unlock()
		s.diag.Disconnect(c.ID())
	})

}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag:  d,
		ws:    websocket.New(websocket.Config{WriteBufferSize: c.WriteBufferSize, ReadBufferSize: c.ReadBufferSize}),
		conn:  make(map[websocket.Connection]string),
		mutex: new(sync.Mutex),
	}

	s.configValue.Store(c)

	return s

}

func (s *Service) Open() error {

	c := s.Config()

	s.ws.OnConnection(s.onConnect) // 注册连接回调函数

	s.Httpd.Server.Get(c.Route, s.ws.Handler()) //将websocket 服务注册到get服务中

	return nil

}

func (s *Service) Close() error {
	s.diag.Close()
	s.mutex.Lock()
	conn := s.conn
	s.mutex.Unlock()
	for c  := range conn {
		c.Disconnect() // 此方法会调用OnDisconnect回调
	}

	s.diag.Closed()

	return nil
}
