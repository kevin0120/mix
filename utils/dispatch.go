package utils

const (
	DEFAULT_BUF_LEN = 1024
)

type DispatchHandler func(data interface{})

// bufLen: 缓冲长度
func CreateDispatch(bufLen int) *Dispatch {
	return &Dispatch{
		buf:         make(chan interface{}, bufLen),
		closing:     make(chan struct{}, 1),
		dispatchers: []DispatchHandler{},
	}
}

type Dispatch struct {
	buf     chan interface{}
	closing chan struct{}

	dispatchers []DispatchHandler
}

func (s *Dispatch) Start() error {
	go s.manage()

	return nil
}

func (s *Dispatch) Release() {
	s.closing <- struct{}{}
}

func (s *Dispatch) Regist(dispatcher DispatchHandler) {
	s.dispatchers = append(s.dispatchers, dispatcher)
}

func (s *Dispatch) Dispatch(data interface{}) {
	s.buf <- data
}

func (s *Dispatch) doDispatch(data interface{}) {
	for _, v := range s.dispatchers {
		v(data)
	}
}

func (s *Dispatch) manage() {
	for {
		select {
		case data := <-s.buf:
			s.doDispatch(data)

		case _ = <-s.closing:
			return
		}
	}
}
