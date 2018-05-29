package httpd

import (
	"log"
	"time"
	"net"
	"net/http"
	"sync"
	"net/url"
	"fmt"
	stdContext "context"

	"github.com/kataras/iris"
)

type Diagnostic interface {
	NewHTTPServerErrorLogger() *log.Logger

	StartingService()
	StoppedService()
	ShutdownTimeout()

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
	ln    net.Listener
	addr  string
	err   chan error

	externalURL string

	server *http.Server
	mu     sync.Mutex
	wg     sync.WaitGroup

	new             chan net.Conn
	active          chan net.Conn
	idle            chan net.Conn
	closed          chan net.Conn
	stop            chan chan struct{}
	shutdownTimeout time.Duration

	Handler *Handler

	diag                  Diagnostic
	httpServerErrorLogger *log.Logger
}


func NewService(c Config, hostname string, d Diagnostic) *Service {

	port, _ := c.Port()
	u := url.URL{
		Host:   fmt.Sprintf("%s:%d", hostname, port),
		Scheme: "http",
	}

	s := &Service{
		addr:            c.BindAddress,
		externalURL:     u.String(),
		err:             make(chan error, 1),
		shutdownTimeout: time.Duration(c.ShutdownTimeout),
		Handler: NewHandler(
			c.LogEnabled,
			c.WriteTracing,
			d,
		),
		diag: d,
		httpServerErrorLogger: d.NewHTTPServerErrorLogger(),
	}

	return s
}

func (s *Service) Open() error {
	s.diag.StartingService()
	app := iris.New()

	// Method:    GET
	// Resource:  http://localhost:8080
	app.Get("/", func(ctx iris.Context) {
		// Bind: {{.message}} with "Hello world!"
		ctx.ViewData("message", "Hello world!")
		// Render template file: ./views/hello.html
		ctx.View("hello.html")
	})

	// Method:    GET
	// Resource:  http://localhost:8080/user/42
	//
	// Need to use a custom regexp instead?
	// Easy,
	// just mark the parameter's type to 'string'
	// which accepts anything and make use of
	// its `regexp` macro function, i.e:
	// app.Get("/user/{id:string regexp(^[0-9]+$)}")
	app.Get("/user/{id:long}", func(ctx iris.Context) {
		userID, _ := ctx.Params().GetInt64("id")
		ctx.Writef("User ID: %d", userID)
	})

	// Start the server using a network address.
	app.Run(iris.Addr(s.addr))
}

func (s *Service) Close() error {
	s.diag.StoppedService()

	timeout := s.shutdownTimeout
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
	defer cancel()
	app.Shutdown(ctx)
}

func (s *Service) AddRoutes(routes []Route) error {
	return s.Handler.AddRoutes(routes)
}