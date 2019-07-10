package io

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/wsnotify"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	configValue atomic.Value
	ios         map[string]*IOModule
	diag        Diagnostic
	WS          *wsnotify.Service
	IONotify
	wsnotify.WSNotify
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

	cfgs := s.config().IOS
	for _, v := range cfgs {
		s.ios[v.SN] = &IOModule{
			cfg: &v,
		}

		err := s.ios[v.SN].Start(s)
		if err != nil {
			s.diag.Error("start io failed", err)
		}
	}

	s.WS.AddNotify(s)

	return nil
}

func (s *Service) Close() error {

	return nil
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
		Type: WS_IO_STATUS,
		Data: IoStatus{
			SN:     sn,
			Status: status,
		},
	})

	s.WS.WSSendIO(string(io))
}

func (s *Service) OnIOStatus(sn string, t string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s type:%s status:%s", sn, t, status))

	io, _ := json.Marshal(wsnotify.WSMsg{
		Type: WS_IO_CONTACT,
		Data: IoContact{
			SN:      sn,
			Type:    t,
			CONTACT: status,
		},
	})

	s.WS.WSSendIO(string(io))
}

func (s *Service) OnWSMsg(data []byte) {
	msg := wsnotify.WSMsg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		s.diag.Error(string(data), err)
		return
	}

	msgData, _ := json.Marshal(msg.Data)
	switch msg.Type {
	case WS_IO_SET:
		// 控制输出
		ioSet := IoSet{}
		err := json.Unmarshal(msgData, &ioSet)
		if err != nil {
			s.diag.Error(string(msgData), err)
			return
		}

		err = s.Write(ioSet.SN, ioSet.Index, ioSet.Status)
		if err != nil {
			s.diag.Error(err.Error(), err)
			return
		}

	case WS_IO_CONTACT:
		// 获取io状态
		ioContact := IoContact{}
		err := json.Unmarshal(msgData, &ioContact)
		if err != nil {
			s.diag.Error(string(msgData), err)
			return
		}

		inputs, outputs, err := s.Read(ioContact.SN)
		if err != nil {
			s.diag.Error(err.Error(), err)
			return
		}

		ioInputs, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_IO_CONTACT,
			Data: IoContact{
				SN:      ioContact.SN,
				Type:    IO_TYPE_INPUT,
				CONTACT: inputs,
			},
		})

		s.WS.WSSendIO(string(ioInputs))

		ioOutputs, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_IO_CONTACT,
			Data: IoContact{
				SN:      ioContact.SN,
				Type:    IO_TYPE_OUTPUT,
				CONTACT: outputs,
			},
		})

		s.WS.WSSendIO(string(ioOutputs))

	case WS_IO_STATUS:
		// 获取连接状态
		ioStatus := IoStatus{}
		err := json.Unmarshal(msgData, &ioStatus)
		if err != nil {
			s.diag.Error(string(msgData), err)
			return
		}

		m, err := s.getIO(ioStatus.SN)
		if err != nil {
			s.diag.Error(err.Error(), err)
			return
		}

		io, _ := json.Marshal(wsnotify.WSMsg{
			Type: WS_IO_STATUS,
			Data: IoStatus{
				SN:     ioStatus.SN,
				Status: m.Status(),
			},
		})

		s.WS.WSSendIO(string(io))

	default:
		// 类型错误
		return
	}
}
