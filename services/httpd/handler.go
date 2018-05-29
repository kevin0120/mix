package httpd

import "github.com/influxdata/kapacitor/auth"

const (
	// Root path for the API
	BasePath = "/aiis/v1"
	// Root path for the preview API
	BasePreviewPath = "/aiis/v1preview"
	// Name of the special user for subscriptions
	SubscriptionUser = "~subscriber"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc interface{}
}

type Handler struct {
	requireAuthentication bool
	exposePprof           bool
	sharedSecret          string

	allowGzip bool

	Version string

	AuthService auth.Interface

	DiagService interface {
		SetLogLevelFromName(lvl string) error
	}

	diag Diagnostic
	// Detailed logging of write path
	// Uses normal logger
	writeTrace bool

	// Log every HTTP access.
	loggingEnabled bool
}

func NewHandler(loggingEnabled bool, writeTrace bool, d Diagnostic ) *Handler {
	h := &Handler{
		diag:                  d,
		writeTrace:            writeTrace,
		loggingEnabled:        loggingEnabled,
	}
	return h
}
