package ts002

import (
	"github.com/kataras/iris/context"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"sync/atomic"
	"time"
)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	httpd       IHttpService

	client *resty.Client
}

func NewService(c Config, d Diagnostic, h IHttpService) *Service {
	ss := &Service{
		diag:  d,
		httpd: h,
	}

	ss.configValue.Store(c)

	return ss
}

func (s *Service) ensureHttpClient() *resty.Client {
	if s.client != nil {
		return s.client
	}
	if client, err := utils.CreateRetryClient(10*time.Second, 5); err != nil {
		s.diag.Error("ensureHttpClient", err)
		return nil
	} else {
		s.client = client
		return client
	}
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {

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
	//todo: 增加所有ts002 接口
	//s.addNewHTTPHandler("")
	return nil
}

func (s *Service) SendNFCData() error {
	//todo: 发送读卡器数据
	//client := s.ensureHttpClient()
	//r := client.R()
	return nil
}
// 收到拧紧结果
func (s Service) onTighteningResult(data interface{}) {
	if data == nil {
		return
	}

}
