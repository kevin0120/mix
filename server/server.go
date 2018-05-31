package server

import (
	"fmt"
	"github.com/masami10/aiis/command"
	"github.com/masami10/aiis/keyvalue"
	"github.com/masami10/aiis/services/diagnostic"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/pmon"
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

	s.appendPmonService()

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
	d := s.DiagService.NewHTTPDHandler()
	srv := httpd.NewService(s.config.HTTP, s.hostname, d, s.DiagService)

	s.HTTPDService = srv

	return nil
}

func (s *Server) appendHTTPDService() {
	s.AppendService("httpd", s.HTTPDService)
}

func (s *Server) appendPmonService() error {
	c := s.config.Pmon
	d := s.DiagService.NewPmonHandler()
	srv := pmon.NewService(c, d)

	srv.HTTPD = s.HTTPDService //httpd服务注入

	s.AppendService("pmon", srv)

	return nil
}

func (s *Server) Open() error {

	if err := s.startServices(); err != nil {
		s.Close()
		return err
	}

	return nil
}

func (s *Server) startServices() error {
	for _, service := range s.Services {
		if err := service.Open(); err != nil {
			return fmt.Errorf("open service %T: %s", service, err)
		}
	}

	return nil
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
