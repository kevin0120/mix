package utils

import (
	"github.com/pkg/errors"
	"log"
	"sync"
)

const (
	DEFAULT_BUF_LEN = 1024
)

type DispatchHandler func(data interface{})

// bufLen: 缓冲长度
func CreateDispatcher(bufLen int) *Dispatcher {
	return &Dispatcher{
		open:        false,
		buf:         make(chan interface{}, bufLen),
		closing:     make(chan struct{}, 1),
		dispatchers: []DispatchHandler{},
	}
}

type Dispatcher struct {
	mtx     sync.Mutex
	open    bool
	buf     chan interface{}
	closing chan struct{}

	dispatchers []DispatchHandler
}

func (s *Dispatcher) getOpen() bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	ret := s.open
	return ret
}

func (s *Dispatcher) setOpen(open bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.open = open
}

func (s *Dispatcher) Start() {
	if !s.getOpen() {
		go s.manage()
		s.setOpen(true)
	}
	return
}

func (s *Dispatcher) Release() {
	if s.getOpen() {
		s.closing <- struct{}{}
	}
}

func (s *Dispatcher) Register(dispatcher DispatchHandler) {
	s.dispatchers = append(s.dispatchers, dispatcher)
}

func (s *Dispatcher) Dispatch(data interface{}) error{
	if !s.getOpen() {
		msg := "Dispatcher Is Not Opened!!!"
		log.Fatalf(msg)
		return errors.New(msg)
	}
	s.buf <- data
	return nil
}

//todo: 限制协程的数量
func (s *Dispatcher) doDispatch(data interface{}) {
	for _, v := range s.dispatchers {
		go v(data)
	}
}

func (s *Dispatcher) manage() {
	for {
		select {
		case data := <-s.buf:
			s.doDispatch(data)

		case <-s.closing:
			s.setOpen(false)
			log.Fatalf("Dispatcher Is Closed!!!")
			return
		}
	}
}
