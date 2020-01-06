package io

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type Service struct {
	configValue   atomic.Value
	ios           map[string]*IOModule
	diag          Diagnostic
	WS            *wsnotify.Service
	DispatcherBus Dispatcher
	DeviceService IDeviceService
	IONotify

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(c Config, d Diagnostic, dp Dispatcher, ds IDeviceService) *Service {

	s := &Service{
		diag:          d,
		DispatcherBus: dp,
		DeviceService: ds,
		ios:           map[string]*IOModule{},
	}

	s.configValue.Store(c)

	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.initWSRequestHandlers()

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	cfgs := s.config().IOS
	for _, v := range cfgs {
		io := NewIOModule(time.Duration(s.config().FlashInteval), v, s.diag, s)

		s.DeviceService.AddDevice(v.SN, io)

		err := io.Start()
		if err != nil {
			s.diag.Error("start io failed", err)
		}
		s.ios[v.SN] = io
	}

	// 注册websocket请求
	s.DispatcherBus.Register(dispatcherbus.DISPATCHER_WS_NOTIFY, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))

	return nil
}

func (s *Service) Close() error {

	for _, dev := range s.ios {
		dev.Stop()
	}

	return nil
}

func (s *Service) initWSRequestHandlers() {
	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_IO_CONTACT: s.OnWSIOContact,
		wsnotify.WS_IO_STATUS:  s.OnWSIOStatus,
		wsnotify.WS_IO_SET:     s.OnWSIOSet,
	})
}

// 获取连接状态
func (s *Service) OnWSIOStatus(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	ioStatus := device.DeviceStatus{}
	err := json.Unmarshal(byteData, &ioStatus)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	m, err := s.getIO(ioStatus.SN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	wsMsg := wsnotify.GenerateWSMsg(msg.SN, wsnotify.WS_IO_STATUS, []device.DeviceStatus{
		{
			SN:     ioStatus.SN,
			Type:   device.BaseDeviceTypeIO,
			Status: m.Status(),
		},
	})

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_IO, wsMsg)
}

// 获取io状态
func (s *Service) OnWSIOContact(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	ioContact := IoContact{}
	err := json.Unmarshal(byteData, &ioContact)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	inputs, outputs, err := s.Read(ioContact.SN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	wsMsg := wsnotify.GenerateWSMsg(msg.SN, wsnotify.WS_IO_CONTACT, IoContact{
		Src:     device.BaseDeviceTypeIO,
		Inputs:  inputs,
		Outputs: outputs,
	})

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_IO, wsMsg)
}

// 控制输出
func (s *Service) OnWSIOSet(c websocket.Connection, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	ioSet := IoSet{}
	err := json.Unmarshal(byteData, &ioSet)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	err = s.Write(ioSet.SN, ioSet.Index, ioSet.Status)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, 0, ""))
}

func (s *Service) Read(sn string) (string, string, error) {
	m, err := s.getIO(sn)
	if err != nil {
		return "", "", err
	}

	return m.Read()
}

func (s *Service) Write(sn string, index uint16, status uint16) error {
	m, err := s.getIO(sn)
	if err != nil {
		return err
	}

	return m.Write(index, status)
}

func (s *Service) getIO(sn string) (*IOModule, error) {
	m := s.ios[sn]
	if m == nil {
		return nil, errors.New("not found")
	}

	return m, nil
}

func (s *Service) OnStatus(sn string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s status:%s", sn, status))

	ioStatus := []device.DeviceStatus{
		{
			SN:     sn,
			Type:   device.BaseDeviceTypeIO,
			Status: status,
		},
	}

	s.DispatcherBus.Dispatch(dispatcherbus.DISPATCHER_DEVICE_STATUS, ioStatus)

	//io, _ := json.Marshal(wsnotify.WSMsg{
	//	Type: wsnotify.WS_IO_STATUS,
	//	Data: []device.DeviceStatus{
	//		{
	//			SN:     sn,
	//			Type:   device.BaseDeviceTypeIO,
	//			Status: status,
	//		},
	//	},
	//})
	//
	//s.NotifyService.WSSendIO(string(io))
}

func (s *Service) OnRecv(string, string) {
	s.diag.Error("OnRecv", errors.New("IO Service Not Support OnRecv"))
}

func (s *Service) OnChangeIOStatus(sn string, t string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s type:%s status:%s", sn, t, status))

	ioContact := IoContact{
		Src: device.BaseDeviceTypeIO,
		SN:  sn,
	}

	if t == IO_TYPE_INPUT {
		ioContact.Inputs = status
	} else {
		ioContact.Outputs = status
	}

	// IO数据输出状态分发
	s.DispatcherBus.Dispatch(dispatcherbus.DISPATCHER_IO, ioContact)
}
