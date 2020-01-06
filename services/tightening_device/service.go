package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

const (
	ModelDesoutterCvi3        = "ModelDesoutterCvi3"
	ModelDesoutterCvi2        = "ModelDesoutterCvi2"
	ModelDesoutterDeltaWrench = "ModelDesoutterDeltaWrench"
)

// TODO: 修改服务中的DISPATCH相关方法
type Service struct {
	diag               Diagnostic
	configValue        atomic.Value
	runningControllers map[string]ITighteningController
	mtxDevices         sync.Mutex
	protocols          map[string]ITighteningProtocol

	StorageService *storage.Service
	DispatcherBus  Dispatcher
	DeviceService  *device.Service
	dispatcherMap  dispatcherbus.DispatcherMap
	Api            *Api

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func (s *Service) loadTighteningController(c Config) {
	for k, deviceConfig := range c.Devices {
		p, err := s.getProtocol(deviceConfig.Protocol)
		if err != nil {
			s.diag.Error("loadTighteningController", err)
			continue
		}
		c, err := p.CreateController(&c.Devices[k], s.DispatcherBus)
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
		DispatcherBus:      dp,
		mtxDevices:         sync.Mutex{},
		runningControllers: map[string]ITighteningController{},
		protocols:          map[string]ITighteningProtocol{},
	}

	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.initGlobalDispatchers()
	s.initWSRequestHandlers()

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

func (s *Service) initGlobalDispatchers() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DISPATCH_RESULT:   utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCH_CURVE:    utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCH_NEW_TOOL: utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) initWSRequestHandlers() {
	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_TOOL_MODE_SELECT: s.OnWS_TOOL_MODE_SELECT,
		wsnotify.WS_TOOL_ENABLE:      s.OnWS_TOOL_ENABLE,
		wsnotify.WS_TOOL_JOB:         s.OnWS_TOOL_JOB,
		wsnotify.WS_TOOL_PSET:        s.OnWS_TOOL_PSET,
		wsnotify.WS_TOOL_PSET_BATCH:  s.OnWS_TOOL_PSET_BATCH,
	})
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.DispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)

	// 启动所有拧紧控制器
	s.startupControllers()

	controllers := s.GetControllers()
	for _, c := range controllers {
		for toolSN, _ := range c.Children() {
			s.doDispatch(dispatcherbus.DISPATCH_NEW_TOOL, toolSN)
		}
	}

	// 注册websocket请求
	s.DispatcherBus.Register(dispatcherbus.DISPATCH_WS_NOTIFY, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))

	return nil
}

func (s *Service) Close() error {

	for name, v := range s.dispatcherMap {
		s.DispatcherBus.Release(name, v.ID)
	}

	// 关闭所有控制器
	s.shutdownControllers()

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) OnWS_TOOL_MODE_SELECT(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolModeSelect{}
	_ = json.Unmarshal(byteData, &req)
	err := s.Api.ToolModeSelect(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_ENABLE(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolControl{}
	_ = json.Unmarshal(byteData, &req)
	err := s.Api.ToolControl(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_JOB(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := JobSet{}
	_ = json.Unmarshal(byteData, &req)
	err := s.Api.ToolJobSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))

}

func (s *Service) OnWS_TOOL_PSET(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := PSetSet{}
	_ = json.Unmarshal(byteData, &req)

	err := s.Api.ToolPSetSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) OnWS_TOOL_PSET_BATCH(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := PSetBatchSet{}
	_ = json.Unmarshal(byteData, &req)

	err := s.Api.ToolPSetBatchSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
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
