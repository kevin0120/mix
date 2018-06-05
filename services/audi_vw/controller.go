package audi_vw

import (
	"github.com/masami10/rush/socket_writer"

	"fmt"
	"time"
	"sync"
)

type ControllerStatus	string

const (
	MINSEQUENCE uint = 1
	MAXSEQUENCE uint = 9999
)

const (
	STATUS_ONLINE  ControllerStatus = "online"
	STATUS_OFFLINE ControllerStatus = "offline"
)


type Controller struct {
	w 			*socket_writer.SocketWriter
	Srv 		*Service
	Status 		ControllerStatus
	sequence	uint // 1~9999
	buffer  	chan string
	Response    map[uint] []byte
	mux_serial		*sync.Mutex
	keep_period		time.Duration
	req_timeout		time.Duration
}


func (c *Controller)GeneratePacket( typ uint, xmlpacket string) (string, uint) {
	header := CVI3Header{}
	header.Init()
	header.MID = c.get_sequence()
	header.SIZ = len(xmlpacket)
	header.TYP = typ
	header_str := header.Serialize()

	return fmt.Sprintf("%s%s", header_str, xmlpacket), header.MID
}

func (c *Controller) get_sequence() uint {
	c.mux_serial.Lock()
	defer c.mux_serial.Unlock()
	s := c.sequence
	if c.sequence == MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence++
	}
	return s
}


func NewController(c Config) Controller {

	return Controller{
		w: socket_writer.NewSocketWriter(),
		buffer: make(chan string, 1), // 创建1个长度的通道，保证1收1发
		Status: STATUS_OFFLINE,
		sequence: MINSEQUENCE,
		mux_serial: new(sync.Mutex),
		keep_period: time.Duration(c.KeepAlivePeriod),
		req_timeout: time.Duration(c.ReqTimeout),
	}
}

func (c *Controller) Start()  {
	go c.keepAlive()

	go c.manage()
}

func (c *Controller) keepAlive() {

	for {
		<- time.After(time.Duration(c.keep_period)) // 周期性发送一次信号
		if c.Status == STATUS_OFFLINE {
			break
		}

		keepAlivePacket, seq := c.GeneratePacket(Header_type_request_with_reply, Xml_heart_beat)
		c.Write(keepAlivePacket,seq)

	}

}

func (s *Controller) Write(data string, seq uint) {

	s.buffer <- data
	s.Response[seq] = []byte{}
}

func (s *Controller) manage() error {

	for {
			<- time.After(time.Duration(s.req_timeout)) //3000毫秒发送一次信号
			v := <- s.buffer
			err := s.w.Write([]byte(v))
			if err != nil {
				s.Srv.diag.Error("Write data fail", err)
		}
	}

}


func (s *Controller) Connect() error {
	return s.w.Connect()
}


func (s *Controller) Close() error {
	return s.w.Close()
}
