package httpd

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

//const (
//	// Root path for the API
//	BasePath = "/rush/v1"
//	// Root path for the preview API
//	BasePreviewPath = "/rush/v1preview"
//	// Name of the special user for subscriptions
//	SubscriptionUser = "~subscriber"
//)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc context.Handler
}

type Handler struct {
	requireAuthentication bool
	exposePprof           bool
	sharedSecret          string
	allowGzip             bool

	Version string

	DiagService interface {
		SetLogLevelFromName(lvl string) error
	}

	diag Diagnostic
	// Detailed logging of write path
	// Uses normal logger
	writeTrace bool

	party *iris.Party

	// Log every HTTP access.
	loggingEnabled bool
}

func NewHandler(loggingEnabled bool, writeTrace bool, d Diagnostic) *Handler {
	h := &Handler{
		diag:           d,
		writeTrace:     writeTrace,
		loggingEnabled: loggingEnabled,
	}
	return h
}

func (h *Handler) AddRoute(r Route) error {
	if len(r.Pattern) > 0 && r.Pattern[0] != '/' {
		return fmt.Errorf("route patterns must begin with a '/' %s", r.Pattern)
	}
	(*h.party).Handle(r.Method, r.Pattern, r.HandlerFunc)
	return nil
}

func (h *Handler) SetParty(p *iris.Party) error {
	h.party = p
	return nil
}
