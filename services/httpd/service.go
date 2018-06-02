package httpd

import (
	stdContext "context"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/masami10/rush/services/diagnostic"
	"log"
	"net/url"
	"strings"
	"time"
	"io/ioutil"
)

const (
	// Root path for the API
	BasePath = "/rush/v1"
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
	addr string
	key  string
	err  chan error

	Handler         []*Handler
	shutdownTimeout time.Duration
	externalURL     string
	server          *iris.Application

	stop chan chan struct{}

	HandlerByNames map[string]int

	cors CorsConfig

	diag Diagnostic

	DiagService interface {
		SetLogLevelFromName(lvl string) error
	}

	httpServerErrorLogger *log.Logger
}

func NewService(doc string, c Config, hostname string, d Diagnostic, disc *diagnostic.Service) *Service {

	port, _ := c.Port()
	u := url.URL{
		Host:   fmt.Sprintf("%s:%d", hostname, port),
		Scheme: "http",
	}
	s := &Service{
		addr:                  c.BindAddress,
		externalURL:           u.String(),
		cors:                  c.Cors,
		server:                iris.New(),
		err:                   make(chan error, 1),
		HandlerByNames:        make(map[string]int),
		shutdownTimeout:       time.Duration(c.ShutdownTimeout),
		diag:                  d,
		DiagService:           disc,
		httpServerErrorLogger: d.NewHTTPServerErrorLogger(),
	}
	s.AddNewHandler(BasePath, c, d, disc)

	r := Route{
		Method:  "GET",
		Pattern: "/healthz",
		HandlerFunc: func(ctx iris.Context) {
			ctx.StatusCode(iris.StatusNoContent)
		},
	}
	s.Handler[0].AddRoute(r)

	r1 := Route{
		Method:  "GET",
		Pattern: "/doc",
		HandlerFunc: func(ctx iris.Context) {
			f, _ := ioutil.ReadFile(doc)
			ctx.Write(f)
			ctx.StatusCode(iris.StatusOK)
		},
	}

	s.Handler[0].AddRoute(r1)


	return s
}

func (s *Service) manage() {
	//println("start mamager")
	var stopDone chan struct{}
	select {
	case stopDone = <-s.stop:
		// if we're already all empty, we're already done
		timeout := s.shutdownTimeout
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		defer close(stopDone)
		s.server.Shutdown(ctx)
		return
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
	err := s.server.Run(s.Addr(), iris.WithoutInterruptHandler)
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

func (s *Service) GetHandlerByName(version string) (*Handler, error) {
	i, ok := s.HandlerByNames[version]
	if !ok {
		// Should be unreachable code
		return nil, fmt.Errorf("cannot get handler By %s", version)
	}

	return s.Handler[i], nil
}

func (s *Service) AddNewHandler(version string, c Config, d Diagnostic, disc *diagnostic.Service) error {
	if _, ok := s.HandlerByNames[version]; ok {
		// Should be unreachable code
		panic("cannot append handler twice")
	}
	crs := cors.New(cors.Options{
		AllowedOrigins:   s.cors.AllowedOrigins,
		AllowCredentials: s.cors.AllowCredentials,
	})
	p := s.server.Party(version, crs).AllowMethods(iris.MethodOptions)
	if p == nil {
		return fmt.Errorf("fail to create the party%s", version)
	}
	h := NewHandler(
		c.LogEnabled,
		c.WriteTracing,
		d,
	)
	h.DiagService = disc
	h.Version = version
	h.party = &p

	i := len(s.Handler)
	s.Handler = append(s.Handler, h)

	s.HandlerByNames[version] = i

	return nil
}
