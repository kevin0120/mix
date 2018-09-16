package server

import (
	"fmt"
	"github.com/masami10/aiis/command"
	"github.com/masami10/aiis/keyvalue"
	"github.com/masami10/aiis/services/changan"
	"github.com/masami10/aiis/services/diagnostic"
	"github.com/masami10/aiis/services/fis"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/odoo"
	"github.com/masami10/aiis/services/pmon"
	"github.com/masami10/aiis/services/rush"
	"github.com/masami10/aiis/services/storage"
	"github.com/masami10/aiis/services/wsnotify"
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
	PmonService  *pmon.Service
	FisService   *fis.Service
	OdooService  *odoo.Service

	StorageService  *storage.Service
	WSNotifyService *wsnotify.Service

	config *Config
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
		return nil, fmt.Errorf("invalid configuration: %s. To generate a valid configuration file run `kapacitord config > kapacitor.generated.conf`.", err)
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

	s.initHTTPDService()

	s.appendWebsocketService()

	err = s.initPMONService()
	if err != nil {
		return nil, err
	}

	s.appendStorageService()

	s.appendOdooService()

	s.appendPmonService()

	s.appendFisService()

	s.appendChanganService()

	s.appendRushService() //必须后于http & storage

	s.appendHTTPDService()

	return s, nil
}

func (s *Server) appendRushService() {
	d := s.DiagService.NewRushHandler()
	srv := rush.NewService(s.config.Rush, d)
	srv.HTTPDService = s.HTTPDService
	srv.StorageService = s.StorageService
	srv.Fis = s.FisService

	s.AppendService("rush", srv)
}

func (s *Server) appendStorageService() {
	d := s.DiagService.NewStorageHandler()
	srv := storage.NewService(s.config.Storage, d)

	s.StorageService = srv
	s.AppendService("storage", srv)
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
	d := s.DiagService.NewHTTPDHandler()
	srv := httpd.NewService(s.config.HTTP, s.hostname, d, s.DiagService)

	s.HTTPDService = srv

	return nil
}

func (s *Server) initPMONService() error {
	d := s.DiagService.NewPmonHandler()
	srv, err := pmon.NewService(s.config.Pmon, d)

	if err != nil {
		return err
	}

	srv.HTTPD = s.HTTPDService //httpd服务注入

	s.PmonService = srv

	return nil
}

func (s *Server) appendHTTPDService() {
	s.AppendService("httpd", s.HTTPDService)
}

func (s *Server) appendPmonService() error {
	s.AppendService("pmon", s.PmonService)

	return nil
}

func (s *Server) appendWebsocketService() error {
	c := s.config.Ws
	d := s.DiagService.NewWebsocketHandler()
	srv := wsnotify.NewService(c, d)

	srv.Httpd = s.HTTPDService //http 服务注入

	s.WSNotifyService = srv
	s.AppendService("websocket", srv)

	return nil
}

func (s *Server) appendFisService() error {
	d := s.DiagService.NewFisHandler()
	c := s.config.Fis

	srv := fis.NewService(d, c, s.PmonService)

	if c.Enable {
		srv.Odoo = s.OdooService
		s.FisService = srv
		s.AppendService("fis", srv)
	}

	return nil
}

func (s *Server) appendChanganService() error {
	d := s.DiagService.NewChanganHandler()
	c := s.config.Changan

	srv := changan.NewService(d, c, s.WSNotifyService)
	s.AppendService("changan", srv)
	return nil
}

func (s *Server) appendOdooService() {
	c := s.config.Odoo
	d := s.DiagService.NewOdooHandler()
	srv := odoo.NewService(c, d)

	s.OdooService = srv
	s.AppendService("odoo", srv)

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

// Err returns an error channel that multiplexes all out of band errors received from all services.
func (s *Server) Err() <-chan error { return s.err }

func (s *Server) Reload() {
	return
}

// Watch if something dies
func (s *Server) watchServices() {
	var err error
	select {
	case err = <-s.HTTPDService.Err():
	}
	s.err <- err
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
