package controller

import (
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"sync"
	"fmt"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Controller interface {
	Start()
	Close() error
	Status() string
	Protocol() string
	LoadController(controller *storage.Controllers)
}

type Protocol interface {
	Parse(msg string) ([]byte, error)
	Write(sn string, buf []byte) error
	AddNewController(cfg Config) Controller
}

type HandlerPackage struct {
	Result interface{}
	Curve interface{}
}

type Service struct {
	configs []Config

	protocols   map[string]Protocol //进行服务注入, serial_no : Protocol
	Controllers map[string]Controller

	DB            *storage.Service
	WS            *wsnotify.Service
	Aiis          *aiis.Service
	Minio         *minio.Service

	handlers      Handlers
	wg            sync.WaitGroup
	closing       chan struct{}

	handle_buffer chan HandlerPackage

	diag Diagnostic
}

func NewService(cs Configs, d Diagnostic, pAudi Protocol, pOpenprotocol Protocol) (*Service, error) {
	s := &Service{
		configs:     cs,
		diag:        d,
		Controllers: map[string]Controller{},
		protocols:   map[string]Protocol{},
		handlers: Handlers{},
		handle_buffer:make(chan HandlerPackage, 1024),
		wg:       sync.WaitGroup{},
		closing:  make(chan struct{}, 1),
	}

	s.handlers.controllerService = s

	for _, c := range cs {
		switch c.Protocol {
		case AUDIPROTOCOL:
			new_controller := pAudi.AddNewController(c)
			s.Controllers[c.SN] = new_controller
			s.protocols[c.SN] = pAudi

		case OPENPROTOCOL:
			new_controller := pOpenprotocol.AddNewController(c)
			s.Controllers[c.SN] = new_controller
			s.protocols[c.SN] = pOpenprotocol

		default:

		}
	}

	return s, nil
}

func (s *Service) InitControllers() {
	for k, v := range s.Controllers {
		c, err := s.DB.CreateController(k)
		if err == nil {
			v.LoadController(&c)
		}
	}
}

func (s *Service) Write(serialNo string, buf []byte) error {

	controller := s.protocols[serialNo]

	return controller.Write(serialNo, buf)
}

func (s *Service) Open() error {
	s.handlers.Init(4)

	for i := 0; i < 4; i++ {
		s.wg.Add(1)
		go s.HandleProcess()
		s.diag.Debug(fmt.Sprintf("init handle process:%d", i+1))
	}

	s.InitControllers()

	return nil
}

func (s *Service) Handle(result interface{}, curve interface{}) {
	pkg := HandlerPackage{
		Result: result,
		Curve: curve,
	}

	s.handle_buffer <- pkg
}

func (s *Service) HandleProcess() {
	for {
		select {
		case msg := <-s.handle_buffer:
			s.handlers.Handle(msg.Result, msg.Curve)

		case <-s.closing:
			s.wg.Done()
			return
		}
	}

}

func (s *Service) Close() error {
	return nil
}
