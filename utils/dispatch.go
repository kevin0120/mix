package utils

const (
	DEFAULT_BUF_LEN = 1024
)

type DispatchHandler func(data interface{})

// bufLen: 缓冲长度
func CreateDispatcher(bufLen int) *Dispatcher {
	return &Dispatcher{
		buf:         make(chan interface{}, bufLen),
		closing:     make(chan struct{}, 1),
		dispatchers: []DispatchHandler{},
	}
}

type Dispatcher struct {
	buf     chan interface{}
	closing chan struct{}

	dispatchers []DispatchHandler
}

func (s *Dispatcher) Start() error {
	go s.manage()

	return nil
}

func (s *Dispatcher) Release() {
	s.closing <- struct{}{}
}

func (s *Dispatcher) Regist(dispatcher DispatchHandler) {
	s.dispatchers = append(s.dispatchers, dispatcher)
}

func (s *Dispatcher) Dispatch(data interface{}) {
	s.buf <- data
}

func (s *Dispatcher) doDispatch(data interface{}) {
	for _, v := range s.dispatchers {
		v(data)
	}
}

func (s *Dispatcher) manage() {
	for {
		select {
		case data := <-s.buf:
			s.doDispatch(data)

		case _ = <-s.closing:
			return
		}
	}
}
