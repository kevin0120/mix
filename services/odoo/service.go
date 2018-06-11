package odoo

import (
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/httpd"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	diag        Diagnostic
	methods	Methods
	DB	   *storage.Service
	Httpd  *httpd.Service
}

func NewService(d Diagnostic) *Service {
	s := &Service{
		diag:      d,
		methods: Methods{},
	}

	s.methods.service = s

	return s
}

func (s *Service) Open() error {
	var r httpd.Route

	r = httpd.Route{
		RouteType:	httpd.ROUTE_TYPE_HTTP,
		Method:  "GET",
		Pattern: "/results",
		HandlerFunc: s.methods.getResults,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:	httpd.ROUTE_TYPE_HTTP,
		Method:  "PATCH",
		Pattern: "/results/{id:int}",
		HandlerFunc: s.methods.patchResult,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:	httpd.ROUTE_TYPE_HTTP,
		Method:  "POST",
		Pattern: "/workorders",
		HandlerFunc: s.methods.postWorkorders,
	}
	s.Httpd.Handler[0].AddRoute(r)

	return nil
}

func (s *Service) Close() error {
	return nil
}
