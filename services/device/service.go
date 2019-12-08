package device

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"reflect"
	"sync"
	"sync/atomic"
)

type IBaseDevice interface {

	// 设备状态
	Status() string

	// 设备类型
	DeviceType() string

	// 子设备
	Children() map[string]IBaseDevice

	// 设备配置
	Config() interface{}

	// 设备运行数据
	Data() interface{}

	//设备序列号唯一追踪号
	SerialNumber() string
}

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

	requestDispatchers map[string]*utils.Dispatcher
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
	s.initRequestDispatchers()

	return nil
}

func (s *Service) Close() error {
	for _, v := range s.requestDispatchers {
		v.Release()
	}

	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initRequestDispatchers() {
	s.requestDispatchers = map[string]*utils.Dispatcher{
		WS_DEVICE_STATUS: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
	}

	s.requestDispatchers[WS_DEVICE_STATUS].Register(s.OnWS_DEVICE_STATUS)

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

		children := []string{}
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

func (s *Service) OnWS_DEVICE_STATUS(data interface{}) {
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
