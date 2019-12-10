package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherBus"
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

// TODO: 修改服务中的DISPATCH相关方法
type Service struct {
	diag               Diagnostic
	configValue        atomic.Value
	runningControllers map[string]ITighteningController
	mtxDevices         sync.Mutex
	protocols          map[string]ITighteningProtocol

	WS             *wsnotify.Service
	StorageService *storage.Service
	DispatcherBus  Dispatcher
	DeviceService  *device.Service
	wsnotify.WSNotify
	dispatcherMap        dispatcherBus.DispatcherMap
	Api                  *Api
	requestDispatcherMap dispatcherBus.DispatcherMap
}

func (s *Service) loadTighteningController(c Config) {
	for _, deviceConfig := range c.Devices {
		p, err := s.getProtocol(deviceConfig.Protocol)
		if err != nil {
			s.diag.Error("loadTighteningController", err)
			continue
		}
		c, err := p.CreateController(&deviceConfig, s.DispatcherBus)
		if err != nil {
			s.diag.Error("Create Controller Failed", err)
			continue
		}
		// TODO: 如果控制器序列号没有配置，则通过探测加入设备列表。
		s.addController(deviceConfig.SN, c)
	}
}

func NewService(c Config, d Diagnostic, protocols []ITighteningProtocol, dp Dispatcher) (*Service, error) {

	s := &Service{
		diag:               d,
		DispatcherBus: dp,
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

func (s *Service) getProtocol(protocolName string) (ITighteningProtocol, error) {
	if p, ok := s.protocols[protocolName]; !ok {
		return nil, errors.New("Protocol Is Not Support")
	} else {
		return p, nil
	}
}

func (s *Service) initGblDispatcher() {

	s.dispatcherMap = dispatcherBus.DispatcherMap{
		dispatcherBus.DISPATCH_RESULT_PREVIEW:            utils.CreateDispatchHandlerStruct(s.OnResult),
		dispatcherBus.DISPATCH_TOOL_STATUS_PREVIEW:       utils.CreateDispatchHandlerStruct(s.OnToolStatus),
		dispatcherBus.DISPATCH_CONTROLLER_STATUS_PREVIEW: utils.CreateDispatchHandlerStruct(s.OnControllerStatus),
		dispatcherBus.DISPATCH_IO_PREVIEW:                utils.CreateDispatchHandlerStruct(s.OnIOInputs),
		dispatcherBus.DISPATCH_CONTROLLER_ID_PREVIEW:     utils.CreateDispatchHandlerStruct(s.OnControllerBarcode),
	}

}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.WS.AddNotify(s)

	s.DispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)

	// 启动所有拧紧控制器
	s.startupControllers()

	s.initAndStartRequestDispatchers()

	controllers := s.GetControllers()
	for _, c := range controllers {
		for toolSN, _ := range c.Children() {
			s.doDispatch(dispatcherBus.DISPATCH_NEW_TOOL, toolSN)
		}
	}

	return nil
}

func (s *Service) Close() error {

	for name, v := range s.dispatcherMap {
		s.DispatcherBus.Release(name, v.ID)
	}

	for name, v := range s.requestDispatcherMap {
		s.DispatcherBus.Release(name, v.ID)
	}

	// 关闭所有控制器
	s.shutdownControllers()

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initAndStartRequestDispatchers() {
	s.requestDispatcherMap = dispatcherBus.DispatcherMap{
		dispatcherBus.WS_TOOL_MODE_SELECT: utils.CreateDispatchHandlerStruct(s.OnWS_TOOL_MODE_SELECT),
		dispatcherBus.WS_TOOL_ENABLE:      utils.CreateDispatchHandlerStruct(s.OnWS_TOOL_ENABLE),
		dispatcherBus.WS_TOOL_JOB:         utils.CreateDispatchHandlerStruct(s.OnWS_TOOL_JOB),
		dispatcherBus.WS_TOOL_PSET:        utils.CreateDispatchHandlerStruct(s.OnWS_TOOL_PSET),
	}

	s.DispatcherBus.LaunchDispatchersByHandlerMap(s.requestDispatcherMap)
}

func (s *Service) dispatchRequest(req *wsnotify.WSRequest) {
	if err := s.DispatcherBus.Dispatch(req.WSMsg.Type, req); err != nil {
		s.diag.Error(fmt.Sprintf("dispatchRequest Request Type: %s Error", req.WSMsg.Type), err)
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
	s.doDispatch(dispatcherBus.DISPATCH_RESULT, data)
}

// 控制器IO变化
func (s *Service) OnIOInputs(data interface{}) {
	s.doDispatch(dispatcherBus.DISPATCH_IO, data)
}

// 控制器状态变化
func (s *Service) OnControllerStatus(data interface{}) {
	s.doDispatch(dispatcherBus.DISPATCH_CONTROLLER_STATUS, data)
}

// 工具状态变化
func (s *Service) OnToolStatus(data interface{}) {
	s.doDispatch(dispatcherBus.DISPATCH_TOOL_STATUS, data)
}

// 控制器条码
func (s *Service) OnControllerBarcode(data interface{}) {
	s.doDispatch(dispatcherBus.DISPATCH_CONTROLLER_ID, data)
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.DispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error(fmt.Sprintf("doDispatch: %s", name), err)
	}
	return
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
