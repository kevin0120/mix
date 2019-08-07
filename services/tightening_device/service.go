package tightening_device

import (
	"encoding/json"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/wsnotify"
	"sync"
	"sync/atomic"
)

const (
	MODEL_DESOUTTER_CVI3         = "MODEL_DESOUTTER_CVI3"
	MODEL_DESOUTTER_DELTA_WRENCH = "MODEL_DESOUTTER_DELTA_WRENCH"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type TighteningDevice interface {
}

type Service struct {
	diag           Diagnostic
	configValue    atomic.Value
	runningDevices map[string]TighteningDevice
	mtxDevices     sync.Mutex

	WS *wsnotify.Service
	wsnotify.WSNotify
}

func NewService(c Config, d Diagnostic, pAudi controller.Protocol, pOpenprotocol controller.Protocol) (*Service, error) {
	srv := &Service{
		diag:           d,
		mtxDevices:     sync.Mutex{},
		runningDevices: map[string]TighteningDevice{},
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

func (s *Service) OnWSMsg(data []byte) {
	msg := wsnotify.WSMsg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(string(data), err)
		return
	}

	//msgData, _ := json.Marshal(msg.Data)
	switch msg.Type {
	case WS_TIGHTENING_DEVICE:
		// 请求取得所有拧紧设备

		devices, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_TIGHTENING_DEVICE,
			Data: []Device{
				{
					SN:     "0001",
					Type:   TIGHTENING_DEVICE_TYPE_CONTROLLER,
					Parent: "",
					Status: TIGHTENING_DEVICE_ONLINE,
				},
				{
					SN:     "1001",
					Type:   TIGHTENING_DEVICE_TYPE_TOOL,
					Parent: "0001",
					Status: TIGHTENING_DEVICE_ONLINE,
				},
			},
		})

		s.WS.WSSendTightening(string(devices))

	default:
		// 类型错误
		return
	}
}
