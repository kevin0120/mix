package aiis

import (
	"github.com/masami10/rush/services/wsnotify"
	"sync/atomic"
	"time"

	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/broker"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"strconv"
	"strings"
	"sync"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
	PutResultDone()
}

type Endpoint struct {
	url     string
	headers map[string]string
	method  string
}

func NewEndpoint(url string, headers map[string]string, method string) *Endpoint {
	return &Endpoint{
		url:     url,
		headers: headers,
		method:  method,
	}
}

type OnOdooStatus func(status string)

type SyncGun func(string) (int64, error)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	endpoints   []*Endpoint
	httpClient  *resty.Client
	rush_port   string
	DB          *storage.Service
	ws          utils.RecConn
	//ws			websocket.Client
	//LocalWSServer *wsnotify.Service
	//Odoo 		*odoo.Service
	//OnOdooStatus OnOdooStatus
	WS          *wsnotify.Service
	updateQueue map[int64]time.Time
	mtx         sync.Mutex

	OdooStatusDispatcher *utils.Dispatcher
	AiisStatusDispatcher *utils.Dispatcher
	SyncGun              SyncGun
	SN                   string
	rpc                  GRPCClient

	TighteningService *tightening_device.Service
	Broker            *broker.Service
	opened            bool
	toolCollector     chan string
	closing           chan struct{}
}

func NewService(c Config, d Diagnostic, rush_port string) *Service {
	e, _ := c.index()
	s := &Service{
		diag:      d,
		endpoints: e,
		rush_port: rush_port,
		rpc: GRPCClient{
			RPCRecvDispatcher:   utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
			RPCStatusDispatcher: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		},
		OdooStatusDispatcher: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		AiisStatusDispatcher: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		updateQueue:          map[int64]time.Time{},
		mtx:                  sync.Mutex{},
		opened:               false,
		toolCollector:        make(chan string, 16),
		closing:              make(chan struct{}, 1),
	}
	s.rpc.RPCRecvDispatcher.Register(s.OnRPCRecv)
	s.rpc.RPCStatusDispatcher.Register(s.OnRPCStatus)
	s.rpc.srv = s
	s.configValue.Store(c)
	return s
}

func (s *Service) AddToQueue(id int64) error {
	defer s.mtx.Unlock()
	s.mtx.Lock()

	_, e := s.updateQueue[id]
	if e {
		return errors.New("exist")
	}

	s.updateQueue[id] = time.Now()

	return nil
}

func (s *Service) RemoveFromQueue(id int64) error {
	defer s.mtx.Unlock()
	s.mtx.Lock()

	_, e := s.updateQueue[id]
	if !e {
		return errors.New("not found")
	}

	delete(s.updateQueue, id)

	return nil
}

func (s *Service) timeoutCheck() {
	defer s.mtx.Unlock()
	s.mtx.Lock()

	wait4Delete := []int64{}
	for k, v := range s.updateQueue {
		if time.Since(v) > time.Duration(s.Config().Timeout) {
			wait4Delete = append(wait4Delete, k)
		}
	}

	for _, id := range wait4Delete {
		delete(s.updateQueue, id)
	}
}

func (s *Service) taskUpdateTimeoutCheck() {
	for {
		s.timeoutCheck()

		time.Sleep(time.Duration(s.Config().Timeout))
	}
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

type TighteningResultHandler = func(data *tightening_device.TighteningResult)

func (s *Service) RegisterTighteningResultHandler(name string, handler TighteningResultHandler) {
	fn := func(data interface{}) {
		d := data.(*tightening_device.TighteningResult)
		handler(d)
	}
	s.TighteningService.GetDispatcher(tightening_device.DISPATCH_RESULT).Register(fn)
}

func (s *Service) ensureHttpClient() *resty.Client {
	if s.httpClient != nil {
		return s.httpClient
	}
	c := s.Config()
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(time.Duration(c.Timeout))
	client.SetContentLength(true)
	// Headers for all request
	client.
		SetRetryCount(c.MaxRetry).
		SetRetryWaitTime(time.Duration(c.PushInterval)).
		SetRetryMaxWaitTime(20 * time.Second)

	s.httpClient = client
	return client
}

func (s *Service) Open() error {
	s.RegisterTighteningResultHandler(tightening_device.DISPATCH_RESULT, s.OnTighteningResult)
	s.TighteningService.GetDispatcher(tightening_device.DISPATCH_NEW_TOOL).Register(s.onNewTool)
	s.Broker.BrokerStatusDispatcher.Register(s.onBrokerStatus)

	s.ensureHttpClient()

	//entry := strings.Split(s.Config().Urls[0], "://")[1]
	//url := url.URL{Scheme: "ws", Host: entry, Path: s.Config().WSResultRoute}
	//s.ws = utils.RecConn{}
	//s.ws.OnConnected = func() {
	//	ws_msg := WSMsg{
	//		Type: WS_REG,
	//		Data: WSRegist{
	//			Rush_SN: s.SN,
	//		},
	//	}
	//
	//	str, _ := json.Marshal(ws_msg)
	//	s.ws.WriteMessage(websocket.TextMessage, str)
	//}
	//
	//s.ws.Dial(url.String(), nil)

	s.OdooStatusDispatcher.Start()
	s.AiisStatusDispatcher.Start()
	go s.taskUpdateTimeoutCheck()
	go s.ResultUploadManager()

	s.rpc.Start()
	s.opened = true
	return nil
}

func (s *Service) Close() error {
	if !s.opened {
		return nil
	}
	s.OdooStatusDispatcher.Release()
	s.AiisStatusDispatcher.Release()
	s.closing <- struct{}{}
	return s.rpc.Stop()
}

func (s *Service) OnRPCStatus(data interface{}) {
	if data == nil {
		return
	}

	status := data.(string)
	// 如果RPC连接断开， 认为ODOO连接也断开
	if status == RPC_OFFLINE {
		s.OdooStatusDispatcher.Dispatch(ODOO_STATUS_OFFLINE)
	}

	s.AiisStatusDispatcher.Dispatch(status)
}

func (s *Service) OnRPCRecv(data interface{}) {
	if data == nil {
		return
	}

	payload := data.(string)
	rpcPayload := RPCPayload{}
	json.Unmarshal([]byte(payload), &rpcPayload)
	str_data, _ := json.Marshal(rpcPayload.Data)

	switch rpcPayload.Type {
	//case TYPE_RESULT:
	//	rp := ResultPatch{}
	//	json.Unmarshal(str_data, &rp)
	//	err := s.DB.UpdateResultByCount(rp.ID, 0, rp.HasUpload)
	//	if err == nil {
	//		s.RemoveFromQueue(rp.ID)
	//		s.diag.Debug(fmt.Sprintf("结果上传成功 ID:%d", rp.ID))
	//	} else {
	//		s.diag.Error(fmt.Sprintf("结果上传失败 ID:%d", rp.ID), err)
	//	}
	//	break

	case TYPE_ODOO_STATUS:
		status := ODOOStatus{}
		json.Unmarshal(str_data, &status)

		s.OdooStatusDispatcher.Dispatch(status.Status)
		break

	case TYPE_MES_STATUS:
		// TODO: 收到mes状态更新, 通知HMI------doing
		//s.WS.WSSend(wsnotify.WS_EVENT_MES,"MES在线")
		body, _ := json.Marshal(wsnotify.GenerateResult(0, "", payload))
		s.WS.WSSend(wsnotify.WS_EVENT_MES, string(body))
		s.diag.Debug(fmt.Sprintf("收到mes状态并推送HMI: %s", payload))

		//case TYPE_ORDER_START:
		//	// TODO: 开工响应------doing
		//	//s.WS.WSSendMes(wsnotify.WS_EVENT_MES,"123","mes允许开工")
		//	break
		//
		//case TYPE_ORDER_FINISH:
		//	// TODO: 完工响应------doing
		//	//s.WS.WSSendMes(wsnotify.WS_EVENT_MES,"123","mes确认完工")
		//	break
	}
}

//func (s *Service) PutResult(msg *WSMsg) error {
//
//}

func (s *Service) PutResult(resultId int64, body *AIISResult) error {

	//var err error
	//for _, endpoint := range s.endpoints {
	//	err = s.putResult(body, fmt.Sprintf(endpoint.url, resultId), endpoint.method)
	//	if err == nil {
	//		// 如果第一次就成功，推出循环
	//		return nil
	//	}
	//}

	//ws_msg := WSMsg{
	//	Type: WS_RESULT,
	//	Data: WSOpResult{
	//		ResultID: resultId,
	//		Result:   body.(AIISResult),
	//		Port:     s.rush_port,
	//	},
	//}
	//
	//str, _ := json.Marshal(ws_msg)
	//s.ws.WriteMessage(websocket.TextMessage, str)
	//s.ws.SendText(string(str))

	result := WSOpResult{
		ResultID: resultId,
		Result:   body,
		Port:     s.rush_port,
	}

	err := s.AddToQueue(result.ResultID)
	if err != nil {
		return nil
	}

	str, _ := json.Marshal(result)
	return s.Broker.Publish(fmt.Sprintf(SUBJECT_RESULTS, body.ToolSN), str)
	//err = s.rpc.RPCSend(string(str))
	//if err != nil {
	//	s.diag.Error("grpc err", err)
	//}
	//
	//return err
}

func (s *Service) PutMesOpenRequest(sn uint64, wsType string, code string, req interface{}, ch <-chan int) (interface{}, error) {
	urlString := s.Config().MesOpenWorkUrl
	url := fmt.Sprintf(urlString, code)
	resp, err := s.httpClient.R().
		SetBody(req).
		Put(url)

	if err != nil {
		body, _ := json.Marshal(wsnotify.GenerateReply(sn, wsType, -2, err.Error()))
		s.WS.WSSend(wsnotify.WS_EVENT_REPLY, string(body))
		<-ch
		return nil, err
	}
	//_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, resp.(string)))
	body, _ := json.Marshal(wsnotify.GenerateResult(sn, wsType, resp.Body()))
	s.WS.WSSend(wsnotify.WS_EVENT_ORDER, string(body))
	<-ch
	return resp.Body(), nil
}

func (s *Service) PutMesFinishRequest(sn uint64, wsType string, code string, req interface{}, ch <-chan int) (interface{}, error) {
	url := fmt.Sprintf(s.Config().MesFinishWorkUrl, code)
	resp, err := s.httpClient.R().
		SetBody(req).
		// or SetError(Error{}).
		Put(url)

	if err != nil {
		body, _ := json.Marshal(wsnotify.GenerateReply(sn, wsType, -2, err.Error()))
		s.WS.WSSend(wsnotify.WS_EVENT_REPLY, string(body))
		<-ch
		return nil, err
	}
	//_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
	body, _ := json.Marshal(wsnotify.GenerateResult(sn, wsType, resp.Body()))
	s.WS.WSSend(wsnotify.WS_EVENT_ORDER, string(body))
	<-ch
	return resp.Body(), nil
}

func (s *Service) PutOrderRequest(reqType string, body interface{}) error {
	msg, _ := json.Marshal(RPCPayload{
		Type: reqType,
		Data: body,
	})

	err := s.rpc.RPCSend(string(msg))
	if err != nil {
		s.diag.Error("Grpc Err", err)
	}

	return err
}

func (s *Service) ResultToAiisResult(result *storage.Results) (AIISResult, error) {
	aiisResult := AIISResult{}
	resultValue := tightening_device.ResultValue{}
	json.Unmarshal([]byte(result.ResultValue), &resultValue)

	psetDefine := tightening_device.PSetDefine{}
	json.Unmarshal([]byte(result.PSetDefine), &psetDefine)

	dbWorkorder, err := s.DB.GetWorkorder(result.WorkorderID, true)
	if err == nil {
		aiisResult.Payload = dbWorkorder.Payload
	}

	aiisResult.CURObjects = append(aiisResult.CURObjects, CURObject{OP: result.Count, File: fmt.Sprintf("%s_%s.json", result.ToolSN, result.TighteningID)})
	aiisResult.ID = result.Id
	aiisResult.WorkorderID = dbWorkorder.WorkorderID
	aiisResult.Control_date = result.UpdateTime.Format(time.RFC3339)
	aiisResult.Measure_degree = resultValue.Wi
	aiisResult.Measure_result = strings.ToLower(result.Result)
	aiisResult.Measure_t_don = resultValue.Ti
	aiisResult.Measure_torque = resultValue.Mi
	aiisResult.Op_time = result.Count
	aiisResult.Pset_m_max = psetDefine.Mp
	aiisResult.Pset_m_min = psetDefine.Mm
	aiisResult.Pset_m_target = psetDefine.Ma
	aiisResult.Pset_m_threshold = psetDefine.Ms
	aiisResult.Pset_strategy = psetDefine.Strategy
	aiisResult.Pset_w_max = psetDefine.Wp
	aiisResult.Pset_w_min = psetDefine.Wm
	aiisResult.Pset_w_target = psetDefine.Wa
	aiisResult.Pset_w_threshold = 1
	aiisResult.UserID = result.UserID
	aiisResult.Seq = result.Seq

	aiisResult.MO_AssemblyLine = dbWorkorder.MO_AssemblyLine
	aiisResult.MO_EquipemntName = dbWorkorder.MO_EquipemntName
	aiisResult.MO_FactoryName = dbWorkorder.MO_FactoryName
	aiisResult.MO_Pin = dbWorkorder.MO_Pin
	aiisResult.MO_Pin_check_code = dbWorkorder.MO_Pin_check_code
	aiisResult.MO_Year = dbWorkorder.MO_Year
	aiisResult.MO_Lnr = dbWorkorder.MO_Lnr
	aiisResult.MO_NutNo = result.NutNo
	aiisResult.MO_Model = dbWorkorder.MO_Model
	aiisResult.Batch = result.Batch
	aiisResult.Vin = dbWorkorder.Track_code
	aiisResult.WorkorderName = dbWorkorder.Code
	aiisResult.Mode = dbWorkorder.Mode
	aiisResult.TighteningId, _ = strconv.ParseInt(result.TighteningID, 10, 64)
	aiisResult.Lacking = "normal"

	//gun, err := s.StorageService.GetGun(result.ToolSN)
	//if err != nil {
	//	gid, err := s.SyncGun(result.ToolSN)
	//	if err == nil {
	//		gun.GunID = gid
	//		gun.Serial = result.ToolSN
	//		s.StorageService.Store(gun)
	//	} else {
	//		gun.GunID = 0
	//	}
	//}

	//aiisResult.GunID = gun.GunID
	//aiisResult.WorkcenterID = dbWorkorder.WorkcenterID
	//aiisResult.ProductID = dbWorkorder.ProductID
	aiisResult.NutID = result.ConsuProductID

	aiisResult.WorkcenterCode = dbWorkorder.WorkcenterCode
	aiisResult.ToolSN = result.ToolSN
	aiisResult.ControllerSN = result.ControllerSN

	aiisResult.Job = fmt.Sprintf("%d", dbWorkorder.JobID)
	aiisResult.Stage = result.Stage

	if result.Result == storage.RESULT_OK {
		aiisResult.Final_pass = ODOO_RESULT_PASS
		if result.Count == 1 {
			aiisResult.One_time_pass = ODOO_RESULT_PASS
		} else {
			aiisResult.One_time_pass = ODOO_RESULT_FAIL
		}

		if s.Config().Recheck {
			if (resultValue.Mi >= result.ToleranceMin && resultValue.Mi <= result.ToleranceMax) &&
				(resultValue.Wi >= result.ToleranceMinDegree && resultValue.Wi <= result.ToleranceMaxDegree) {
				aiisResult.QualityState = QUALITY_STATE_PASS
				aiisResult.ExceptionReason = ""
			} else {
				aiisResult.QualityState = QUALITY_STATE_EX
				aiisResult.ExceptionReason = QUALITY_STATE_EX
			}
		} else {
			aiisResult.QualityState = QUALITY_STATE_PASS
			aiisResult.ExceptionReason = ""
		}

	} else {
		aiisResult.Final_pass = ODOO_RESULT_FAIL
		aiisResult.One_time_pass = ODOO_RESULT_FAIL
	}

	return aiisResult, nil
}

func (s *Service) ResultUploadManager() error {
	for {

		results, err := s.DB.ListUnuploadResults()
		if err == nil {
			for _, v := range results {
				aiisResult, err := s.ResultToAiisResult(&v)
				if err == nil {
					s.PutResult(v.Id, &aiisResult)
				}
			}
		}

		time.Sleep(time.Duration(s.Config().ResultUploadInteval))
	}
}

// 收到控制器结果
func (s *Service) OnTighteningResult(data *tightening_device.TighteningResult) {
	if data == nil {
		return
	}

	tighteningResult := data
	dbResult, err := s.DB.GetResultByID(tighteningResult.ID)
	if err != nil {
		s.diag.Error("Get Result Failed", err)
	}

	aiisResult, err := s.ResultToAiisResult(dbResult)
	if err == nil {
		s.PutResult(dbResult.Id, &aiisResult)
	}
}

func (s *Service) collectTools() {
	for {
		select {
		case toolSN := <-s.toolCollector:
			s.Broker.Subscribe(fmt.Sprintf(SUBJECT_RESULTS_RESP, toolSN), s.onResultResp)

		case <-s.closing:
			return
		}
	}
}

func (s *Service) onBrokerStatus(data interface{}) {
	if data == nil {
		return
	}

	brokerStatus := data.(bool)
	if !brokerStatus {
		return
	}

	go s.collectTools()
}

// 检测到新工具
func (s *Service) onNewTool(data interface{}) {
	if data == nil {
		return
	}

	toolSN := data.(string)
	s.toolCollector <- toolSN
}

func (s *Service) onResultResp(message *broker.BrokerMessage) ([]byte, error) {
	if message == nil {
		return nil, nil
	}

	rpcPayload := RPCPayload{}
	json.Unmarshal(message.Body, &rpcPayload)
	str_data, _ := json.Marshal(rpcPayload.Data)
	rp := ResultPatch{}
	json.Unmarshal(str_data, &rp)
	err := s.DB.UpdateResultByCount(rp.ID, 0, rp.HasUpload)
	if err == nil {
		s.RemoveFromQueue(rp.ID)
		s.diag.Debug(fmt.Sprintf("结果上传成功 ID:%d", rp.ID))
	} else {
		s.diag.Error(fmt.Sprintf("结果上传失败 ID:%d", rp.ID), err)
	}

	return nil, nil
}
