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
	Debug(msg string)
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
	wsnotify.WSNotify
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

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "POST",
		Pattern:     "/ak2",
		HandlerFunc: s.methods.postAK2,
	}
	s.Httpd.Handler[0].AddRoute(r)

	//r = httpd.Route{
	//	RouteType:   httpd.ROUTE_TYPE_HTTP,
	//	Method:      "PUT",
	//	Pattern:     "/test-protocol",
	//	HandlerFunc: s.methods.testProtocol,
	//}
	//s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/ws-test",
		HandlerFunc: s.methods.wsTest,
	}
	s.Httpd.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/storage",
		HandlerFunc: s.methods.getStorage,
	}
	s.Httpd.Handler[0].AddRoute(r)

	s.ODOO.WS.AddNotify(s)

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

		inputs := openprotocol.IOMonitor{
			ControllerSN: k,
			Inputs:       v.Inputs(),
		}

		msg, _ = json.Marshal(inputs)
		c.Emit(wsnotify.WS_EVENT_IO, string(msg))

		// 推送工具状态
		// TODO:工具状态
		//tools := v.Tools()
		//for sn, status := range tools {
		//	ts := wsnotify.WSToolStatus{
		//		ToolSN: sn,
		//		Status: status,
		//	}
		//	msg, _ = json.Marshal(ts)
		//
		//	c.Emit(wsnotify.WS_EVETN_TOOL, string(msg))
		//}
	}

	odooStatus := wsnotify.WSOdooStatus{
		Status: s.ODOO.Status(),
	}

	str, _ := json.Marshal(odooStatus)
	c.Emit(wsnotify.WS_EVENT_ODOO, string(str))
}

func (s *Service) OnWSMsg(c websocket.Connection, data []byte) {
	msg := wsnotify.WSMsg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(string(data), err)
		return
	}

	sData, _ := json.Marshal(msg.Data)
	switch msg.Type {
	case WS_ORDER_LIST:
		// 请求获取工单列表
		workorders, err := s.DB.Workorders("")
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			return
		}

		body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, workorders))

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))

	case WS_ORDER_DETAIL:
		// 请求获取工单详情
		orderReq := WSOrderReq{}
		err := json.Unmarshal(sData, &orderReq)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			return
		}

		w, err := s.DB.WorkorderStep(orderReq.ID)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
			return
		}

		body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, w))
		//fmt.Println(string(body))
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_ORDER, string(body))

	case WS_ORDER_UPDATE:
		// 更新工单状态
		orderReq := WSOrderReq{}
		err := json.Unmarshal(sData, &orderReq)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			return
		}

		_, err = s.DB.UpdateWorkorder(&storage.Workorders{
			Id:     orderReq.ID,
			Status: orderReq.Status,
		})

		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
			return
		}

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))

	case WS_ORDER_STEP_UPDATE:
		// 更新工步状态

		orderReq := WSOrderReq{}
		err := json.Unmarshal(sData, &orderReq)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			return
		}

		_, err = s.DB.UpdateStep(&storage.Steps{
			Id:     orderReq.ID,
			Status: orderReq.Status,
		})

		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
			return
		}

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
	}
}
