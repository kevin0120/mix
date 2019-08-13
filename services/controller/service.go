package controller

import (
	"fmt"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"sync"
)

const (
	HANDLE_TYPE_RESULT = "result"
	HANDLE_TYPE_CURVE  = "curve"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Controller interface {
	Start()
	Close() error
	Protocol() string
	Inputs() string
	LoadController(controller *storage.Controllers)
	Tools() map[string]string
	device.Device
}

type Protocol interface {
	Parse(msg string) ([]byte, error)
	Write(sn string, buf []byte) error
	AddNewController(cfg ControllerConfig) Controller
	AddDevice(cfg DeviceConfig, ts interface{}) Controller
}

type HandlerPackage struct {
	result interface{}
	curve  interface{}
}

type Service struct {
	config Config

	protocols   map[string]Protocol //进行服务注入, serial_no : Protocol
	Controllers map[string]Controller

	DB    *storage.Service
	WS    *wsnotify.Service
	Aiis  *aiis.Service
	Minio *minio.Service

	Handlers Handlers
	wg       sync.WaitGroup
	closing  chan struct{}

	handle_buffer chan *HandlerPackage

	diag Diagnostic
}

func NewService(cs Config, d Diagnostic, pAudi Protocol, pOpenprotocol Protocol) (*Service, error) {
	s := &Service{
		config:        cs,
		diag:          d,
		Controllers:   map[string]Controller{},
		protocols:     map[string]Protocol{},
		Handlers:      Handlers{},
		handle_buffer: make(chan *HandlerPackage, 1024),
		wg:            sync.WaitGroup{},
		closing:       make(chan struct{}, 1),
	}

	s.Handlers.controllerService = s

	for _, c := range cs.Configs {
		switch c.Protocol {
		case AUDIPROTOCOL:
			newController := pAudi.AddNewController(c)
			s.Controllers[c.SN] = newController
			s.protocols[c.SN] = pAudi

		case OPENPROTOCOL:
			newController := pOpenprotocol.AddNewController(c)
			s.Controllers[c.SN] = newController
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
	if !s.config.Enable {
		return nil
	}

	s.Handlers.Init(s.config.Workers)

	for i := 0; i < s.config.Workers; i++ {
		s.wg.Add(1)
		go s.HandleProcess()
		s.diag.Debug(fmt.Sprintf("init handle process:%d", i+1))
	}

	s.InitControllers()

	return nil
}

func (s *Service) Handle(result interface{}, curve interface{}) {
	pkg := HandlerPackage{
		result: result,
		curve:  curve,
	}

	s.handle_buffer <- &pkg
}

func (s *Service) HandleProcess() {
	for {
		select {
		case data := <-s.handle_buffer:
			s.Handlers.Handle(data.result, data.curve)
		case <-s.closing:
			s.wg.Done()
			return
		}
	}
}

func (s *Service) Close() error {
	for i := 0; i < s.config.Workers; i++ {
		s.closing <- struct{}{}
		s.diag.Debug(fmt.Sprintf("stop handler process:%d", i+1))
	}
	return nil
}

func (s *Service) GetControllers() *map[string]Controller {
	return &s.Controllers
}
