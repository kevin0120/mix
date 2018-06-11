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
	recv_flag		bool
	cfg				controller.Config
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

		c.Srv.diag.Debug(fmt.Sprintf("civ3:%s %s\n", c.cfg.SN, c.Status))

		if c.Status == STATUS_OFFLINE {
			c.Close()

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

		seq := c.get_sequence()
		keepAlivePacket, seq := GeneratePacket(seq, Header_type_keep_alive, Xml_heart_beat)
		c.Write([]byte(keepAlivePacket), seq)

		<- time.After(time.Duration(c.keep_period)) // 周期性发送一次信号
	}
}

// 订阅数据
func (c *Controller) subscribe() {

	sdate, stime := utils.GetDateTime()
	xml_subscribe := fmt.Sprintf(Xml_subscribe, sdate, stime)

	seq := c.get_sequence()
	subscribe_packet, seq := GeneratePacket(seq, Header_type_request_with_reply, xml_subscribe)

	c.Write([]byte(subscribe_packet), seq)

}

func (c *Controller) Write(buf []byte, seq uint) {
	c.buffer <- buf
}

// 异步发送
func (c *Controller) manage() {

	for {
		v := <- c.buffer
		err := c.w.Write([]byte(v))
		if err != nil {
			c.Srv.diag.Error("Write data fail", err)
			break
		}

		<- time.After(time.Duration(c.req_timeout)) //300毫秒发送一次信号
	}
}


func (c *Controller) Connect() error {
	c.Status = STATUS_OFFLINE
	c.sequence = 0

	c.Response = ResponseQueue {
		Results: map[uint]string{},
	}

	c.Srv.diag.Debug(fmt.Sprintf("CVI3:%s connecting ...\n", c.cfg.SN))

	for {
		err := c.w.Connect(DAIL_TIMEOUT)
		if err != nil {
			c.Srv.diag.Error("connect err", err)
		} else {
			break
		}

		time.Sleep(time.Duration(c.req_timeout))
	}

	c.update_status(STATUS_ONLINE)

	// 启动发送
	go c.manage()

	// 启动心跳
	go c.keepAlive()

	return nil
}


func (c *Controller) Close() error {
	return c.w.Close()
}

// 客户端读取
func (c *Controller) Read(conn net.Conn){
	defer conn.Close()

	buffer := make([]byte, c.Srv.config().ReadBufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		c.recv_flag = true

		msg := string(buffer[0:n])

		//fmt.Printf("%s\n", string(buffer))

		// 处理应答
		header_str := msg[0: HEADER_LEN]
		header := CVI3Header{}
		header.Deserialize(header_str)

		c.Response.update(header.MID, header_str)

	}
}

// PSet程序设定
func (c *Controller) PSet(pset int, workorder_id int64, reseult_id int64, count int) (uint, error) {

	sdate, stime := utils.GetDateTime()
	xml_pset := fmt.Sprintf(Xml_pset, sdate, stime, c.cfg.SN, workorder_id, reseult_id, count, pset)

	seq := c.get_sequence()
	pset_packet, seq := GeneratePacket(seq, Header_type_request_with_reply, xml_pset)

	c.Response.Add(seq, "")
	c.Write([]byte(pset_packet), seq)

	return seq, nil
}