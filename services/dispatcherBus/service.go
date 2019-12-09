package dispatcherBus

import (
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/utils"
	"sync"
	"sync/atomic"
)

// name: handlerName
type DispatcherMap map[string]utils.DispatchHandler

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag        Diagnostic
	configValue atomic.Value

	dispatchers    map[string]*utils.Dispatcher
	dispatchersMtx sync.Mutex

	registers    map[string][]utils.DispatchHandler
	registersMtx sync.Mutex
}

func NewService(c Config, d Diagnostic) (*Service, error) {

	srv := &Service{
		diag:           d,
		dispatchers:    map[string]*utils.Dispatcher{},
		dispatchersMtx: sync.Mutex{},

		registers:    map[string][]utils.DispatchHandler{},
		registersMtx: sync.Mutex{},
	}

	srv.configValue.Store(c)

	return srv, nil
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	for _, v := range s.dispatchers {
		if v != nil {
			v.Release()
		}
	}

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) getRegistersAndRemove(name string) ([]utils.DispatchHandler, error) {
	s.registersMtx.Lock()
	defer s.registersMtx.Unlock()

	r, exist := s.registers[name]
	if !exist {
		return nil, errors.New("Registers Not Found")
	}

	delete(s.registers, name)

	return r, nil
}

func (s *Service) appendRegisters(name string, handler utils.DispatchHandler) {
	s.registersMtx.Lock()
	defer s.registersMtx.Unlock()

	r, exist := s.registers[name]
	if !exist {
		// 不存在则创建
		s.registers[name] = []utils.DispatchHandler{handler}
		return
	}

	// 存在则追加
	r = append(r, handler)
	s.registers[name] = r
}

func (s *Service) getDispatcher(name string) (*utils.Dispatcher, error) {
	s.dispatchersMtx.Lock()
	defer s.dispatchersMtx.Unlock()

	d, exist := s.dispatchers[name]
	if !exist {
		return nil, errors.New("Dispatcher Not Found")
	}

	return d, nil
}

func (s *Service) Create(name string, len int) error {
	s.dispatchersMtx.Lock()
	defer s.dispatchersMtx.Unlock()

	_, exist := s.dispatchers[name]
	if exist {
		return errors.New("Dispatcher Already Exist")
	}

	s.dispatchers[name] = utils.CreateDispatcher(len)

	// 追加注册handler
	registers, err := s.getRegistersAndRemove(name)
	if err == nil {
		for _, handler := range registers {
			s.dispatchers[name].Register(handler)
		}
	}

	return nil
}

func (s *Service) Start(name string) error {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		return err
	}

	dispatcher.Start()

	return nil
}

func (s *Service) Release(name string) error {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		return err
	}

	dispatcher.Release()

	return nil
}

func (s *Service) Register(name string, handler utils.DispatchHandler) {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		// 如果dispatcher还没创建， 将handler加入注册列表等待创建后注册
		s.appendRegisters(name, handler)
		return
	}

	dispatcher.Register(handler)
}

func (s *Service) Dispatch(name string, data interface{}) error {
	dispatcher, err := s.getDispatcher(name)
	if err != nil {
		return err
	}

	return dispatcher.Dispatch(data)
}

// create, register and start
func (s *Service) LaunchDispatchersByHandlerMap(dispatcherMap DispatcherMap) {
	for name, handler := range dispatcherMap {
		err := s.Create(name, utils.DEFAULT_BUF_LEN)
		if err != nil {
			s.diag.Debug(err.Error())
		}
		s.Register(name, handler)
		s.Start(name)
	}
}
