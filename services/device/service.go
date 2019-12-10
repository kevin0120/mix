package device

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/DispatcherBus"
	"github.com/masami10/rush/services/wsnotify"
	"reflect"
	"sync"
	"sync/atomic"
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

	WS *wsnotify.Service
	wsnotify.WSNotify

	DispatcherBus *dispatcherBus.Service
	dispatcherMap dispatcherBus.DispatcherMap
}

func NewService(c Config, d Diagnostic) (*Service, error) {

	srv := &Service{
		diag:           d,
		runningDevices: map[string]IBaseDevice{},
		mtxDevices:     sync.Mutex{},
	}

	srv.configValue.Store(c)

	return srv, nil
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.WS.AddNotify(s)

	s.dispatcherMap = dispatcherBus.DispatcherMap{
		dispatcherBus.DISPATCHER_WS_DEVICE_STATUS: s.OnWSDeviceStatus,
	}
	s.DispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)

	return nil
}

func (s *Service) Close() error {
	for name, _ := range s.dispatcherMap {
		s.DispatcherBus.Release(name)
	}

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) dispatchRequest(req *wsnotify.WSRequest) {
	switch req.WSMsg.Type {
	case wsnotify.NotifywsDeviceStatus:
		s.DispatcherBus.Dispatch(dispatcherBus.DISPATCHER_WS_DEVICE_STATUS, req)
	}
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

func (s *Service) OnWSDeviceStatus(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	devices := s.fetchAllDevices()
	body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, devices))

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_DEVICE, string(body))
}

func (s *Service) OnWSMsg(c websocket.Connection, data []byte) {
	msg := wsnotify.WSMsg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(string(data), err)
		return
	}

	s.dispatchRequest(&wsnotify.WSRequest{
		C:     c,
		WSMsg: &msg,
	})
}
