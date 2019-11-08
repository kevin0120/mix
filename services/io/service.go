package io

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	configValue   atomic.Value
	ios           map[string]*IOModule
	diag          Diagnostic
	WS            *wsnotify.Service
	DeviceService *device.Service
	IONotify
	wsnotify.WSNotify

	requestDispatchers map[string]*utils.Dispatcher
}

func NewService(c Config, d Diagnostic) *Service {

	s := &Service{
		diag: d,
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
		s.ios[v.SN] = &IOModule{
			cfg:           v,
			flashInterval: time.Duration(s.config().FlashInteval),
		}

		s.DeviceService.AddDevice(v.SN, s.ios[v.SN])

		err := s.ios[v.SN].Start(s)
		if err != nil {
			s.diag.Error("start io failed", err)
		}
	}

	s.WS.AddNotify(s)

	s.initRequestDispatchers()

	return nil
}

func (s *Service) Close() error {
	for _, v := range s.requestDispatchers {
		v.Release()
	}

	for _, dev := range s.ios {
		dev.Stop()
	}

	return nil
}

func (s *Service) initRequestDispatchers() {
	s.requestDispatchers = map[string]*utils.Dispatcher{
		WS_IO_STATUS:  utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_IO_CONTACT: utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
		WS_IO_SET:     utils.CreateDispatcher(utils.DEFAULT_BUF_LEN),
	}

	s.requestDispatchers[WS_IO_STATUS].Register(s.OnWS_IO_STATUS)
	s.requestDispatchers[WS_IO_CONTACT].Register(s.OnWS_IO_CONTACT)
	s.requestDispatchers[WS_IO_SET].Register(s.OnWS_IO_SET)

	for _, v := range s.requestDispatchers {
		v.Start()
	}
}

// 获取连接状态
func (s *Service) OnWS_IO_STATUS(data interface{}) {
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
				Type:   device.DEVICE_TYPE_IO,
				Status: m.Status(),
			},
		},
	})

	s.WS.WSSendIO(string(io))
}

// 获取io状态
func (s *Service) OnWS_IO_CONTACT(data interface{}) {
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
			Src:     device.DEVICE_TYPE_IO,
			SN:      ioContact.SN,
			Inputs:  inputs,
			Outputs: outputs,
		},
	})

	s.WS.WSSendIO(string(ioContacts))
}

// 控制输出
func (s *Service) OnWS_IO_SET(data interface{}) {
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
	d, exist := s.requestDispatchers[req.WSMsg.Type]
	if exist {
		d.Dispatch(req)
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
		Type: device.WS_DEVICE_STATUS,
		Data: []device.DeviceStatus{
			{
				SN:     sn,
				Type:   device.DEVICE_TYPE_IO,
				Status: status,
			},
		},
	})

	s.WS.WSSendIO(string(io))
}

func (s *Service) OnIOStatus(sn string, t string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s type:%s status:%s", sn, t, status))

	ioContact := IoContact{
		Src: device.DEVICE_TYPE_IO,
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
