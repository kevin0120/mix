package utils

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"log"
	"sync"
)

const (
	DefaultDispatcherBufLen = 1024
)

type dispatchHandler func(data interface{})

type DispatchHandlerStruct struct {
	ID      string
	Handler dispatchHandler
}

func CreateDispatchHandlerStruct(h dispatchHandler) *DispatchHandlerStruct {
	if h == nil {
		return nil
	}
	return &DispatchHandlerStruct{
		ID:      uuid.NewV4().String(),
		Handler: h,
	}
}

// bufLen: 缓冲长度
func CreateDispatcher(bufLen int) *Dispatcher {
	return &Dispatcher{
		open:        false,
		buf:         make(chan interface{}, bufLen),
		closing:     make(chan struct{}, 1),
		dispatchers: map[string]dispatchHandler{},
	}
}

type Dispatcher struct {
	mtx     sync.Mutex
	open    bool
	buf     chan interface{}
	closing chan struct{}

	dispatchers map[string]dispatchHandler
}

func (s *Dispatcher) removeHandler(handler string) {
	if _, ok := s.dispatchers[handler]; ok {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		delete(s.dispatchers, handler)
	}
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

func (s *Dispatcher) Release(handler string) {
	if !s.getOpen() {
		return
	}
	s.removeHandler(handler)
	if len(s.dispatchers) == 0 {
		s.closing <- struct{}{}
	}
}

func (s *Dispatcher) Register(key string, dispatcher dispatchHandler) string {
	if key == "" {
		key = uuid.NewV4().String()
	}
	s.dispatchers[key] = dispatcher
	return key
}

func (s *Dispatcher) Dispatch(data interface{}) error {
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
