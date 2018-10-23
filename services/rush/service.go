package rush

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/aiis/services/changan"
	"github.com/masami10/aiis/services/fis"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/wsnotify"
	"gopkg.in/resty.v1"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type cResult struct {
	r    *OperationResult
	id   int64
	ip   string
	port string
}

type RushResult struct {
	HasUpload bool `json:"has_upload"`
}

type Service struct {
	HTTPDService *httpd.Service
	workers      int
	Opened       bool
	wg           sync.WaitGroup
	chResult     chan cResult
	closing      chan struct{}
	configValue  atomic.Value
	httpClient   *resty.Client
	route        string

	ws            *websocket.Server
	clientManager wsnotify.WSClientManager

	StorageService interface {
		UpdateResults(result *OperationResult, id int64, sent int) error
	}

	Fis     *fis.Service
	Changan *changan.Service

	diag Diagnostic
}

func NewService(c Config, d Diagnostic) *Service {
	if c.Enable {
		s := Service{
			diag:     d,
			workers:  c.Workers,
			Opened:   false,
			chResult: make(chan cResult, c.Workers),
			route:    c.Route,

			ws: websocket.New(websocket.Config{
				WriteBufferSize: c.WSWriteBufferSize,
				ReadBufferSize:  c.WSReadBufferSize,
				ReadTimeout:     websocket.DefaultWebsocketPongTimeout, //此作为readtimeout, 默认 如果有ping没有发送也成为read time out
			}),
			clientManager: wsnotify.WSClientManager{},
		}

		s.clientManager.Init()
		s.configValue.Store(c)
		return &s
	}

	return nil
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {

	r := httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/operation.results/{result_id:long}",
		HandlerFunc: s.getResultUpdate,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/fis.results",
		HandlerFunc: s.putFisResult,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "GET",
		Pattern:     "/healthz",
		HandlerFunc: s.getHealthz,
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	s.ws.OnConnection(s.onConnect)
	r = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_WS,
		Method:      "GET",
		Pattern:     s.Config().WSRoute,
		HandlerFunc: s.ws.Handler(),
	}
	s.HTTPDService.Handler[0].AddRoute(r)

	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(time.Duration(s.Config().Timeout))
	client.SetContentLength(true)
	// Headers for all request
	client.SetHeaders(s.Config().Headers)
	client.
		SetRetryCount(s.Config().MaxRetry).
		SetRetryWaitTime(time.Duration(s.Config().PushInterval)).
		SetRetryMaxWaitTime(20 * time.Second)

	s.httpClient = client

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)

		go s.run()
	}

	s.Opened = true

	return nil
}

func (s *Service) onConnect(c websocket.Connection) {

	c.OnMessage(func(data []byte) {
		ws_msg := WSMsg{}
		err := json.Unmarshal(data, &ws_msg)
		if err != nil {
			s.diag.Error("ws error", err)
			return
		}

		str_data, _ := json.Marshal(ws_msg.Data)

		switch ws_msg.Type {
		case WS_REG:
			reg := WSRegist{}
			err := json.Unmarshal(str_data, &reg)
			if err != nil {
				Msg := map[string]string{"msg": "regist msg error"}
				msg, err := json.Marshal(Msg)
				if err != nil {
					c.Emit(wsnotify.WS_EVENT_REG, msg)
				}

				c.Disconnect()
				return
			}

			_, exist := s.clientManager.GetClient(reg.Rush_SN)
			if exist {
				Msg := fmt.Sprintf("client with sn:%s already exists", reg.Rush_SN)
				msgs := map[string]string{"msg": Msg}
				regStrs, err := json.Marshal(msgs)
				if err != nil {
					c.Emit(wsnotify.WS_EVENT_REG, regStrs)
				}

				c.Disconnect()
			} else {
				// 将客户端加入列表
				s.clientManager.AddClient(reg.Rush_SN, c)
				Msg := map[string]string{"msg": "OK"}
				msg, err := json.Marshal(Msg)
				if err != nil {
					c.Emit(wsnotify.WS_EVENT_REG, msg)
				}
			}

		case WS_RESULT:
			op_result := WSOpResult{}
			err := json.Unmarshal(str_data, &op_result)
			if err != nil {
				s.diag.Error("ws result error", err)
				return
			}

			rush_ip := strings.Split(c.Context().RemoteAddr(), ":")[0]
			if strings.Contains(rush_ip, ":") {
				kvs := strings.Split(rush_ip, ":")
				rush_ip = kvs[0]
			}

			cr := cResult{
				r:    &op_result.Result,
				id:   op_result.ResultID,
				ip:   rush_ip,
				port: op_result.Port,
			}

			s.chResult <- cr
		}

	})

	c.OnDisconnect(func() {
		s.clientManager.RemoveClient(c.ID())
	})

	c.OnError(func(err error) {
		s.diag.Error("Connection get error", err)
		c.Disconnect()
	})

}

func (s *Service) putFisResult(ctx iris.Context) {

	var r OperationResult
	err := ctx.ReadJSON(&r)

	if err != nil {
		ctx.Writef(fmt.Sprintf("Result Params from Odoo wrong: %s", err.Error()))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if s.Fis != nil {
		fis_result := s.OperationToFisResult(&r)
		fis_err := s.Fis.PushResult(&fis_result)
		if fis_err != nil {
			ctx.Writef(fmt.Sprintf("Push fis err: %s", fis_err.Error()))
			ctx.StatusCode(iris.StatusBadRequest)
		} else {
			ctx.StatusCode(iris.StatusOK)
		}
	}

}

func (s *Service) getResultUpdate(ctx iris.Context) {
	resultId, err := ctx.Params().GetInt64("result_id")

	if err != nil {
		ctx.Writef("error while trying to parse resultId parameter")
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	var r OperationResult
	err = ctx.ReadJSON(&r)

	if err != nil {
		ctx.Writef(fmt.Sprintf("Result Params from Rush wrong: %s", err))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	rush_port := ctx.GetHeader("rush_port")
	rush_ip := strings.Split(ctx.Request().RemoteAddr, ":")[0]
	if strings.Contains(rush_ip, ":") {
		kvs := strings.Split(rush_ip, ":")
		rush_ip = kvs[0]
	}

	cr := cResult{
		r:    &r,
		id:   resultId,
		ip:   rush_ip,
		port: rush_port,
	}

	s.chResult <- cr

	ctx.StatusCode(iris.StatusNoContent)
}

func (s *Service) getHealthz(ctx iris.Context) {

	ctx.StatusCode(iris.StatusNoContent)
	return
}

func (s *Service) run() {
	for {
		select {
		case r := <-s.chResult:
			s.HandleResult(&r)

		case <-s.closing:
			s.wg.Done()
			return
		}
	}
}

func (s *Service) Close() error {
	if !s.Opened {
		return nil
	}

	for i := 0; i < s.workers; i++ {
		ss := struct{}{}
		s.closing <- ss
	}

	s.wg.Wait()
	return nil
}

func (s *Service) OperationToFisResult(r *OperationResult) fis.FisResult {
	var result fis.FisResult
	result.Init()

	result.EquipemntName = r.EquipemntName
	result.FactoryName = r.FactoryName
	result.Year = r.Year
	result.Pin = r.Pin
	result.PinCheckCode = r.Pin_check_code
	result.AssemblyLine = r.AssemblyLine
	result.ResultID = fmt.Sprintf("%s-%s-%s-%02d", r.Model, s.Fis.Config().FactoryCode, r.NutNo, r.Seq)
	result.Lnr = r.Lnr

	valueResult := 1

	if strings.ToUpper(r.MeasureResult) == "OK" {
		result.ResultValue = "IO__"
	} else {
		result.ResultValue = "NIO_"
		valueResult = 0
	}

	result.Dat = r.ControlDate
	result.SystemType = s.Fis.Config().SystemType
	result.SoftwareVersion = s.Fis.Config().SoftwareVersion
	result.Mode = s.Fis.Config().Mode

	// 扭矩结果
	var rv fis.FisResultValue
	rv.Value = r.MeasureTorque
	rv.ID = fis.FIS_ID_NM
	rv.Unit = fis.FIS_UNIT_NM
	rv.Measure = valueResult
	result.Values = append(result.Values, rv)

	// 角度结果
	rv.Value = r.MeasureDegree
	rv.ID = fis.FIS_ID_DEG
	rv.Unit = fis.FIS_UNIT_DEG
	rv.Measure = valueResult
	result.Values = append(result.Values, rv)

	return result
}

func (s *Service) OperationToChanganResult(r *OperationResult) changan.TighteningResults {
	result := changan.TighteningResults{}

	result.Spent = 0
	result.Angle = r.MeasureDegree
	result.Torque = r.MeasureTorque
	result.Result = r.MeasureResult
	result.Mode = r.Mode
	result.AngleMax = r.PsetWMax
	result.AngleMin = r.PsetWMin
	result.AngleTarget = r.PsetWTarget
	result.Batch = r.Batch
	result.Cartype = r.Model
	result.ControllerSn = r.ControllerSN
	result.Exception = r.ExceptionReason
	result.Strategy = r.PsetStrategy
	result.TighteningId = r.TighteningId
	result.ToolSn = r.ToolSN
	result.TorqueMax = r.PsetMMax
	result.TorqueMin = r.PsetMMin
	result.TorqueTarget = r.PsetMTarget
	result.UpdateTime = r.ControlDate
	result.Vin = r.Vin
	result.WorkcenterCode = r.WorkcenterCode

	return result
}

func (s *Service) PatchResultFlag(result_id int64, has_upload bool, ip string, port string) error {
	if s.httpClient == nil {
		return errors.New("rush http client is nil")
	}

	rush_result := RushResult{}
	rush_result.HasUpload = has_upload
	r := s.httpClient.R().SetBody(rush_result)

	resp, err := r.Patch(fmt.Sprintf("http://%s%s%s/%d", ip, port, s.route, result_id))
	if err != nil {
		return fmt.Errorf("patch result flag failed: %s\n", err)
	} else {
		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("patch result flag failed: %s\n", resp.Status())
		}
	}

	return nil
}

func (s *Service) HandleResult(cr *cResult) {

	// 结果推送fis
	sent := 1
	if s.Fis != nil {
		fisResult := s.OperationToFisResult(cr.r)

		e := s.Fis.PushResult(&fisResult)
		if e != nil {
			sent = 0
			s.diag.Error("push result to fis error", e)
		}
	}

	if s.Changan != nil {
		changanResult := s.OperationToChanganResult(cr.r)

		if !s.Changan.AndonDB.InsertResult(&changanResult) {
			s.diag.Error("insert andon result failed", nil)
		}
	}

	// 结果保存数据库
	err := s.StorageService.UpdateResults(cr.r, cr.id, sent)
	if err != nil {
		s.diag.Error("update result error", err)
	} else {
		// 更新masterpc结果上传标识
		s.PatchResultFlag(cr.r.ID, true, cr.ip, cr.port)
	}
}
