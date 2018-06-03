package wsnotify

import (
	"github.com/kataras/iris/websocket"
	"sync"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	Close()
	Closed()
}

type Service struct {
	diag Diagnostic

	ws *websocket.Server

	conn map[websocket.Connection]bool

	mutex *sync.Mutex //为了连接进行锁操作

}

func (s *Service) onConnect(c websocket.Connection) {

	s.mutex.Lock()
	s.conn[c] = true
	s.mutex.Unlock()

	c.OnDisconnect(func() {
		s.mutex.Lock()
		delete(s.conn, c)
		s.mutex.Unlock()
		s.diag.Disconnect(c.ID())
	})

}

func NewService(c Config, d Diagnostic) *Service {

	return &Service{
		diag:  d,
		ws:    websocket.New(websocket.Config{WriteBufferSize: c.WriteBufferSize, ReadBufferSize: c.ReadBufferSize}),
		conn:  make(map[websocket.Connection]bool),
		mutex: new(sync.Mutex),
	}

}

func (s *Service) Open() error {

	s.ws.OnConnection(s.onConnect) // 注册连接回调函数

	return nil

}

func (s *Service) Close() error {
	s.diag.Close()
	s.mutex.Lock()
	conn := s.conn
	s.mutex.Unlock()
	for c := range conn {
		c.Disconnect() // 此方法会调用OnDisconnect回调
	}

	s.diag.Closed()

	return nil
}
