package io

import (
	"encoding/json"
	"github.com/kataras/iris/websocket"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/wsnotify"
)

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
