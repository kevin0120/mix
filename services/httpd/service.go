package httpd

import (
	"net/url"
	"fmt"
	"github.com/kataras/iris"
	"time"
	"log"
	stdContext "context"
	"strings"
)

type Diagnostic interface {
	NewHTTPServerErrorLogger() *log.Logger

	StartingService()
	StoppedService()
	ShutdownTimeout()
	AuthenticationEnabled(enabled bool)

	ListeningOn(addr string, proto string)

	WriteBodyReceived(body string)

	HTTP(
		host string,
		username string,
		start time.Time,
		method string,
		uri string,
		proto string,
		status int,
		referer string,
		userAgent string,
		reqID string,
		duration time.Duration,
	)

	Error(msg string, err error)
	RecoveryError(
		msg string,
		err string,
		host string,
		username string,
		start time.Time,
		method string,
		uri string,
		proto string,
		status int,
		referer string,
		userAgent string,
		reqID string,
		duration time.Duration,
	)
}

type Service struct {
	addr  string
	key   string
	err   chan error

	Handler *Handler
	shutdownTimeout time.Duration
	externalURL string
	server *iris.Application

	stop            chan chan struct{}

	diag                  Diagnostic
}

func NewService(c Config, hostname string, d Diagnostic) *Service {

	port, _ := c.Port()
	u := url.URL{
		Host:   fmt.Sprintf("%s:%d", hostname, port),
		Scheme: "http",
	}
	s := &Service{
		addr:        c.BindAddress,
		externalURL: u.String(),
		err:         make(chan error, 1),
		shutdownTimeout: time.Duration(c.ShutdownTimeout),
		Handler: NewHandler(
			c.LogEnabled,
			c.WriteTracing,
			d,
		),
	}

	return s
}

func (s *Service) manage() {

	select {
	case <-s.stop:
		// if we're already all empty, we're already done
		timeout := s.shutdownTimeout
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}


// Close closes the underlying listener.
func (s *Service) Close() error {
	defer s.diag.StoppedService()
	// If server is not set we were never started
	if s.server == nil {
		return nil
	}
	// Signal to manage loop we are stopping
	stopping := make(chan struct{})
	s.stop <- stopping

	<-stopping
	s.server = nil
	return nil
}

func (s *Service) serve() {
	err := s.server.Run(s.Addr())
	// The listener was closed so exit
	// See https://github.com/golang/go/issues/4373
	if !strings.Contains(err.Error(), "closed") {
		s.err <- fmt.Errorf("listener failed: addr=%s, err=%s", s.Addr(), err)
	} else {
		s.err <- nil
	}
}

// Open starts the service
func (s *Service) Open() error {
	s.diag.StartingService()

	s.server = iris.New()

	s.stop = make(chan chan struct{})

	go s.manage()
	go s.serve()
	return nil
}

func (s *Service) Addr() iris.Runner {
	return iris.Addr(s.addr)
}

func (s *Service) Err() <-chan error {
	return s.err
}

func (s *Service) URL() string {

	return "http://" + s.server.ConfigurationReadOnly().GetVHost()
}

// URL that should resolve externally to the server HTTP endpoint.
// It is possible that the URL does not resolve correctly  if the hostname config setting is incorrect.
func (s *Service) ExternalURL() string {
	return s.externalURL
}

