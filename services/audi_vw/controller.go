package audi_vw

import (
	"github.com/masami10/rush/socket_writer"
	"github.com/masami10/rush/utils"
	"fmt"
	"time"
	"sync"
	"net"
	"github.com/masami10/rush/services/controller"
)

type ControllerStatus	string

const (
	MINSEQUENCE uint = 1
	MAXSEQUENCE uint = 9999
	DAIL_TIMEOUT = time.Duration(300 * time.Millisecond)
	MAX_KEEP_ALIVE_CHECK = 3
	READ_BUF_LEN = 65535
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
	buffer  	chan []byte
	Response    ResponseQueue
	mtx_serial		*sync.Mutex
	keep_period		time.Duration
	req_timeout		time.Duration
	mtx_status		sync.Mutex
	RemoteConn		net.Conn
	recv_flag		bool
	cfg				controller.Config

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
	c.mtx_serial.Lock()
	defer c.mtx_serial.Unlock()
	s := c.sequence
	if c.sequence == MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence++
	}
	return s
}

func (c *Controller) GetStatus() (ControllerStatus) {
	defer c.mtx_status.Unlock()

	c.mtx_status.Lock()
	return c.Status
}

func (c *Controller) update_status(status ControllerStatus) {
	defer c.mtx_status.Unlock()

	c.mtx_status.Lock()

	if status != c.Status {
		c.Status = status
		// 将最新状态推送给hmi
		//go client.Parent.FUNCStatus(client.Config.SN, client.Status)
		//fmt.Printf("civ3:%s %s\n", client.Config.SN, client.Status)

		if c.Status == STATUS_OFFLINE {
			c.Close()
			c.RemoteConn.Close()

			// 断线重连
			go c.Connect()
		}
	}
}


func NewController(c Config) Controller {

	return Controller {
		buffer: make(chan []byte, 1024), // 创建1个长度的通道，保证1收1发
		Status: STATUS_OFFLINE,
		sequence: MINSEQUENCE,
		mtx_serial: new(sync.Mutex),
		keep_period: time.Duration(c.KeepAlivePeriod),
		req_timeout: time.Duration(c.ReqTimeout),
		mtx_status: sync.Mutex{},
	}
}

func (c *Controller) Start()  {

	c.w = socket_writer.NewSocketWriter(c.cfg.RemoteIP)

	// 启动心跳检测
	go c.keep_alive_check()

	c.Connect()

	// 订阅数据
	c.subscribe()
}

// 心跳检测
func (c *Controller) keep_alive_check() {

	for {

		for i:=0; i < MAX_KEEP_ALIVE_CHECK; i++ {
			if c.recv_flag == true {
				c.update_status(STATUS_ONLINE)
				c.recv_flag = false
				time.Sleep(c.keep_period)

				break
			} else {
				if i == (MAX_KEEP_ALIVE_CHECK - 1) {
					c.update_status(STATUS_OFFLINE)
				}
			}

			time.Sleep(c.keep_period)
		}

	}
}

func (c *Controller) keepAlive() {

	for {
		<- time.After(time.Duration(c.keep_period)) // 周期性发送一次信号
		if c.Status == STATUS_OFFLINE {
			break
		}

		keepAlivePacket, seq := c.GeneratePacket(Header_type_request_with_reply, Xml_heart_beat)
		c.Write([]byte(keepAlivePacket), seq)
	}
}

// 订阅数据
func (c *Controller) subscribe() {

	sdate, stime := utils.GetDateTime()
	xml_subscribe := fmt.Sprintf(Xml_subscribe, sdate, stime)

	subscribe_packet, seq := c.GeneratePacket(Header_type_request_with_reply, xml_subscribe)

	c.Write([]byte(subscribe_packet), seq)

}

func (s *Controller) Write(buf []byte, seq uint) {

	s.buffer <- buf
	s.Response.update(seq, buf)
}

// 异步发送
func (s *Controller) manage() {

	for {
		<- time.After(time.Duration(s.req_timeout)) //300毫秒发送一次信号

		v := <- s.buffer
		err := s.w.Write([]byte(v))
		if err != nil {
			s.Srv.diag.Error("Write data fail", err)
			break
		}
	}
}


func (s *Controller) Connect() error {
	s.Status = STATUS_OFFLINE
	s.sequence = 0

	for {
		err := s.w.Connect(DAIL_TIMEOUT)
		if err != nil {
			return err
		} else {
			break
		}

		time.Sleep(time.Duration(s.req_timeout))
	}

	s.update_status(STATUS_ONLINE)
	s.Response = ResponseQueue {
		Results: map[uint][]byte{},
	}

	// 启动发送
	go s.manage()

	// 启动心跳
	go s.keepAlive()

	return nil
}


func (s *Controller) Close() error {
	return s.w.Close()
}

// 客户端读取
func (s *Controller) ClientRead(conn net.Conn){
	defer conn.Close()

	buffer := make([]byte, READ_BUF_LEN)

	for {
		//msg, err := reader.ReadString('\n')
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		s.recv_flag = true

		msg := string(buffer[0:n])

		//fmt.Printf("%s\n", msg)

		// 处理应答
		header_str := msg[0: HEADER_LEN]
		header := CVI3Header{}
		header.Deserialize(header_str)

		s.Response.update(header.MID, []byte(header_str))
	}
}

// PSet程序设定
func (s *Controller) PSet(pset int, workorder_id int, reseult_id int, count int) (uint, error) {

	sdate, stime := utils.GetDateTime()
	xml_pset := fmt.Sprintf(Xml_pset, sdate, stime, s.cfg.SN, workorder_id, reseult_id, count, pset)

	pset_packet, seq := s.GeneratePacket(Header_type_request_with_reply, xml_pset)

	s.Write([]byte(pset_packet), seq)

	return seq, nil
}