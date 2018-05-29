// Provides a server type for starting and configuring a Kapacitor server.
package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"

	"github.com/masami10/aiis/services/diagnostic"
	"github.com/masami10/aiis/keyvalue"
	"github.com/masami10/aiis/services/httpd"

	"github.com/pkg/errors"
)

// BuildInfo represents the build details for the server code.
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

// Server represents a container for the metadata and storage data and services.
// It is built using a Config and it manages the startup and shutdown of all
// services in the proper order.
type Server struct {
	dataDir  string
	hostname string

	config *Config

	err chan error

	Commander command.Commander

	HTTPDService          *httpd.Service

	SessionService        *diagnostic.SessionService


	// List of services in startup order
	Services []Service
	// Map of service name to index in Services list
	ServicesByName map[string]int

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
		config:           c,
		err:              make(chan error),
		DiagService:      diagService,
		Diag:             d,
		ServicesByName:   make(map[string]int),
		Commander:        c.Commander,
	}


	// Append Kapacitor services.
	s.initHTTPDService()


	s.appendSessionService()

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


func (s *Server) initHTTPDService() {
	d := s.DiagService.NewHTTPDHandler()
	srv := httpd.NewService(s.config.HTTP, s.hostname, d)

	srv.Handler.DiagService = s.DiagService

	srv.Handler.Version = s.BuildInfo.Version

	s.HTTPDService = srv
	s.TaskMaster.HTTPDService = srv
}

func (s *Server) appendHTTPDService() {
	s.AppendService("httpd", s.HTTPDService)
}

func (s *Server) appendSessionService() {
	srv := s.DiagService.SessionService
	srv.HTTPDService = s.HTTPDService

	s.AppendService("session", srv)
}


// Err returns an error channel that multiplexes all out of band errors received from all services.
func (s *Server) Err() <-chan error { return s.err }

// Open opens all the services.
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
	}
	s.err <- err
}

// Close shuts down the meta and data stores and all services.
func (s *Server) Close() error {

	// Close all services that write points first.
	if err := s.HTTPDService.Close(); err != nil {
		s.Diag.Error("error closing httpd service", err)
	}

	// Close services now that all tasks are stopped.
	for i := len(s.Services) - 1; i >= 0; i-- {
		service := s.Services[i]
		s.Diag.Debug("closing service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
		err := service.Close()
		if err != nil {
			s.Diag.Error("error closing service", err, keyvalue.KV("service", fmt.Sprintf("%T", service)))
		}
		s.Diag.Debug("closed service", keyvalue.KV("service", fmt.Sprintf("%T", service)))
	}

	// Finally close the task master
	return nil
}


func (s *Server) Reload() {
	return
}


// Service represents a service attached to the server.
type Service interface {
	Open() error
	Close() error
}

type tcpaddr struct{ host string }

func (a *tcpaddr) Network() string { return "tcp" }
func (a *tcpaddr) String() string  { return a.host }

