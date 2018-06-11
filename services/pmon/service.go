package pmon

import (
	"sync/atomic"
	"fmt"
	"github.com/pkg/errors"
	"github.com/masami10/aiis/services/httpd"
	"log"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	HTTPD		 	*httpd.Service
	configValue 	atomic.Value
	Channels 		map[string]*Channel
	err 			chan error
	diag 			Diagnostic
}

type PMONEventHandler func(error, []rune , interface{}) //事件号， 内容

func NewService(conf Config, d Diagnostic) (*Service, error)  {
	s := &Service{ err: make(chan error), diag:d}
	c, err := PmonNewConfig(conf.Path)
	if err != nil {
		return nil, err
	}
	s.configValue.Store(c)
	s.Channels = make(map[string]*Channel, len(c.Channels))// 通道长度
	connections := make(map[string]*Connection, len(c.Connections)) // 初始化长度
	for name, conn := range c.Connections {
		addr := fmt.Sprintf("udp://%s:%d",conn.Address[0],conn.Port)
		connections[name] =  NewConnection(addr, name, c.WaitResp) //waitResponse 作为其读取Timeout
		connections[name].SetDispatcher(s) //将服务注入进行通道分发
	}
	for cname,  channel:= range c.Channels {
		connectKey := fmt.Sprintf("Port%d",channel.Port)
		s.Channels[cname] = NewChannel(channel) //因为从远端传来的T/R相反，所以进行反转
		s.Channels[cname].SetConnection(connections[connectKey])
		connections[connectKey].AppendChannel(cname,channel.SNoT, channel.SNoR )
	}
	return s, nil
}

func (s *Service)Config() PmonConfig {
	return s.configValue.Load().(PmonConfig)
}

func (s *Service) Open()  error{
//func (s *Service) PmonInit( e PMONEventHandler, traceLevel int )  error{
	for _, ch := range s.Channels {
		err := ch.Start()
		if err != nil {
			return errors.Wrap(err, "Open connection fail")
		}
	}
	go s.run()
	return nil
}

func (s *Service) Close()  error{
	return nil
}

func (s *Service) PmonRegistryEvent(e PMONEventHandler, channelNumber string, ud interface{})  error{
	if _, ok := s.Channels[channelNumber]; !ok {
		log.Printf("not found channel %s", channelNumber)
		return nil
	}
	ch := s.Channels[channelNumber]
	ch.RegistryHandler(e, ud)
	return nil
}

func (s *Service) run() {
	for {
		err := <- s.err
		log.Printf("fail %s", err)
	}
}

func (s *Service) SendPmonMessage( msgType PMONSMGTYPE , channelNumber string , data string)  error{
	if _, ok := s.Channels[channelNumber]; !ok {
		log.Printf("not found channel %s", channelNumber)
		return nil
	}
	ch := s.Channels[channelNumber]
	x, err := ch.PMONGenerateMsg(msgType, data)
	if err != nil {
		log.Printf("Generation %s msg fail", msgType)
		return errors.Wrap(err, "SendPmonMessage")
	}
	ch.Write([]byte(x), msgType)
	return nil
}

func (s *Service) SendData( msgId int ,channelNumber string , data string )  error{
	return nil
}


func (s *Service) Dispatch(pkg PmonPackage, chName string) {
	if _, ok := s.Channels[chName]; !ok {
		log.Printf("not found channel %s", chName)
		return
	}
	s.Channels[chName].recvBuf <- pkg //将数据发送到通道中
}
