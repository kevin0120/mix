package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
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
}

func NewService(c Config, d Diagnostic, protocols []ITighteningProtocol) (*Service, error) {

	srv := &Service{
		diag:               d,
		mtxDevices:         sync.Mutex{},
		runningControllers: map[string]ITighteningController{},
		protocols:          map[string]ITighteningProtocol{},
	}

	srv.configValue.Store(c)

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

	// 启动所有拧紧控制器
	s.startupControllers()

	return nil
}

func (s *Service) Close() error {

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
	//s.diag.Debug(fmt.Sprintf("OnWSMsg Recv New Message: %# 20X", data))
	err := json.Unmarshal(data, &msg)
	if err != nil {
		//s.diag.Error(fmt.Sprintf("OnWSMsg Fail With Message: %# 20X", data), err)
		return
	}
	switch msg.Type {
	case WS_TOOL_MODE_SELECT:
		req := ToolModeSelect{}
		strData, _ := json.Marshal(msg.Data)
		_ = json.Unmarshal([]byte(strData), &req)
		tool, err := s.getTool(req.ControllerSN, req.ToolSN)
		if err != nil {
			return
		}

		err = tool.ModeSelect(req.Mode)
		msg := wsnotify.WSMsg{
			Type: msg.Type,
			SN:   msg.SN,
			Data: err,
		}
		reply, _ := json.Marshal(msg)

		err = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, string(reply))
		if err != nil {
			s.diag.Error(fmt.Sprintf("WS_TOOL_ENABLE Send WS_EVENT_REPLY WebSocket Message: %#v Fail", msg), err)
		}

	case WS_TOOL_ENABLE:
		req := ToolEnable{}
		strData, _ := json.Marshal(msg.Data)
		_ = json.Unmarshal([]byte(strData), &req)
		tool, err := s.getTool(req.ControllerSN, req.ToolSN)
		if err != nil {
			return
		}

		err = tool.ToolControl(req.Enable)
		msg := wsnotify.WSMsg{
			Type: WS_TOOL_ENABLE,
			SN:   msg.SN,
			Data: err,
		}
		reply, _ := json.Marshal(msg)

		err = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, string(reply))
		if err != nil {
			s.diag.Error(fmt.Sprintf("WS_TOOL_ENABLE Send WS_EVENT_REPLY WebSocket Message: %#v Fail", msg), err)
		}

	case WS_TOOL_JOB:
		reply := wsnotify.WSMsg{
			Type: WS_TOOL_JOB,
			SN:   msg.SN,
			Data: Reply{
				Result: 0,
				Msg:    "",
			},
		}

		s.WS.WSSendReply(&reply)

	case WS_TOOL_PSET:
		ds := s.StorageService
		if ds == nil {
			s.diag.Error("WS_TOOL_PSET Fail ", errors.New("Please Inject Storage Service First"))
			return
		}
		var req PSetSet
		bData, _ := json.Marshal(msg.Data)

		_ = json.Unmarshal(bData, &req)

		_ = ds.UpdateTool(&storage.Guns{
			Serial: req.ToolSN,
			Trace:  string(bData),
		})

		tool, err := s.getTool(req.ControllerSN, req.ToolSN)
		if err != nil {
			return
		}

		rt := tool.SetPSet(req.PSet)
		reply, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_TOOL_PSET,
			SN:   msg.SN,
			Data: rt,
		})

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, string(reply))

	default:
		// 类型错误
		return
	}
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
