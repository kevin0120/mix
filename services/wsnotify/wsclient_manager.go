package wsnotify

import (
	"github.com/kataras/iris/websocket"
	"sync"
)

type WSClientManager struct {
	conn  map[string]websocket.Connection // sn : Connection
	mutex sync.Mutex                      //为了连接进行锁操作
}

func (s *WSClientManager) Init() {
	s.conn = map[string]websocket.Connection{}
	s.mutex = sync.Mutex{}
}

func (s *WSClientManager) AddClient(sn string, c websocket.Connection) {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	s.conn[sn] = c
}

func (s *WSClientManager) RemoveClientBySN(sn string) {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	if sn != "" {
		delete(s.conn, sn)
	}
}

func (s *WSClientManager) RemoveClient(cid string) {
	defer s.mutex.Unlock()

	s.mutex.Lock()

	var key string = ""
	for k, v := range s.conn {
		if v.ID() == cid {
			key = k
			break
		}
	}

	if key != "" {
		delete(s.conn, key)
	}
}

func (s *WSClientManager) NotifyALL(evt string, payload string) {
	defer s.mutex.Unlock()

	s.mutex.Lock()

	for _, v := range s.conn {
		_ = v.Emit(evt, payload)
	}
}

func (s *WSClientManager) CloseAll() {
	defer s.mutex.Unlock()

	s.mutex.Lock()

	for _, v := range s.conn {
		v.Disconnect()
	}
}

func (s *WSClientManager) GetClient(sn string) (websocket.Connection, bool) {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	v, e := s.conn[sn]
	return v, e
}
