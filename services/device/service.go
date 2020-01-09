package device

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
	"go.uber.org/atomic"
	"sync"

	"github.com/masami10/rush/services/wsnotify"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag           Diagnostic
	configValue    atomic.Value
	runningDevices map[string]IBaseDevice
	mtxDevices     sync.Mutex
	dispatcherBus  Dispatcher
	dispatcherMap  dispatcherbus.DispatcherMap

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(c Config, d Diagnostic, dp Dispatcher) (*Service, error) {

	s := &Service{
		diag:           d,
		runningDevices: map[string]IBaseDevice{},
		mtxDevices:     sync.Mutex{},
		dispatcherBus:  dp,
	}

	s.configValue.Store(c)

	s.setupGlobalDispatchers()
	s.setupWSRequestHandlers()

	return s, nil
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.initDispatcherRegisters()
	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)

	return nil
}

func (s *Service) Close() error {
	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initDispatcherRegisters() {
	// 注册websocket请求
	s.dispatcherBus.Register(dispatcherbus.DISPATCHER_WS_NOTIFY, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))
}

func (s *Service) setupGlobalDispatchers() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DISPATCHER_DEVICE_STATUS: utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCHER_READER_DATA:   utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCHER_SCANNER_DATA:  utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCHER_IO:            utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) setupWSRequestHandlers() {
	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_DEVICE_STATUS: s.OnWSDeviceStatus,
	})
}

func (s *Service) AddDevice(sn string, d IBaseDevice) {
	defer s.mtxDevices.Unlock()
	s.mtxDevices.Lock()

	_, exist := s.runningDevices[sn]
	if exist {
		return
	}

	s.runningDevices[sn] = d
}

func (s *Service) fetchAllDevices() []DeviceStatus {
	defer s.mtxDevices.Unlock()
	s.mtxDevices.Lock()

	var devices []DeviceStatus
	for k, v := range s.runningDevices {
		devices = append(devices, DeviceStatus{
			SN:       k,
			Type:     v.DeviceType(),
			Status:   v.Status(),
			Children: v.Children(),
			Config:   v.Config(),
			Data:     v.Data(),
		})

		for cSN, c := range v.Children() {
			devices = append(devices, DeviceStatus{
				SN:     cSN,
				Type:   c.DeviceType(),
				Status: c.Status(),
				Config: c.Config(),
				Data:   c.Data(),
			})
		}
	}

	return devices
}
