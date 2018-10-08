package pmon

import (
	"fmt"
	"github.com/masami10/aiis/services/httpd"
	"github.com/masami10/aiis/utils"
	"github.com/pkg/errors"
	"log"
	"sync/atomic"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	Opened      bool
	rawConf     atomic.Value
	HTTPD       *httpd.Service
	configValue atomic.Value
	Channels    map[string]*Channel
	err         chan error
	stop        chan chan struct{}
	diag        Diagnostic
}

type PMONEventHandler func(error, []rune, interface{}) //事件号， 内容

func NewService(conf Config, d Diagnostic) (*Service, error) {
	s := &Service{
		err:    make(chan error, 1),
		diag:   d,
		Opened: false,
		stop:   make(chan chan struct{}),
	}

	s.rawConf.Store(conf)
	return s, nil
}

func (s *Service) Config() PmonConfig {
	return s.configValue.Load().(PmonConfig)
}

func (s *Service) Open() error {
	defer func() {
		go s.run()
	}()
	conf := s.rawConf.Load().(Config)
	if !conf.Enable {
		return nil
	}
	exist, err := utils.FileIsExist(conf.Path)
	if !exist {
		return fmt.Errorf("Pmon Configuration path %s not exist ", conf.Path)
	}
	c, err := PmonNewConfig(conf.Path)
	if err != nil {
		return err
	}
	s.configValue.Store(c)
	s.Channels = make(map[string]*Channel, len(c.Channels))         // 通道长度
	connections := make(map[string]*Connection, len(c.Connections)) // 初始化长度
	for name, conn := range c.Connections {
		addr := fmt.Sprintf("udp://%s:%d", conn.Address[0], conn.Port)
		connections[name] = NewConnection(addr, name, c.WaitResp, conf.Workers) //waitResponse 作为其读取Timeout
		connections[name].SetDispatcher(s)                                      //将服务注入进行通道分发
	}
	for cname, channel := range c.Channels {
		connectKey := fmt.Sprintf("Port%d", channel.Port)
		s.Channels[cname] = NewChannel(channel) //因为从远端传来的T/R相反，所以进行反转
		s.Channels[cname].SetConnection(connections[connectKey])
		if _, ok := s.Config().RestartPoints[cname]; ok {
			s.Channels[cname].RestartPoint = s.Config().RestartPoints[cname]
		} else {
			s.Channels[cname].RestartPoint = ""
		}

		s.Channels[cname].Service = s
		connections[connectKey].AppendChannel(cname, channel.SNoT, channel.SNoR)
	}
	for _, ch := range s.Channels {
		err := ch.Start()
		if err != nil {
			return errors.Wrap(err, "Open connection fail")
		}
	}
	s.Opened = true
	return nil
}

func (s *Service) Close() error {
	if !s.Opened {
		return nil
	}
	for _, ch := range s.Channels {
		err := ch.Stop()
		if err != nil {
			return err
		}
	}
	stopping := make(chan struct{})
	s.stop <- stopping

	<-stopping
	return nil
}

func (s *Service) Err() <-chan error {
	return s.err
}

func (s *Service) PmonRegistryEvent(e PMONEventHandler, channelNumber string, ud interface{}) error {
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
		select {
		case err := <-s.err:
			log.Printf("Pmon Service error msg %s", err)
		case s := <-s.stop:
			close(s)
		}
	}
}

func (s *Service) SendPmonMessage(msgType PMONSMGTYPE, channelNumber string, data string) error {
	if _, ok := s.Channels[channelNumber]; !ok {
		log.Printf("not found channel %s", channelNumber)
		return fmt.Errorf("Cannot Found the channel: %s\n", channelNumber)
	}
	ch := s.Channels[channelNumber]
	x, err := ch.PMONGenerateMsg(msgType, data)
	if err != nil {
		log.Printf("Generation %s msg fail", msgType)
		return errors.Wrap(err, "SendPmonMessage")
	}
	for _, s := range x {
		if err := ch.Write([]byte(s), msgType); err != nil {
			return errors.Wrap(err, "SendPmonMessage Write")
		}
	}

	return nil
}

func (s *Service) SendData(msgId int, channelNumber string, data string) error {
	return nil
}

func (s *Service) Dispatch(pkg PmonPackage, chName string) {
	if _, ok := s.Channels[chName]; !ok {
		log.Printf("not found channel %s", chName)
		return
	}
	s.Channels[chName].recvBuf <- pkg //将数据发送到通道中
}
