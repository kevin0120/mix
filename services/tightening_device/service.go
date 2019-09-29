package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
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

	wsnotify.WSNotify
	utils.DispatchHandler

	dispatchers map[string]*utils.Dispatcher
	Api         *Api
}

func NewService(c Config, d Diagnostic, protocols []ITighteningProtocol) (*Service, error) {

	srv := &Service{
		diag:               d,
		mtxDevices:         sync.Mutex{},
		runningControllers: map[string]ITighteningController{},
		protocols:          map[string]ITighteningProtocol{},
		dispatchers: map[string]*utils.Dispatcher{
			DISPATCH_RESULT:            utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
			DISPATCH_CONTROLLER_STATUS: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
			DISPATCH_IO:                utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
			DISPATCH_TOOL_STATUS:       utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
			DISPATCH_CONTROLLER_ID:     utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		},
	}

	srv.configValue.Store(c)
	srv.Api = &Api{
		service: srv,
	}

	// 载入支持的协议
	for _, protocol := range protocols {
		srv.protocols[protocol.Name()] = protocol
	}

	// 根据配置加载所有拧紧控制器
	for _, deviceConfig := range c.Devices {

		c, err := srv.protocols[deviceConfig.Protocol].CreateController(&deviceConfig)
		if err != nil {
			srv.diag.Error("Create Controller Failed", err)
			continue
		}

		c.GetDispatch(DISPATCH_RESULT).Register(srv.OnResult)
		c.GetDispatch(DISPATCH_TOOL_STATUS).Register(srv.OnToolStatus)
		c.GetDispatch(DISPATCH_CONTROLLER_STATUS).Register(srv.OnControllerStatus)
		c.GetDispatch(DISPATCH_IO).Register(srv.OnIOInputs)

		// TODO: 如果控制器序列号没有配置，则通过探测加入设备列表。
		srv.addController(deviceConfig.SN, c)
	}

	return srv, nil
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

	return nil
}

func (s *Service) Close() error {

	for _, v := range s.dispatchers {
		v.Release()
	}

	// 关闭所有控制器
	s.shutdownControllers()

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

// TODO: case封装成函数
func (s *Service) OnWSMsg(c websocket.Connection, data []byte) {
	var msg wsnotify.WSMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(fmt.Sprintf("OnWSMsg Fail With Message: %# 20X", data), err)
		return
	}

	if msg.Data == nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, "Msg Data Should Not Be nil"))
		return
	}

	byteData, _ := json.Marshal(msg.Data)

	switch msg.Type {
	case WS_TOOL_MODE_SELECT:
		req := ToolModeSelect{}
		_ = json.Unmarshal(byteData, &req)
		err := s.Api.ToolModeSelect(&req)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			return
		}

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))

	case WS_TOOL_ENABLE:
		req := ToolControl{}
		_ = json.Unmarshal(byteData, &req)
		err := s.Api.ToolControl(&req)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			return
		}

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))

	case WS_TOOL_JOB:
		req := JobSet{}
		_ = json.Unmarshal(byteData, &req)
		err := s.Api.ToolJobSet(&req)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
			return
		}

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))

	case WS_TOOL_PSET:
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

	default:
		// 类型错误
		return
	}
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

func (s *Service) GetDispatcher(name string) *utils.Dispatcher {
	return s.dispatchers[name]
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

	tool, err := controller.GetTool(toolSN)
	if err != nil {
		return nil, err
	}

	return tool, nil
}

func (s *Service) startupControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	for _, c := range s.runningControllers {
		err := c.Start()
		if err != nil {
			s.diag.Error("Startup Controller Failed", err)
		}
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
