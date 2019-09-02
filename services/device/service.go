package device

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/wsnotify"
	"reflect"
	"sync"
	"sync/atomic"
)

type IDevice interface {

	// 设备状态
	Status() string

	// 设备类型
	DeviceType() string

	// 子设备
	Children() map[string]IDevice

	// 设备配置
	Config() interface{}

	// 设备运行数据
	Data() interface{}
}

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag           Diagnostic
	configValue    atomic.Value
	runningDevices map[string]IDevice
	mtxDevices     sync.Mutex

	WS *wsnotify.Service
	wsnotify.WSNotify
}

func NewService(c Config, d Diagnostic) (*Service, error) {

	srv := &Service{
		diag:           d,
		runningDevices: map[string]IDevice{},
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

	return nil
}

func (s *Service) Close() error {

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) AddDevice(sn string, d IDevice) {
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

	devices := []DeviceStatus{}
	for k, v := range s.runningDevices {

		devices = append(devices, DeviceStatus{
			SN:       k,
			Type:     v.DeviceType(),
			Status:   v.Status(),
			Children: reflect.ValueOf(v.Children()).MapKeys(),
			Config:   v.Config(),
			Data:     v.Data(),
		})
	}

	return devices
}

func (s *Service) OnWSMsg(c websocket.Connection, data []byte) {
	msg := wsnotify.WSMsg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(string(data), err)
		return
	}

	//sData, _ := json.Marshal(msg.Data)
	switch msg.Type {
	case WS_DEVICE_STATUS:
		devices := s.fetchAllDevices()
		body, _ := json.Marshal(wsnotify.GenerateResult(msg.SN, msg.Type, devices))

		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_DEVICE, string(body))
	}
}
