package ts002

import (
	"github.com/kataras/iris/context"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/pkg/errors"
	"sync/atomic"
)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	httpd       IHttpService

	IO               *io.Service
	TighteningDevice *tightening_device.Service
}

func NewService(c Config, d Diagnostic, h IHttpService) *Service {
	ss := &Service{
		diag:  d,
		httpd: h,
	}

	ss.configValue.Store(c)

	return ss
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	s.addTS002HTTPHandlers()

	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) getDefaultHandler() (*httpd.Handler, error) {

	if s.httpd == nil {
		return nil, errors.New("Please Inject Http Service First")
	}
	return s.httpd.GetHandlerByName(httpd.BasePath)
}

func (s *Service) addNewHTTPHandler(method, pattern string, handler context.Handler) error {

	if h, err := s.getDefaultHandler(); err != nil {
		s.diag.Error("addNewHandler", err)
		return err
	} else {
		r := Route{
			RouteType:   httpd.ROUTE_TYPE_HTTP,
			Method:      method,
			Pattern:     pattern,
			HandlerFunc: handler,
		}
		return h.AddRoute(r)
	}

}

func (s *Service) addTS002HTTPHandlers() error {

	s.addNewHTTPHandler("PUT", "/alarm", s.putAlarmReq)
	s.addNewHTTPHandler("PUT", "/pset", s.putPSetReq)
	s.addNewHTTPHandler("PUT", "/io", s.putIOReq)

	return nil
}

// 报警控制
func (s *Service) alarmControl(req *RushAlarmReq) error {
	if req == nil {
		return errors.New("alarmControl: Req Is Nil")
	}

	return nil
}

// PSet下发控制
func (s *Service) psetControl(req *RushPSetReq) error {
	if req == nil {
		return errors.New("psetControl: Req Is Nil")
	}

	return s.TighteningDevice.Api.ToolPSetByIP(&tightening_device.PSetSet{
		WorkorderID: 0,
		PSet:        req.PSet,
		IP:          req.ToolID,
		PointID:     req.PointID,
	})
}

// IO输出控制请求
func (s *Service) ioControl(req *RushIOControlReq) error {
	if req == nil {
		return errors.New("ioControl: Req Is Nil")
	}

	return nil
}

// 收到拧紧结果
func (s *Service) onTighteningResult(data interface{}) {
	if data == nil {
		return
	}

	tighteningResult := data.(*tightening_device.TighteningResult)
	if tighteningResult.MeasureResult == tightening_device.RESULT_NOK {
		// 如果结果NOK，则触发报警
	}
}
