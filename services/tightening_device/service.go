package tightening_device

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/controller"
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

	WS             *wsnotify.Service
	StorageService *storage.Service

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

func (s *Service) OnWSMsg(c websocket.Connection, data []byte) {
	var msg wsnotify.WSMsg
	s.diag.Debug(fmt.Sprintf("OnWSMsg Recv New Message: %# 20X", data))
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(fmt.Sprintf("OnWSMsg Fail With Message: %# 20X", data), err)
		return
	}
	switch msg.Type {
	case WS_TOOL_ENABLE:
		req := ToolEnable{}
		strData, _ := json.Marshal(msg.Data)
		_ = json.Unmarshal([]byte(strData), &req)
		device, _ := s.GetDevice(req.ToolSN)
		rt := device.Enable(&req)
		msg := wsnotify.WSMsg{
			Type: WS_TOOL_ENABLE,
			SN:   msg.SN,
			Data: rt,
		}
		reply, _ := json.Marshal(msg)

		err := wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, string(reply))
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

		_ = ds.UpdateGun(&storage.Guns{
			Serial: req.ToolSN,
			Trace:  string(bData),
		})

		device, _ := s.GetDevice(req.ToolSN)
		rt := device.SetPSet(&req)
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
		return nil, errors.New(fmt.Sprintf("Device Serial Number: %s Not Fount", sn))
	}

	return td, nil
}
