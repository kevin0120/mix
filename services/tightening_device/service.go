package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

const (
	ModelDesoutterCvi3        = "ModelDesoutterCvi3"
	ModelDesoutterDeltaWrench = "ModelDesoutterDeltaWrench"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag               Diagnostic
	configValue        atomic.Value
	runningControllers map[string]ITighteningController
	mtxDevices         sync.Mutex
	protocols          map[string]ITighteningProtocol

	WS             *wsnotify.Service
	StorageService *storage.Service

	DeviceService *device.Service

	wsnotify.WSNotify
	utils.DispatchHandler

	dispatchers map[string]*utils.Dispatcher
	Api         *Api

	requestDispatchers map[string]*utils.Dispatcher
}

func (s *Service) initGblDispatcher() {
	s.dispatchers = map[string]*utils.Dispatcher{
		DISPATCH_RESULT:            utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		DISPATCH_CONTROLLER_STATUS: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		DISPATCH_IO:                utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		DISPATCH_TOOL_STATUS:       utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		DISPATCH_CONTROLLER_ID:     utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		DISPATCH_NEW_TOOL:          utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
	}
}

func (s *Service) loadTighteningController(c Config) {
	for _, deviceConfig := range c.Devices {
		p , err := s.getProtocol(deviceConfig.Protocol)
		if err != nil {
			s.diag.Error("loadTighteningController", err)
			continue
		}
		c, err := p.CreateController(&deviceConfig)
		if err != nil {
			s.diag.Error("Create Controller Failed", err)
			continue
		}

		c.GetDispatch(DISPATCH_RESULT).Register(s.OnResult)
		c.GetDispatch(DISPATCH_TOOL_STATUS).Register(s.OnToolStatus)
		c.GetDispatch(DISPATCH_CONTROLLER_STATUS).Register(s.OnControllerStatus)
		c.GetDispatch(DISPATCH_IO).Register(s.OnIOInputs)
		c.GetDispatch(DISPATCH_CONTROLLER_ID).Register(s.OnControllerBarcode)

		// TODO: 如果控制器序列号没有配置，则通过探测加入设备列表。
		s.addController(deviceConfig.SN, c)
	}
}

func NewService(c Config, d Diagnostic, protocols []ITighteningProtocol) (*Service, error) {

	s := &Service{
		diag:               d,
		mtxDevices:         sync.Mutex{},
		runningControllers: map[string]ITighteningController{},
		protocols:          map[string]ITighteningProtocol{},
	}

	s.initGblDispatcher()

	s.configValue.Store(c)
	s.Api = &Api{
		s,
	}

	// 载入支持的协议
	for _, protocol := range protocols {
		s.protocols[protocol.Name()] = protocol
	}

	// 根据配置加载所有拧紧控制器
	s.loadTighteningController(c)
	return s, nil
}

func (s *Service)getProtocol(protocolName string) (ITighteningProtocol, error) {
	if p, ok := s.protocols[protocolName]; !ok {
		return nil, errors.New("Protocol Is Not Support")
	}else {
		return p, nil
	}
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.WS.AddNotify(s)

	for _, v := range s.dispatchers {
		v.Start()
	}

	// 启动所有拧紧控制器
	s.startupControllers()

	s.initRequestDispatchers()

	controllers := s.GetControllers()
	for _, c := range controllers {
		for toolSN, _ := range c.Children() {
			s.GetDispatcher(DISPATCH_NEW_TOOL).Dispatch(toolSN)
		}
	}

	return nil
}

func (s *Service) Close() error {

	for _, v := range s.dispatchers {
		v.Release()
	}

	for _, v := range s.requestDispatchers {
		v.Release()
	}

	// 关闭所有控制器
	s.shutdownControllers()

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initRequestDispatchers() {
	s.requestDispatchers = map[string]*utils.Dispatcher{
		WS_TOOL_MODE_SELECT: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_TOOL_ENABLE:      utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_TOOL_JOB:         utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_TOOL_PSET:        utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
	}

	s.requestDispatchers[WS_TOOL_MODE_SELECT].Register(s.OnWS_TOOL_MODE_SELECT)
	s.requestDispatchers[WS_TOOL_ENABLE].Register(s.OnWS_TOOL_ENABLE)
	s.requestDispatchers[WS_TOOL_JOB].Register(s.OnWS_TOOL_JOB)
	s.requestDispatchers[WS_TOOL_PSET].Register(s.OnWS_TOOL_PSET)

	for _, v := range s.requestDispatchers {
		v.Start()
	}
}

func (s *Service) dispatchRequest(req *wsnotify.WSRequest) {
	d, exist := s.requestDispatchers[req.WSMsg.Type]
	if exist {
		d.Dispatch(req)
	}
}

func (s *Service) OnWS_TOOL_MODE_SELECT(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	byteData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	req := ToolModeSelect{}
	_ = json.Unmarshal(byteData, &req)
	err := s.Api.ToolModeSelect(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}
func (s *Service) OnWS_TOOL_ENABLE(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	byteData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	req := ToolControl{}
	_ = json.Unmarshal(byteData, &req)
	err := s.Api.ToolControl(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_JOB(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	byteData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	req := JobSet{}
	_ = json.Unmarshal(byteData, &req)
	err := s.Api.ToolJobSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))

}

func (s *Service) OnWS_TOOL_PSET(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	byteData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	ds := s.StorageService
	if ds == nil {
		s.diag.Error("WS_TOOL_PSET Fail ", errors.New("Please Inject Storage Service First"))
		return
	}
	req := PSetSet{}
	_ = json.Unmarshal(byteData, &req)

	err := s.Api.ToolPSetSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWSMsg(c websocket.Connection, data []byte) {
	var msg wsnotify.WSMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(fmt.Sprintf("OnWSMsg Fail With Message: %# 20X", data), err)
		return
	}

	s.dispatchRequest(&wsnotify.WSRequest{
		C:     c,
		WSMsg: &msg,
	})
}

// 收到结果
func (s *Service) OnResult(data interface{}) {
	s.GetDispatcher(DISPATCH_RESULT).Dispatch(data)
}

// 控制器IO变化
func (s *Service) OnIOInputs(data interface{}) {
	s.GetDispatcher(DISPATCH_IO).Dispatch(data)
}

// 控制器状态变化
func (s *Service) OnControllerStatus(data interface{}) {
	s.GetDispatcher(DISPATCH_CONTROLLER_STATUS).Dispatch(data)
}

// 工具状态变化
func (s *Service) OnToolStatus(data interface{}) {
	s.GetDispatcher(DISPATCH_TOOL_STATUS).Dispatch(data)
}

// 控制器条码
func (s *Service) OnControllerBarcode(data interface{}) {
	s.GetDispatcher(DISPATCH_CONTROLLER_ID).Dispatch(data)
}

func (s *Service) GetDispatcher(name string) *utils.Dispatcher {
	if d, ok := s.dispatchers[name]; ok {
		return d
	} else {
		err := errors.Errorf("Dispatcher : %sIs Not Existed!", name)
		s.diag.Error("GetDispatcher", err)
		return nil
	}
}

func (s *Service) GetControllers() map[string]ITighteningController {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	return s.runningControllers
}

func (s *Service) addController(controllerSN string, controller ITighteningController) {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	_, exist := s.runningControllers[controllerSN]
	if exist {
		return
	}

	s.runningControllers[controllerSN] = controller
}

func (s *Service) getController(controllerSN string) (ITighteningController, error) {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	td, exist := s.runningControllers[controllerSN]
	if !exist {
		return nil, errors.New(fmt.Sprintf("Controller %s Not Fount", controllerSN))
	}

	return td, nil
}

func (s *Service) getTool(controllerSN string, toolSN string) (ITighteningTool, error) {
	controller, err := s.getController(controllerSN)
	if err != nil {
		return nil, err
	}

	tool, err := controller.GetToolViaSerialNumber(toolSN)
	if err != nil {
		return nil, err
	}

	return tool, nil
}

func (s *Service) startupControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	for sn, c := range s.runningControllers {
		err := c.Start()
		if err != nil {
			s.diag.Error("Startup Controller Failed", err)
		}

		s.DeviceService.AddDevice(sn, c)
	}
}

func (s *Service) shutdownControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	for _, c := range s.runningControllers {
		err := c.Stop()
		if err != nil {
			s.diag.Error("Shutdown Controller Failed", err)
		}
	}
}
