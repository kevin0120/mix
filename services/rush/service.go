package rush

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/aiis/services/aiis"
	"github.com/masami10/aiis/services/changan"
	"github.com/masami10/aiis/services/fis"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/services/odoo"
	"github.com/masami10/aiis/services/storage"
	"github.com/masami10/aiis/services/wsnotify"
	"gopkg.in/resty.v1"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	Info(msg string)
	Debug(msg string)
}

type CResult struct {
	Result *storage.OperationResult
	ID     int64
	IP     string
	Port   string
	Stream *aiis.RPCAiis_RPCNodeServer
}

type RushResult struct {
	ID        int64 `json:"id"`
	HasUpload bool  `json:"has_upload"`
}

type Service struct {
	HTTPDService *httpd.Service
	workers      int
	Opened       bool
	wg           sync.WaitGroup
	chResult     chan CResult
	closing      chan struct{}
	configValue  atomic.Value
	httpClient   *resty.Client
	route        string

	ws            *websocket.Server
	clientManager wsnotify.WSClientManager

	StorageService interface {
		UpdateResults(result *storage.OperationResult, id int64, sent int) error
		BatchSave(results []*storage.ResultObject) error
	}

	Fis     *fis.Service
	Changan *changan.Service

	results chan *storage.ResultObject
	Odoo    *odoo.Service
	diag    Diagnostic

	rpc aiis.GRPCServer
}

func NewService(c Config, d Diagnostic) *Service {
	if c.Enable {
		s := Service{
			diag:     d,
			workers:  c.Workers,
			Opened:   false,
			chResult: make(chan CResult, c.Workers*4),
			route:    c.Route,

			ws: websocket.New(websocket.Config{
				WriteBufferSize: c.WSWriteBufferSize,
				ReadBufferSize:  c.WSReadBufferSize,
				ReadTimeout:     websocket.DefaultWebsocketPongTimeout, //此作为readtimeout, 默认 如果有ping没有发送也成为read time out
			}),
			clientManager: wsnotify.WSClientManager{},
			results:       make(chan *storage.ResultObject, c.BatchSaveRowsLimit),
			rpc:           aiis.GRPCServer{},
		}

		s.rpc.RPCRecv = s.OnRPCRecv
		s.rpc.RPCNewClient = s.OnRPCNewClinet
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

	go s.TaskResultsBatchSave()

	s.rpc.Start(s.Config().GRPCPort)

	s.Opened = true

	s.Odoo.OnStatus = s.OnOdooStatus

	return nil
}

func (s *Service) OnRPCNewClinet(stream *aiis.RPCAiis_RPCNodeServer) {

	odoo_status := odoo.ODOOStatus{
		Status: s.Odoo.Status(),
	}

	payload := aiis.RPCPayload{
		Type: aiis.TYPE_ODOO_STATUS,
		Data: odoo_status,
	}

	str, _ := json.Marshal(payload)
	s.rpc.RPCSend(stream, string(str))
	s.diag.Info("new rpc client")
}

func (s *Service) OnRPCRecv(stream *aiis.RPCAiis_RPCNodeServer, payload string) {
	s.diag.Debug(fmt.Sprintf("收到结果: %s\n", payload))

	op_result := WSOpResult{}
	err := json.Unmarshal([]byte(payload), &op_result)
	if err != nil {
		s.diag.Error("op result error", err)
		return
	}

	cr := CResult{
		Result: &op_result.Result,
		ID:     op_result.ResultID,
		Stream: stream,
	}

	s.AddResultTask(cr)
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

			cr := CResult{
				Result: &op_result.Result,
				ID:     op_result.ResultID,
				IP:     rush_ip,
				Port:   op_result.Port,
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

func (s *Service) AddResultTask(cr CResult) {
	s.chResult <- cr
}

func (s *Service) putFisResult(ctx iris.Context) {

	var r storage.OperationResult
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
	var r storage.OperationResult
	err = ctx.ReadJSON(&r)
	//ctx.Request().Body.Read()

	if err != nil {
		ctx.Writef(fmt.Sprintf("Result Params from Rush wrong: %s", err))
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	rush_port := ctx.GetHeader("rush_port")
	rush_ip := ctx.GetHeader("rush_ip")

	cr := CResult{
		Result: &r,
		ID:     resultId,
		IP:     rush_ip,
		Port:   rush_port,
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

	s.rpc.Stop()

	return nil
}

func (s *Service) OperationToFisResult(r *storage.OperationResult) fis.FisResult {
	var result fis.FisResult
	result.Init()

	result.EquipemntName = r.EquipemntName
	result.FactoryName = r.FactoryName
	result.Year = r.Year
	result.Pin = r.Pin
	result.PinCheckCode = r.Pin_check_code
	result.AssemblyLine = r.AssemblyLine
	result.ResultID = fmt.Sprintf("%s", r.NutNo)
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

func (s *Service) OperationToChanganResult(r *storage.OperationResult) changan.TighteningResults {
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

func (s *Service) PatchResultFlag(stream *aiis.RPCAiis_RPCNodeServer, result_id int64, has_upload bool, ip string, port string) error {
	//if s.httpClient == nil {
	//	return errors.New("rush http client is nil")
	//}
	//
	//rush_result := RushResult{}
	//rush_result.HasUpload = has_upload
	//r := s.httpClient.R().SetBody(rush_result)
	//
	//s_port := port
	//if s_port == "" {
	//	s_port = "80"
	//}
	//url := fmt.Sprintf("http://%s:%s%s/%d", ip, s_port, s.route, result_id)
	//resp, err := r.Patch(url)
	//if err != nil {
	//	return fmt.Errorf("patch result flag failed: %s\n", err)
	//} else {
	//	if resp.StatusCode() != http.StatusOK {
	//		return fmt.Errorf("patch result flag failed: %s\n", resp.Status())
	//	}
	//}

	rushResult := RushResult{
		ID:        result_id,
		HasUpload: has_upload,
	}

	payload := aiis.RPCPayload{
		Type: aiis.TYPE_RESULT,
		Data: rushResult,
	}

	str, _ := json.Marshal(payload)
	s.rpc.RPCSend(stream, string(str))

	return nil
}

func (s *Service) HandleResult(cr *CResult) {

	// 结果推送fis
	sent := 1
	if s.Fis != nil {
		fisResult := s.OperationToFisResult(cr.Result)

		e := s.Fis.PushResult(&fisResult)
		if e != nil {
			sent = 0
			s.diag.Error("push result to fis error", e)
		}
	}

	if s.Changan != nil {
		changanResult := s.OperationToChanganResult(cr.Result)

		if !s.Changan.AndonDB.InsertResult(&changanResult) {
			s.diag.Error("insert andon result failed", nil)
		}
	}

	json_str, _ := json.Marshal(cr.Result)
	json_obj := map[string]interface{}{}
	json.Unmarshal(json_str, &json_obj)
	result := &storage.ResultObject{
		OR:     json_obj,
		ID:     cr.ID,
		Send:   sent,
		IP:     cr.IP,
		Port:   cr.Port,
		Stream: cr.Stream,
	}

	s.AddResult(result)

	// 结果保存数据库
	//err := s.StorageService.UpdateResults(cr.r, cr.id, sent)
	//
	//if err != nil {
	//	s.diag.Error("update result error", err)
	//} else {
	//	// 更新masterpc结果上传标识
	//	//s.PatchResultFlag(cr.r.ID, true, cr.ip, cr.port)
	//}
}

func (s *Service) AddResult(r *storage.ResultObject) {
	s.results <- r
}

func (s *Service) TaskResultsBatchSave() {
	//idx := 0
	c := s.configValue.Load().(Config)
	//results := make([]*storage.ResultObject, c.BatchSaveRowsLimit)
	results := []*storage.ResultObject{}

	for {
		select {
		case <-time.After(time.Duration(c.BatchSaveTimeLimit)):
			if len(results) > 0 {
				if s.StorageService.BatchSave(results) == nil {
					for _, v := range results {
						resultID := int64(v.OR["id"].(float64))
						if resultID > 0 {
							s.PatchResultFlag(v.Stream, int64(v.OR["id"].(float64)), true, v.IP, v.Port)
						}
					}
				}
				results = []*storage.ResultObject{}
			}

		case data := <-s.results:
			results = append(results, data)

			if len(results) == s.Config().BatchSaveRowsLimit {
				if s.StorageService.BatchSave(results) == nil {
					for _, v := range results {
						resultID := int64(v.OR["id"].(float64))
						if resultID > 0 {
							s.PatchResultFlag(v.Stream, int64(v.OR["id"].(float64)), true, v.IP, v.Port)
						}
					}
				}
				results = []*storage.ResultObject{}
			}
		}
	}
}

func (s *Service) OnOdooStatus(status string) {
	odoo_status := odoo.ODOOStatus{
		Status: status,
	}

	payload := aiis.RPCPayload{
		Type: aiis.TYPE_ODOO_STATUS,
		Data: odoo_status,
	}

	str, _ := json.Marshal(payload)
	s.rpc.RPCSendAll(string(str))
}
