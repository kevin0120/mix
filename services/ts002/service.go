package ts002

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/context"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/io"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	httpd       IHttpService

	validator *validator.Validate

	mesAPI *MesAPI

	io            IIOService
	tightening    ITightening
	dispatcherBus IDispatcher
	storage       IStorage

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(c Config, d Diagnostic, h IHttpService, io IIOService, tightening ITightening, dispatcher IDispatcher, stroage IStorage) *Service {
	s := &Service{
		diag:          d,
		httpd:         h,
		io:            io,
		tightening:    tightening,
		dispatcherBus: dispatcher,
		storage:       stroage,
		validator:     validator.New(),
	}

	s.configValue.Store(c)

	s.setupWSRequestHandlers()

	return s
}

func (s *Service) Config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) setupWSRequestHandlers() {
	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_TOOL_ENABLE: s.OnWS_TOOL_ENABLE,
	})
}

func (s *Service) ensureValidator() *validator.Validate {
	if s.validator != nil {
		return s.validator
	}

	cc := validator.New()
	s.validator = cc
	return cc
}

func (s *Service) initDispatcherRegisters() {

	// 注册websocket请求
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))

	// 接收读卡器数据
	//s.dispatcherBus.Register(dispatcherbus.DispatcherReaderData, utils.CreateDispatchHandlerStruct(s.onNFCData))
	//
	//// 接收拧紧结果
	//s.dispatcherBus.Register(dispatcherbus.DispatcherResult, utils.CreateDispatchHandlerStruct(s.onTighteningResult))
}

func (s *Service) Open() error {
	c := s.Config()
	if !c.Enable {
		return nil
	}

	if err := s.addTS002HTTPHandlers(); err != nil {
		s.diag.Error("Open Error", err)
		return err
	}

	if mes, err := NewMesAPI(c.MesApiConfig, s.diag); err != nil {
		s.diag.Error("Open NewMesAPI Error", err)
		return err
	} else {
		s.mesAPI = mes
	}

	s.initDispatcherRegisters()

	//go s.doHealthCheck()

	return nil
}

func (s *Service) doHealthCheck() {
	s.mesAPI.healthCheck()
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

func (s *Service) addNewHTTPHandler(method, pattern string, handler context.Handler) {

	if h, err := s.getDefaultHandler(); err != nil {
		s.diag.Error("addNewHandler getDefaultHandler", err)
		return
	} else {
		r := Route{
			RouteType:   httpd.ROUTE_TYPE_HTTP,
			Method:      method,
			Pattern:     pattern,
			HandlerFunc: handler,
		}
		if err := h.AddRoute(r); err != nil {
			s.diag.Error(fmt.Sprintf("addNewHTTPHandler AddRoute: %s Error", pattern), err)
		}
	}

}

func (s *Service) addTS002HTTPHandlers() error {

	s.addNewHTTPHandler("PUT", "/alarm", s.putAlarmReq)
	s.addNewHTTPHandler("PUT", "/pset", s.putPSetReq)
	s.addNewHTTPHandler("PUT", "/io", s.putIOReq)

	s.addNewHTTPHandler("PUT", "/equipments/sync", s.putSyncEquipments)

	return nil
}

func (s *Service) validateRequestPayload(req interface{}) error {
	cc := s.ensureValidator()
	if err := cc.Struct(req); err != nil {
		s.diag.Error("validateRequestPayload Error", err)
		return err
	}
	return nil
}

func (s *Service) ioONDuration(sn string, idx int, duration time.Duration) {
	if err := s.ioDoAction(sn, idx, io.OutputStatusOn); err != nil {
		s.diag.Error("Write ON Error", err)
		return
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := s.io.Write(sn, uint16(idx), io.OutputStatusOff); err != nil {
				s.diag.Error("Write OFF Error", err)
			}
			return
		}
	}
}

func (s *Service) ioDoAction(sn string, idx int, status uint16) error {
	if err := s.io.Write(sn, uint16(idx), status); err != nil {
		e := errors.Wrapf(err, "Write io Serial Number: %s, Idx: %d Error", sn, idx)
		s.diag.Error("Write ON Error", e)
		return err
	}
	return nil
}

// 报警控制
func (s *Service) alarmControl(req *RushAlarmReq) error {
	if req == nil {
		return errors.New("alarmControl: Req Is Nil")
	}
	if err := s.validateRequestPayload(req); err != nil {
		return err
	}
	if _, existed := mapMESStatusIO[req.Status]; !existed {
		return errors.Errorf("alarmControl: Status:%s Is Not Support", req.Status)
	}
	c := s.Config()
	iList := c.IOAlarm

	for _, IOIdx := range iList {
		idx := IOIdx / 8 //IO模块索引
		rr := IOIdx % 8  // 真实IO的位数
		sn := s.io.GetIOSerialNumberByIdx(idx)
		go s.ioONDuration(sn, rr, time.Duration(c.IOAlarmLast))
	}

	return nil
}

// PSet下发控制
func (s *Service) psetControl(req *RushPSetReq) error {
	if req == nil {
		return errors.New("psetControl: Req Is Nil")
	}
	if err := s.validateRequestPayload(req); err != nil {
		return err
	}

	return s.tightening.ToolPSetByIP(&tightening_device.PSetSet{
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

	if err := s.validateRequestPayload(req); err != nil {
		return err
	}

	if err := s.validateRequestPayload(req); err != nil {
		return err
	}
	if _, existed := mapMESStatusIO[req.Status]; !existed {
		return errors.Errorf("ioControl: Status:%s Is Not Support", req.Status)
	}
	iList := req.Outputs

	for _, IOIdx := range iList {
		idx := IOIdx / 8 //IO模块索引
		rr := IOIdx % 8  // 真实IO的位数
		sn := s.io.GetIOSerialNumberByIdx(idx)
		if err := s.ioDoAction(sn, rr, mapMESStatusIO[req.Status]); err != nil {
			return err
		}
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
		err := s.alarmControl(&RushAlarmReq{
			Status: "on",
		})

		if err != nil {
			s.diag.Error("Trigger Alarm Failed", err)
		}
	}

	result := MesResultUploadReq{
		UUID:         tighteningResult.PointID,
		ActualAngle:  tighteningResult.MeasureAngle,
		ActualTorque: tighteningResult.MeasureTorque,
		Flag:         mapMeasureResult[tighteningResult.MeasureResult],
	}

	// 上传拧紧结果
	err := s.mesAPI.sendResultData(&result)
	if err != nil {
		s.diag.Error("Upload Result Failed", err)
	}
}

// 收到读卡器信息
func (s *Service) onNFCData(data interface{}) {
	if data == nil {
		return
	}
	code := data.(string)
	if code == "" || s.mesAPI == nil {
		err := errors.Errorf("NFC Data Is: %s, or MES API Is Empty", code)
		s.diag.Error("onNFCData", err)
		return
	}
	if err := s.mesAPI.sendNFCData(code); err != nil {
		s.diag.Error("sendNFCData Error", err)
		//return
	}

	// 如果成功则开锁
	iList := s.Config().IOLocker
	for _, IOIdx := range iList {
		idx := IOIdx / 8 //IO模块索引
		rr := IOIdx % 8  // 真实IO的位数
		sn := s.io.GetIOSerialNumberByIdx(idx)
		err := s.io.Write(sn, uint16(rr), io.OutputStatusOn)
		if err != nil {
			s.diag.Error(fmt.Sprintf("Locker Control Error SN:%s Output:%d", sn, rr), err)
		}
	}
}

func (s *Service) saveEquipments(equipments []Equipment) error {
	for _, v := range equipments {
		if v.Type != "tightening_gun" && v.Type != "tightening_wrench" {
			continue
		}

		body, _ := json.Marshal(v.Location)
		if err := s.storage.UpdateToolLocation(v.EquipmentSN, string(body)); err != nil {
			s.diag.Error("UpdateToolLocation Failed", err)
		}
	}

	return nil
}
