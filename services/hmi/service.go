package hmi

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/audi_vw"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/odoo"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
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

	s.ControllerService.WS.OnNewClient = s.OnNewHmiClient

	var r httpd.Route

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/tool-enable",
		HandlerFunc: s.methods.putToolControl,
	}
	s.Httpd.Handler[0].AddRoute(r)

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
		Pattern:     "/psets-manual",
		HandlerFunc: s.methods.putManualPSets,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/psets",
		HandlerFunc: s.methods.getPSetList,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/pset-detail",
		HandlerFunc: s.methods.getPSetDetail,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/jobs",
		HandlerFunc: s.methods.getJobList,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/job-detail",
		HandlerFunc: s.methods.getJobDetail,
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
		Pattern:     "/jobs-manual",
		HandlerFunc: s.methods.putManualJobs,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/controller-mode",
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
		Pattern:     "/next-workorder",
		HandlerFunc: s.methods.getNextWorkorder,
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

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/io-sets",
		HandlerFunc: s.methods.putIOSet,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/scanner",
		HandlerFunc: s.methods.putBarcodeTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/input-test",
		HandlerFunc: s.methods.putIOInputTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/result-test",
		HandlerFunc: s.methods.putResultTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/job-control",
		HandlerFunc: s.methods.putJobControll,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/operation/{code:string}",
		HandlerFunc: s.methods.getRoutingOpertions,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/local-results",
		HandlerFunc: s.methods.getLocalResults,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/workorders",
		HandlerFunc: s.methods.listWorkorders,
	}
	s.Httpd.Handler[0].AddRoute(r)

	return nil

}

func (s *Service) Close() error {
	s.diag.Close()
	s.diag.Closed()

	return nil
}

func (s *Service) resetResult(id int64) error {
	// 重置结果
	err := s.DB.ResetResult(id)
	if err != nil {
		return err
	}

	// 删除对应曲线
	err = s.DB.DeleteCurvesByResult(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) OnNewHmiClient(c websocket.Connection) {
	controllers := s.ControllerService.GetControllers()

	for k, v := range *controllers {
		s := wsnotify.WSStatus{
			SN:     k,
			Status: string(v.Status()),
		}

		msg, _ := json.Marshal(s)

		c.Emit(wsnotify.WS_EVENT_CONTROLLER, string(msg))
	}

	odooStatus := wsnotify.WSOdooStatus{
		Status: s.ODOO.Status(),
	}

	str, _ := json.Marshal(odooStatus)
	c.Emit(wsnotify.WS_EVENT_ODOO, string(str))
}
