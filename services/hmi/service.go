package hmi

import (
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/storage"
)

type Diagnostic interface {
	Error(msg string, err error)
	Disconnect(id string)
	Close()
	Closed()
}

type Service struct {
	diag              Diagnostic
	methods           Methods
	DB                *storage.Service
	Httpd             *httpd.Service
	ODOO              *odoo.Service
	AudiVw            *audi_vw.Service
	OpenProtocol      *openprotocol.Service
	ControllerService *controller.Service
	SN                string
}

func NewService(d Diagnostic) *Service {

	s := &Service{
		diag:    d,
		methods: Methods{},
	}

	s.methods.service = s

	return s
}

func (s *Service) Open() error {

	var r httpd.Route

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/psets",
		HandlerFunc: s.methods.putPSets,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/jobs",
		HandlerFunc: s.methods.putJobs,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/enable-job-mode",
		HandlerFunc: s.methods.enableJobMode,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/workorder",
		HandlerFunc: s.methods.getWorkorder,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/controller-status",
		HandlerFunc: s.methods.getStatus,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/healthz",
		HandlerFunc: s.methods.getHealthz,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/hmi-results",
		HandlerFunc: s.methods.getHmiResults,
	}
	s.Httpd.Handler[0].AddRoute(r)

	return nil

}

func (s *Service) Close() error {
	s.diag.Close()
	s.diag.Closed()

	return nil
}
