package audi_vw

import (
	"github.com/masami10/rush/socket_writer"
	"github.com/masami10/rush/utils"
	"fmt"
	"time"
	"sync"
	"net"
	"github.com/masami10/rush/services/controller"
	"encoding/json"
	"github.com/masami10/rush/services/wsnotify"
)

type ControllerStatusType	string

const (
	MINSEQUENCE uint = 1
	MAXSEQUENCE uint = 9999
	DAIL_TIMEOUT = time.Duration(1 * time.Second)
	MAX_KEEP_ALIVE_CHECK = 3
)

const (
	STATUS_ONLINE  ControllerStatusType = "online"
	STATUS_OFFLINE ControllerStatusType = "offline"
)


type Controller struct {
	w 				*socket_writer.SocketWriter
	Srv 			*Service
	Status 			ControllerStatusType
	sequence		uint // 1~9999
	buffer  		chan []byte
	Response    	ResponseQueue
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

	if c.sequence == MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence++
	}
	return c.sequence
}

func (c *Controller) GetStatus() (ControllerStatusType) {
	defer c.mtx_status.Unlock()

	c.mtx_status.Lock()
	return c.Status
}

func (c *Controller) update_status(status ControllerStatusType) {
	defer c.mtx_status.Unlock()

	c.mtx_status.Lock()

	if status != c.Status {
		c.Status = status

		// 将最新状态推送给hmi
		s := wsnotify.WSStatus{
			SN: c.cfg.SN,
			Status: string(status),
		}

		msg, _ := json.Marshal(s)
		c.Srv.WS.WSSendControllerStatus(string(msg))

		fmt.Printf("civ3:%s %s\n", c.cfg.SN, c.Status)

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
		buffer: make(chan []byte, 1024),
		Status: STATUS_OFFLINE,
		sequence: MINSEQUENCE,
		mtx_serial: new(sync.Mutex),
		keep_period: time.Duration(c.KeepAlivePeriod),
		req_timeout: time.Duration(c.ReqTimeout),
		mtx_status: sync.Mutex{},
	}
}

func (c *Controller) Start()  {

	c.w = socket_writer.NewSocketWriter(fmt.Sprintf("tcp://%s:%d", c.cfg.RemoteIP, c.cfg.Port), c)

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
		if c.Status == STATUS_OFFLINE {
			break
		}

		keepAlivePacket, seq := c.GeneratePacket(Header_type_keep_alive, Xml_heart_beat)
		c.Write([]byte(keepAlivePacket), seq)

		<- time.After(time.Duration(c.keep_period)) // 周期性发送一次信号
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
}

// 异步发送
func (s *Controller) manage() {

	for {
		v := <- s.buffer
		err := s.w.Write([]byte(v))
		if err != nil {
			s.Srv.diag.Error("Write data fail", err)
			break
		}

		<- time.After(time.Duration(s.req_timeout)) //300毫秒发送一次信号
	}
}


func (s *Controller) Connect() error {
	s.Status = STATUS_OFFLINE
	s.sequence = 0

	s.Response = ResponseQueue {
		Results: map[uint]string{},
	}

	fmt.Printf("CVI3:%s connecting ...\n", s.cfg.SN)

	for {
		err := s.w.Connect(DAIL_TIMEOUT)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		} else {
			break
		}

		time.Sleep(time.Duration(s.req_timeout))
	}

	s.update_status(STATUS_ONLINE)

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
func (s *Controller) Read(conn net.Conn){
	defer conn.Close()

	buffer := make([]byte, s.Srv.config().ReadBufferSize)

	for {
		//msg, err := reader.ReadString('\n')
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		s.recv_flag = true

		msg := string(buffer[0:n])

		//fmt.Printf("%s\n", string(buffer))

		// 处理应答
		header_str := msg[0: HEADER_LEN]
		header := CVI3Header{}
		header.Deserialize(header_str)

		s.Response.update(header.MID, header_str)

	}
}

// PSet程序设定
func (s *Controller) PSet(pset int, workorder_id int64, reseult_id int64, count int) (uint, error) {

	sdate, stime := utils.GetDateTime()
	xml_pset := fmt.Sprintf(Xml_pset, sdate, stime, s.cfg.SN, workorder_id, reseult_id, count, pset)

	pset_packet, seq := s.GeneratePacket(Header_type_request_with_reply, xml_pset)

	s.Response.Add(seq, "")
	s.Write([]byte(pset_packet), seq)

	return seq, nil
}