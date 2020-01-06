package device

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
	"reflect"
	"sync"
	"sync/atomic"

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
	WS             *wsnotify.Service
	DispatcherBus  Dispatcher
	dispatcherMap  dispatcherbus.DispatcherMap

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(c Config, d Diagnostic) (*Service, error) {

	s := &Service{
		diag:           d,
		runningDevices: map[string]IBaseDevice{},
		mtxDevices:     sync.Mutex{},
	}

	s.configValue.Store(c)

	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.initGlobalDispatchers()
	s.initWSRequestHandlers()

	return s, nil
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.DispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)
	// 注册websocket请求
	s.DispatcherBus.Register(dispatcherbus.DISPATCHER_WS_NOTIFY, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))

	return nil
}

func (s *Service) Close() error {
	s.DispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initGlobalDispatchers() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DISPATCHER_DEVICE_STATUS: utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCHER_READER_DATA:   utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCHER_SCANNER_DATA:  utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DISPATCHER_IO:            utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) initWSRequestHandlers() {
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

		var children []string
		if v.Children() != nil {
			for _, child := range reflect.ValueOf(v.Children()).MapKeys() {
				children = append(children, child.String())
			}
		}

		devices = append(devices, DeviceStatus{
			SN:       k,
			Type:     v.DeviceType(),
			Status:   v.Status(),
			Children: children,
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

func (s *Service) OnWSDeviceStatus(c websocket.Connection, msg *wsnotify.WSMsg) {
	devices := s.fetchAllDevices()
	body, _ := json.Marshal(wsnotify.GenerateWSMsg(msg.SN, msg.Type, devices))

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_DEVICE, string(body))
}

func (s *Service) SendDeviceStatus(deviceStatus DeviceStatus) {
	body, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.WS_DEVICE_STATUS,
		Data: []DeviceStatus{
			deviceStatus,
		},
	})

	s.WS.NotifyAll(wsnotify.WS_EVENT_DEVICE, string(body))
}
