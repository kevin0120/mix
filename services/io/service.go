package io

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherBus"
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
	DeviceService *device.Service
	IONotify
	wsnotify.WSNotify

	//DispatcherBus *DispatcherBus.Service
	dispatcherMap dispatcherBus.DispatcherMap
}

func NewService(c Config, d Diagnostic, dp Dispatcher) *Service {

	s := &Service{
		diag: d,
		DispatcherBus: dp,
		ios:  map[string]*IOModule{},
	}

	s.configValue.Store(c)

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

	s.WS.AddNotify(s)

	s.dispatcherMap = dispatcherBus.DispatcherMap{
		dispatcherBus.DISPATCHER_WS_IO_STATUS:  utils.CreateDispatchHandlerStruct(s.OnWSIOStatus),
		dispatcherBus.DISPATCHER_WS_IO_CONTACT: utils.CreateDispatchHandlerStruct(s.OnWSIOContact),
		dispatcherBus.DISPATCHER_WS_IO_SET:     utils.CreateDispatchHandlerStruct(s.OnWSIOSet),
	}
	s.DispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)

	return nil
}

func (s *Service) Close() error {
	for name, v := range s.dispatcherMap {
		s.DispatcherBus.Release(name, v.ID)
	}

	for _, dev := range s.ios {
		dev.Stop()
	}

	return nil
}

// 获取连接状态
func (s *Service) OnWSIOStatus(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	msgData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	ioStatus := IoStatus{}
	err := json.Unmarshal(msgData, &ioStatus)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	m, err := s.getIO(ioStatus.SN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	io, _ := json.Marshal(wsnotify.WSMsg{
		Type: WS_IO_STATUS,
		SN:   msg.SN,
		Data: []device.DeviceStatus{
			{
				SN:     ioStatus.SN,
				Type:   device.BaseDeviceTypeIO,
				Status: m.Status(),
			},
		},
	})

	s.WS.WSSendIO(string(io))
}

// 获取io状态
func (s *Service) OnWSIOContact(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	msgData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	ioContact := IoContact{}
	err := json.Unmarshal(msgData, &ioContact)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -1, err.Error()))
		return
	}

	inputs, outputs, err := s.Read(ioContact.SN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SN, msg.Type, -2, err.Error()))
		return
	}

	ioContacts, _ := json.Marshal(wsnotify.WSMsg{
		Type: WS_IO_CONTACT,
		Data: IoContact{
			Src:     device.BaseDeviceTypeIO,
			SN:      ioContact.SN,
			Inputs:  inputs,
			Outputs: outputs,
		},
	})

	s.WS.WSSendIO(string(ioContacts))
}

// 控制输出
func (s *Service) OnWSIOSet(data interface{}) {
	if data == nil {
		return
	}

	wsRequest := data.(*wsnotify.WSRequest)
	msgData, _ := json.Marshal(wsRequest.WSMsg.Data)
	c := wsRequest.C
	msg := wsRequest.WSMsg

	ioSet := IoSet{}
	err := json.Unmarshal(msgData, &ioSet)
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

func (s *Service) dispatchRequest(req *wsnotify.WSRequest) {
	switch req.WSMsg.Type {
	case WS_IO_CONTACT:
		s.DispatcherBus.Dispatch(dispatcherBus.DISPATCHER_WS_IO_CONTACT, req)

	case WS_IO_SET:
		s.DispatcherBus.Dispatch(dispatcherBus.DISPATCHER_WS_IO_SET, req)

	case WS_IO_STATUS:
		s.DispatcherBus.Dispatch(dispatcherBus.DISPATCHER_WS_IO_STATUS, req)
	}
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

	io, _ := json.Marshal(wsnotify.WSMsg{
		Type: wsnotify.NotifywsDeviceStatus,
		Data: []device.DeviceStatus{
			{
				SN:     sn,
				Type:   device.BaseDeviceTypeIO,
				Status: status,
			},
		},
	})

	s.WS.WSSendIO(string(io))
}

func (s *Service)OnRecv(string, string){
	s.diag.Error("OnRecv",errors.New("IO Service Not Support OnRecv"))
}

func (s *Service) OnIOStatus(sn string, t string, status string) {
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

	io, _ := json.Marshal(wsnotify.WSMsg{
		Type: WS_IO_CONTACT,
		Data: ioContact,
	})

	s.WS.WSSendIO(string(io))
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
