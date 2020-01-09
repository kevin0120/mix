package aiis

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"sync/atomic"
	"time"

	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"strconv"
	"strings"
	"sync"
)

type Service struct {
	configValue    atomic.Value
	diag           Diagnostic
	httpClient     *resty.Client
	storageService IStorageService
	dispatcherBus  Dispatcher
	notifyService  INotifyService
	updateQueue    map[int64]time.Time
	mtx            sync.Mutex
	dispatcherMap  dispatcherbus.DispatcherMap
	transport      ITransport

	opened  bool
	closing chan struct{}
}

func NewService(c Config, d Diagnostic, dp Dispatcher, ss IStorageService, bs IBrokerService, ns INotifyService) *Service {
	s := &Service{
		diag:           d,
		dispatcherBus:  dp,
		updateQueue:    map[int64]time.Time{},
		opened:         false,
		closing:        make(chan struct{}, 1),
		storageService: ss,
		notifyService:  ns,
	}
	s.configValue.Store(c)
	s.initGlbDispatcher()
	s.initTransport(bs, dp)

	return s
}

func (s *Service) initTransport(bs IBrokerService, dispatcherBus Dispatcher) error {
	switch s.Config().TransportType {
	case TRANSPORT_TYPE_GRPC:
		s.transport = NewGRPCClient(s.diag, s.Config())

	case TRANSPORT_TYPE_BROKER:
		s.transport = NewBrokerClient(s.diag, bs, dispatcherBus)

	default:
		s.transport = NewGRPCClient(s.diag, s.Config())
	}

	s.transport.SetResultPatchHandler(s.onResultPatch)
	s.transport.SetServiceStatusHandler(s.onServiceStatus)
	s.transport.SetStatusHandler(s.onTransportStatus)
	return nil
}

func (s *Service) initGlbDispatcher() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DISPATCHER_SERVICE_STATUS: utils.CreateDispatchHandlerStruct(nil),
	}
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
		return errors.Errorf("RemoveFromQueue Queue ID: %d Not Found", id)
	}

	delete(s.updateQueue, id)

	return nil
}

func (s *Service) uploadQueueTimeoutCheck() {
	defer s.mtx.Unlock()
	s.mtx.Lock()

	var wait4Delete []int64
	for k, v := range s.updateQueue {
		if time.Since(v) > time.Duration(s.Config().Timeout) {
			wait4Delete = append(wait4Delete, k)
		}
	}

	for _, id := range wait4Delete {
		delete(s.updateQueue, id)
	}
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
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
	s.ensureHttpClient()
	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)

	s.dispatcherBus.Register(dispatcherbus.DISPATCHER_RESULT, utils.CreateDispatchHandlerStruct(s.onTighteningResult))

	go s.manage()

	s.transport.Start()

	s.opened = true
	return nil
}

func (s *Service) Close() error {
	if !s.opened {
		return nil
	}

	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	s.transport.Stop()
	s.closing <- struct{}{}
	return s.transport.Stop()
}

func (s *Service) PutResult(body *AIISResult) {

	err := s.AddToQueue(body.ID)
	if err != nil {
		return
	}

	err = s.transport.SendResult(body)
	if err != nil {
		s.diag.Error("Publish Tool Result Failed", err)
	}
}

func (s *Service) ResultToAiisResult(result *storage.Results) (AIISResult, error) {
	aiisResult := AIISResult{}
	resultValue := tightening_device.ResultValue{}
	json.Unmarshal([]byte(result.ResultValue), &resultValue)

	psetDefine := tightening_device.PSetDefine{}
	json.Unmarshal([]byte(result.PSetDefine), &psetDefine)

	dbWorkorder, err := s.storageService.GetWorkOrder(result.WorkorderID, true)
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

	aiisResult.NutID = result.ConsuProductID

	aiisResult.WorkcenterCode = dbWorkorder.WorkcenterCode
	aiisResult.ToolSN = result.ToolSN
	aiisResult.ControllerSN = result.ControllerSN

	aiisResult.Job = fmt.Sprintf("%d", dbWorkorder.JobID)
	aiisResult.Stage = result.Stage

	if result.Result == storage.RESULT_OK {
		aiisResult.Final_pass = tightening_device.RESULT_PASS
		if result.Count == 1 {
			aiisResult.One_time_pass = tightening_device.RESULT_PASS
		} else {
			aiisResult.One_time_pass = tightening_device.RESULT_FAIL
		}

		if s.Config().Recheck {
			if (resultValue.Mi >= result.ToleranceMin && resultValue.Mi <= result.ToleranceMax) &&
				(resultValue.Wi >= result.ToleranceMinDegree && resultValue.Wi <= result.ToleranceMaxDegree) {
				aiisResult.QualityState = tightening_device.RESULT_PASS
				aiisResult.ExceptionReason = ""
			} else {
				aiisResult.QualityState = tightening_device.RESULT_EXCEPTION
				aiisResult.ExceptionReason = tightening_device.RESULT_EXCEPTION
			}
		} else {
			aiisResult.QualityState = tightening_device.RESULT_PASS
			aiisResult.ExceptionReason = ""
		}

	} else {
		aiisResult.Final_pass = tightening_device.RESULT_FAIL
		aiisResult.One_time_pass = tightening_device.RESULT_FAIL
	}

	return aiisResult, nil
}

// 收到控制器结果
func (s *Service) onTighteningResult(data interface{}) {
	if data == nil {
		return
	}

	tighteningResult := data.(tightening_device.TighteningResult)
	dbResult, err := s.storageService.GetResultByID(tighteningResult.ID)
	if err != nil {
		s.diag.Error("Get Result Failed", err)
	}

	aiisResult, err := s.ResultToAiisResult(dbResult)
	if err == nil {
		s.PutResult(&aiisResult)
	}
}

func (s *Service) patchResult(rp *ResultPatch) {
	err := s.storageService.UpdateResultByCount(rp.ID, 0, rp.HasUpload)
	if err == nil {
		s.RemoveFromQueue(rp.ID)
		s.diag.Debug(fmt.Sprintf("结果上传成功 ID:%d", rp.ID))
	} else {
		s.diag.Error(fmt.Sprintf("结果上传失败 ID:%d", rp.ID), err)
	}
}

func (s *Service) reuploadResult() error {
	results, err := s.storageService.ListUnUploadResults()
	if err != nil {
		return err
	}

	for _, v := range results {
		aiisResult, err := s.ResultToAiisResult(&v)
		if err == nil {
			s.PutResult(&aiisResult)
		}
	}

	return nil
}

func (s *Service) manage() {
	for {
		select {
		case <-time.After(time.Duration(s.Config().ResultUploadInteval)):
			err := s.reuploadResult()
			if err != nil {
				s.diag.Error("Reupload Result Failed", err)
			}

		case <-time.After(time.Duration(s.Config().Timeout)):
			s.uploadQueueTimeoutCheck()

		case <-s.closing:
			return
		}
	}
}

// 服务状态变化
func (s *Service) onServiceStatus(status ServiceStatus) {
	s.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_SERVICE_STATUS, status)
}

// 传输连接状态变化
func (s *Service) onTransportStatus(status string) {
	s.dispatcherBus.Dispatch(dispatcherbus.DISPATCHER_SERVICE_STATUS, ServiceStatus{
		Name:   SERVICE_AIIS,
		Status: status,
	})
}

// 收到结果上传反馈
func (s *Service) onResultPatch(rp ResultPatch) {
	s.patchResult(&rp)
}
