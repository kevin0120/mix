package server

import (
	"fmt"
	"github.com/masami10/rush/command"
	"github.com/masami10/rush/keyvalue"
	"github.com/masami10/rush/services/aiis"
	"github.com/masami10/rush/services/diagnostic"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/minio"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/controller"
)

type BuildInfo struct {
	Version  string
	Commit   string
	Branch   string
	Platform string
}

type Diagnostic interface {
	Debug(msg string, ctx ...keyvalue.T)
	Info(msg string, ctx ...keyvalue.T)
	Error(msg string, err error, ctx ...keyvalue.T)
}

// Service represents a service attached to the server.
type Service interface {
	Open() error
	Close() error
}

type Server struct {
	dataDir  string
	hostname string

	HTTPDService *httpd.Service

	AudiVWService *audi_vw.Service

	config       *Config
	// List of services in startup order
	Services []Service

	ServicesByName map[string]int
	err            chan error

	BuildInfo   BuildInfo
	Commander   command.Commander
	DiagService *diagnostic.Service
	Diag        Diagnostic
}

// New returns a new instance of Server built from a config.
func New(c *Config, buildInfo BuildInfo, diagService *diagnostic.Service) (*Server, error) {
	err := c.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %s. To generate a valid configuration file run `rushd config > rush.generated.conf`. ", err)
	}
	d := diagService.NewServerHandler()
	s := &Server{
		config:         c,
		DiagService:    diagService,
		Diag:           d,
		BuildInfo:      buildInfo,
		dataDir:        c.DataDir,
		hostname:       c.Hostname,
		ServicesByName: make(map[string]int),
		err:            make(chan error),
		Commander:      c.Commander,
	}

	if err := s.initHTTPDService(); err != nil {
		return nil, errors.Wrap(err, "init httpd service")
	}

	if err := s.initAudiVWDService(); err != nil {
		return nil, errors.Wrap(err, "init Audi/VW service")
	}

	s.appendMinioService()

	s.appendAiisService()

	s.appendOdooService()

	s.appendWebsocketService()

	s.appendStorageService()

	if err := s.appendControllersService(); err != nil {
		return nil, errors.Wrap(err, "Controllers service")
	}

	s.appendAudiVWService() //此服务必须在控制器服务后进行append

	s.appendHTTPDService()

	return s, nil
}

func (s *Server) AppendService(name string, srv Service) {
	if _, ok := s.ServicesByName[name]; ok {
		// Should be unreachable code
		panic("cannot append service twice")
	}
	i := len(s.Services)
	s.Services = append(s.Services, srv)
	s.ServicesByName[name] = i
}

func (s *Server) initHTTPDService() error {
	p := s.config.DocPath
	exist, err := utils.FileIsExist(p)
	if err != nil || !exist {
		return fmt.Errorf("API File %s is not exist", p)
	}
	d := s.DiagService.NewHTTPDHandler()
	srv := httpd.NewService(p, s.config.HTTP, s.hostname, d, s.DiagService)

	s.HTTPDService = srv

	return nil
}

func (s *Server) initAudiVWDService() error {
	c := s.config.AudiVW
	d := s.DiagService.NewHTTPDHandler()
	srv := audi_vw.NewService(c, d)

	s.AudiVWService = srv

	return nil
}

func (s *Server) appendAudiVWService() {
	s.AppendService("audi/vw", s.AudiVWService)
}

func (s *Server) appendHTTPDService() {
	s.AppendService("httpd", s.HTTPDService)
}

func (s *Server) appendMinioService() error {
	c := s.config.Minio
	d := s.DiagService.NewMinioHandler()
	srv := minio.NewService(c, d)

	s.AppendService("minio", srv)

	return nil
}


func (s *Server) appendControllersService() error {
	c := s.config.Contollers
	d := s.DiagService.NewControllerHandler()
	srv, err := controller.NewService(c, d, s.AudiVWService)

	if err != nil{
		return errors.Wrap(err, "append Controller service fail")
	}

	s.AppendService("controller", srv)

	return nil
}

func (s *Server) appendAiisService() error {
	c := s.config.Aiis
	d := s.DiagService.NewAiisHandler()
	srv := aiis.NewService(c, d)

	s.AppendService("aiis", srv)

	return nil
}

func (s *Server) appendOdooService() error {
	c := s.config.Odoo
	d := s.DiagService.NewOdooHandler()
	srv := odoo.NewService(c, d)

	s.AppendService("odoo", srv)

	return nil
}

func (s *Server) appendWebsocketService() error {
	c := s.config.Ws
	d := s.DiagService.NewWebsocketHandler()
	srv := wsnotify.NewService(c, d)

	srv.Httpd = s.HTTPDService //http 服务注入

	s.AppendService("websocket", srv)

	return nil
}


func (s *Server) appendStorageService() error {
	c := s.config.Storage
	d := s.DiagService.NewStorageHandler()
	srv := storage.NewService(c, d)

	s.AppendService("storage", srv)

	return nil
}

func (s *Server) Open() error {

	if err := s.startServices(); err != nil {
		s.Close()
		return err
	}

	go s.watchServices()

	return nil
}

func (s *Server) startServices() error {
	for _, service := range s.Services {
		s.Diag.Debug("opening service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
		if err := service.Open(); err != nil {
			return fmt.Errorf("open service %T: %s", service, err)
		}
		s.Diag.Debug("opened service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
	}

	return nil
}

// Watch if something dies
func (s *Server) watchServices() {
	var err error
	select {
	case err = <-s.HTTPDService.Err():
	case err = <-s.AudiVWService.Err():
	}
	s.err <- err
}

// Err returns an error channel that multiplexes all out of band errors received from all services.
func (s *Server) Err() <-chan error { return s.err }

func (s *Server) Reload() {
	return
}

func (s *Server) Close() error {

	if err := s.HTTPDService.Close(); err != nil {
		s.Diag.Error("error closing httpd service", err)
	}

	for i := len(s.Services) - 1; i >= 0; i-- {
		service := s.Services[i]
		s.Diag.Debug("closing service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
		err := service.Close()
		if err != nil {
			s.Diag.Error("error closing service", err, keyvalue.KV("service", fmt.Sprintf("%T", service)))
		}
		s.Diag.Debug("closed service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
	}

	return nil
}
