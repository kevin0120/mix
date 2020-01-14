package aiis

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"go.uber.org/atomic"
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

	closing chan struct{}
}

func NewService(c Config, d Diagnostic, dp Dispatcher, ss IStorageService, bs IBrokerService, ns INotifyService) *Service {
	s := &Service{
		diag:           d,
		dispatcherBus:  dp,
		updateQueue:    map[int64]time.Time{},
		closing:        make(chan struct{}, 1),
		storageService: ss,
		notifyService:  ns,
	}
	s.configValue.Store(c)
	s.setupGlbDispatcher()
	s.setupTransport(bs, dp)

	return s
}

func (s *Service) setupTransport(bs IBrokerService, dispatcherBus Dispatcher) {
	switch s.Config().TransportType {
	case TransportTypeGrpc:
		s.transport = NewGRPCClient(s.diag, s.Config())

	case TransportTypeBroker:
		s.transport = NewBrokerClient(s.diag, bs, dispatcherBus)

	default:
		s.transport = NewGRPCClient(s.diag, s.Config())
	}

	s.transport.SetResultPatchHandler(s.onResultPatch)
	s.transport.SetServiceStatusHandler(s.onServiceStatus)
	s.transport.SetStatusHandler(s.onTransportStatus)
}

func (s *Service) setupGlbDispatcher() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DispatcherServiceStatus: utils.CreateDispatchHandlerStruct(nil),
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

func (s *Service) RemoveFromQueue(id int64) {
	defer s.mtx.Unlock()
	s.mtx.Lock()

	_, e := s.updateQueue[id]
	if e {
		delete(s.updateQueue, id)
	}
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

	s.dispatcherBus.Register(dispatcherbus.DispatcherResult, utils.CreateDispatchHandlerStruct(s.onTighteningResult))

	go s.manage()

	if err := s.transport.Start(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Close() error {
	if err := s.transport.Stop(); err != nil {
		s.diag.Error("Stop Transport Failed", err)
	}

	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	s.closing <- struct{}{}

	return nil
}

func (s *Service) PutResult(body *PublishResult) {

	err := s.AddToQueue(body.ID)
	if err != nil {
		return
	}

	err = s.transport.SendResult(body)
	if err != nil {
		s.diag.Error("Publish Tool Result Failed", err)
	}
}

func (s *Service) ResultToAiisResult(result *storage.Results) (PublishResult, error) {
	aiisResult := PublishResult{}
	resultValue := tightening_device.ResultValue{}
	if err := json.Unmarshal([]byte(result.ResultValue), &resultValue); err != nil {
		s.diag.Error("Unmarshal ResultValue Failed", err)
	}

	psetDefine := tightening_device.PSetDefine{}
	if err := json.Unmarshal([]byte(result.PSetDefine), &psetDefine); err != nil {
		s.diag.Error("Unmarshal PSetDefine Failed", err)
	}

	dbWorkorder, err := s.storageService.GetWorkOrder(result.WorkorderID, true)
	if err != nil {
		return aiisResult, err
	}

	aiisResult.Payload = dbWorkorder.Payload
	aiisResult.CURObjects = append(aiisResult.CURObjects, CURObject{OP: result.Count, File: fmt.Sprintf("%s_%s.json", result.ToolSN, result.TighteningID)})
	aiisResult.ID = result.Id
	aiisResult.WorkorderID = dbWorkorder.WorkorderID
	aiisResult.ControlDate = result.UpdateTime.Format(time.RFC3339)
	aiisResult.MeasureDegree = resultValue.Wi
	aiisResult.MeasureResult = strings.ToLower(result.Result)
	aiisResult.MeasureTDon = resultValue.Ti
	aiisResult.MeasureTorque = resultValue.Mi
	aiisResult.OpTime = result.Count
	aiisResult.PsetMMax = psetDefine.Mp
	aiisResult.PsetMMin = psetDefine.Mm
	aiisResult.PsetMTarget = psetDefine.Ma
	aiisResult.PsetMThreshold = psetDefine.Ms
	aiisResult.PsetStrategy = psetDefine.Strategy
	aiisResult.PsetWMax = psetDefine.Wp
	aiisResult.PsetWMin = psetDefine.Wm
	aiisResult.PsetWTarget = psetDefine.Wa
	aiisResult.PsetWThreshold = 1
	aiisResult.UserID = result.UserID
	aiisResult.Seq = result.Seq

	aiisResult.MoAssemblyline = dbWorkorder.MO_AssemblyLine
	aiisResult.MoEquipemntname = dbWorkorder.MO_EquipemntName
	aiisResult.MoFactoryname = dbWorkorder.MO_FactoryName
	aiisResult.MoPin = dbWorkorder.MO_Pin
	aiisResult.MoPinCheckCode = dbWorkorder.MO_Pin_check_code
	aiisResult.MoYear = dbWorkorder.MO_Year
	aiisResult.MoLnr = dbWorkorder.MO_Lnr
	aiisResult.MoNutno = result.NutNo
	aiisResult.MoModel = dbWorkorder.MO_Model
	aiisResult.Batch = result.Batch
	aiisResult.Vin = dbWorkorder.TrackCode
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
		aiisResult.FinalPass = tightening_device.RESULT_PASS
		if result.Count == 1 {
			aiisResult.OneTimePass = tightening_device.RESULT_PASS
		} else {
			aiisResult.OneTimePass = tightening_device.RESULT_FAIL
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
		aiisResult.FinalPass = tightening_device.RESULT_FAIL
		aiisResult.OneTimePass = tightening_device.RESULT_FAIL
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
	s.doDispatch(dispatcherbus.DispatcherServiceStatus, status)
}

// 传输连接状态变化
func (s *Service) onTransportStatus(status string) {
	s.doDispatch(dispatcherbus.DispatcherServiceStatus, ServiceStatus{
		Name:   ServiceAiis,
		Status: status,
	})
}

// 收到结果上传反馈
func (s *Service) onResultPatch(rp ResultPatch) {
	s.patchResult(&rp)
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error("Dispatch Failed", err)
	}
}
