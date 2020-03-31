package tightening_device

import (
	"fmt"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"sync"
)

const (
	ModelDesoutterCvi3        = "ModelDesoutterCvi3"
	ModelDesoutterCvi2        = "ModelDesoutterCvi2"
	ModelDesoutterDeltaWrench = "ModelDesoutterDeltaWrench"
	ModelDesoutterConnect     = "ModelDesoutterConnect"
)

type Service struct {
	diag               Diagnostic
	configValue        atomic.Value
	runningControllers map[string]ITighteningController
	mtxDevices         sync.Mutex
	protocols          map[string]ITighteningProtocol

	storageService IStorageService
	dispatcherBus  Dispatcher
	deviceService  IDeviceService
	dispatcherMap  dispatcherbus.DispatcherMap

	ioService IOService

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

		c, err := p.NewController(&c.Devices[k], s.dispatcherBus)
		if err != nil {
			s.diag.Error("Create Controller Failed", err)
			continue
		}

		c.SetSerialNumber(deviceConfig.SN)
		s.ioService.AddModule(fmt.Sprintf(TIGHTENING_CONTROLLER_IO_SN_FORMAT, c.SerialNumber()), c.CreateIO())
		s.addController(deviceConfig.SN, c)
	}
}

func NewService(c Config, d Diagnostic, protocols []ITighteningProtocol, dp Dispatcher, ds IDeviceService, db IStorageService, io IOService) (*Service, error) {

	s := &Service{
		diag:               d,
		dispatcherBus:      dp,
		runningControllers: map[string]ITighteningController{},
		protocols:          map[string]ITighteningProtocol{},
		deviceService:      ds,
		storageService:     db,
		ioService:          io,
	}

	s.setupGlobalDispatchers()
	s.setupWSRequestHandlers()

	s.configValue.Store(c)

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

func (s *Service) setupGlobalDispatchers() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DispatcherCurve: utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherJob:   utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) setupWSRequestHandlers() {
	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_TOOL_MODE_SELECT:       s.OnWS_TOOL_MODE_SELECT,
		wsnotify.WS_TOOL_ENABLE:            s.OnWS_TOOL_ENABLE,
		wsnotify.WS_TOOL_JOB:               s.OnWS_TOOL_JOB,
		wsnotify.WS_TOOL_PSET:              s.OnWS_TOOL_PSET,
		wsnotify.WS_TOOL_PSET_BATCH:        s.OnWS_TOOL_PSET_BATCH,
		wsnotify.WS_TOOL_PSET_LIST:         s.OnWS_TOOL_PSET_LIST,
		wsnotify.WS_TOOL_PSET_DETAIL:       s.OnWS_TOOL_PSET_DETAIL,
		wsnotify.WS_TOOL_JOB_LIST:          s.OnWS_TOOL_JOB_LIST,
		wsnotify.WS_TOOL_JOB_DETAIL:        s.OnWS_TOOL_JOB_DETAIL,
		wsnotify.WS_TOOL_RESULT_MANUAL_SET: s.OnWS_TOOL_RESULT_MANUAL_SET,
	})
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)
	s.initDispatcherRegisters()

	// 启动所有拧紧控制器
	s.startupControllers()

	return nil
}

func (s *Service) Close() error {

	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)

	// 关闭所有控制器
	s.shutdownControllers()

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initDispatcherRegisters() {
	// 注册websocket请求
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error(fmt.Sprintf("doDispatch: %s", name), err)
	}
	return
}

func (s *Service) getControllers() map[string]ITighteningController {
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
		return nil, errors.New(fmt.Sprintf("Controller %s Not Found", controllerSN))
	}

	return td, nil
}

func (s *Service) getTool(controllerSN string, toolSN string) (ITighteningTool, error) {

	for _, c := range s.runningControllers {
		for _, t := range c.Children() {
			if t.SerialNumber() == toolSN {
				return t.(ITighteningTool), nil
			}
		}
	}

	return nil, errors.New("Tool Not Found")
}

func (s *Service) startupControllers() {
	s.mtxDevices.Lock()
	defer s.mtxDevices.Unlock()

	for sn, c := range s.runningControllers {
		err := c.Start()
		if err != nil {
			s.diag.Error("Startup Controller Failed", err)
		}

		s.deviceService.AddDevice(sn, c)
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
