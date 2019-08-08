package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
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
	SetModel(model string)
	SetJob(r *JobSet) Reply
	SetPSet(r *PSetSet) Reply
	Enable(r *ToolEnable) Reply
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

	for _, device := range c.Devices {
		switch device.Protocol {
		//case controller.AUDIPROTOCOL:

		case controller.OPENPROTOCOL:
			pOpenprotocol.AddDevice(device, srv)

		default:

		}
	}

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
	fmt.Println(string(data))
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

	case WS_TOOL_ENABLE:
		reply, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_TOOL_ENABLE,
			SN:   msg.SN,
			Data: Reply{
				Result: 0,
				Msg:    "",
			},
		})

		s.WS.WSSendReply(string(reply))

	case WS_TOOL_JOB:
		reply, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_TOOL_JOB,
			SN:   msg.SN,
			Data: Reply{
				Result: 0,
				Msg:    "",
			},
		})

		s.WS.WSSendReply(string(reply))

	case WS_TOOL_PSET:
		//device, _ := s.GetDevice("0")
		//req := PSetSet{}
		//strData, _ := json.Marshal(msg.Data)
		//json.Unmarshal([]byte(strData), &req)
		//
		//rt := device.SetPSet(&req)
		//reply, _ := json.Marshal(wsnotify.WSMsg{
		//	Type: WS_TOOL_PSET,
		//	SN: msg.SN,
		//	Data: rt,
		//})

		reply, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_TOOL_PSET,
			SN:   msg.SN,
			Data: Reply{
				Result: 0,
				Msg:    "",
			},
		})

		s.WS.WSSendReply(string(reply))

	default:
		// 类型错误
		return
	}
}

func (s *Service) AddDevice(sn string, td TighteningDevice) {
	defer s.mtxDevices.Unlock()
	s.mtxDevices.Lock()

	_, exist := s.runningDevices[sn]
	if exist {
		return
	}

	s.runningDevices[sn] = td
}

func (s *Service) GetDevice(sn string) (TighteningDevice, error) {
	defer s.mtxDevices.Unlock()
	s.mtxDevices.Lock()

	td, exist := s.runningDevices[sn]
	if !exist {
		return td, errors.New("device not found")
	}

	return td, nil
}
