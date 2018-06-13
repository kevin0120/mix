package audi_vw

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/socket_writer"
	"github.com/masami10/rush/utils"
	"net"
	"sync"
	"time"
	"sync/atomic"
)

type ControllerStatusType string

const (
	MINSEQUENCE          uint = 1
	MAXSEQUENCE          uint = 9999
	DAIL_TIMEOUT              = time.Duration(5 * time.Second)
	MAX_KEEP_ALIVE_CHECK      = 3
)

const (
	STATUS_ONLINE  ControllerStatusType = "online"
	STATUS_OFFLINE ControllerStatusType = "offline"
)

type Controller struct {
	w                 *socket_writer.SocketWriter
	Srv               *Service
	StatusValue       atomic.Value

	keepAliveCount	  atomic.Value
	response		  chan string
	sequence          uint // 1~9999
	buffer            chan []byte
	//Response          ResponseQueue
	mux_seq           sync.Mutex
	keep_period       time.Duration
	req_timeout       time.Duration
	recv_flag         bool
	keepaliveDeadLine atomic.Value
	closing           chan chan struct{}
	cfg               controller.Config
}

func (c *Controller) KeepAliveCount() int {
	return c.keepAliveCount.Load().(int)
}

func (c *Controller) updateKeepAliveCount(i int) {
	c.keepAliveCount.Store(i)
}

func (c *Controller) Sequence() uint {
	c.mux_seq.Lock()
	defer c.mux_seq.Unlock()

	if c.sequence == MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence++
	}
	return c.sequence
}

func (c *Controller) Status() ControllerStatusType {

	return c.StatusValue.Load().(ControllerStatusType)
}

func (c *Controller) updateStatus(status ControllerStatusType) {

	if status != c.Status() {

		c.StatusValue.Store(status)

		if status == STATUS_OFFLINE {
			c.Close()

			// 断线重连
			go c.Connect()
		}

		// 将最新状态推送给hmi
		s := wsnotify.WSStatus{
			SN:     c.cfg.SN,
			Status: string(status),
		}

		msg, _ := json.Marshal(s)
		c.Srv.WS.WSSendControllerStatus(string(msg))

		c.Srv.diag.Debug(fmt.Sprintf("civ3:%s %s\n", c.cfg.SN, status))

	}
}

func NewController(c Config) Controller {

	cont := Controller{
		buffer:      make(chan []byte, 1024),
		response: 	 make(chan string),
		closing:     make(chan chan struct{}),
		sequence:    MINSEQUENCE,
		mux_seq:     sync.Mutex{},
		keep_period: time.Duration(c.KeepAlivePeriod),
		req_timeout: time.Duration(c.ReqTimeout),
	}

	cont.StatusValue.Store(STATUS_OFFLINE)

	return cont
}

func (c *Controller) Start() {

	c.w = socket_writer.NewSocketWriter(fmt.Sprintf("tcp://%s:%d", c.cfg.RemoteIP, c.cfg.Port), c)

	// 启动心跳检测
	//go c.keep_alive_check()

	c.Connect()

	// 订阅数据
	c.subscribe()
}

func (c *Controller) manage()  {
	for {
		select {
		case <- time.After(c.keep_period):
			count := c.KeepAliveCount()
			if c.Status() == STATUS_OFFLINE {
				continue
			}
			if  count >= MAX_KEEP_ALIVE_CHECK {
				go c.updateStatus(STATUS_OFFLINE)
				c.updateKeepAliveCount(0)
				continue
			}
			if c.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.send_keepalive()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.updateKeepAliveCount(count+1)
			}
		case v := <-c.buffer:
			err := c.w.Write([]byte(v))
			if err != nil {
				c.Srv.diag.Error("Write data fail", err)
			}
			time.Sleep(c.req_timeout)
		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *Controller) send_keepalive(){
	if c.Status() == STATUS_OFFLINE {
		return
	}

	seq := c.Sequence()
	keepAlivePacket, seq := GeneratePacket(seq, Header_type_keep_alive, Xml_heart_beat)
	c.Write([]byte(keepAlivePacket), seq)
}

//// 心跳检测
//func (c *Controller) keep_alive_check() {
//
//	for i := 0; i < MAX_KEEP_ALIVE_CHECK; i++ {
//		if c.recv_flag == true {
//			c.updateStatus(STATUS_ONLINE)
//			c.recv_flag = false
//			time.Sleep(c.keep_period)
//
//			break
//		} else {
//			if i == (MAX_KEEP_ALIVE_CHECK - 1) {
//				c.updateStatus(STATUS_OFFLINE)
//			}
//		}
//
//		time.Sleep(c.keep_period)
//	}
//
//}

// 订阅数据
func (c *Controller) subscribe() {

	sdate, stime := utils.GetDateTime()
	xml_subscribe := fmt.Sprintf(Xml_subscribe, sdate, stime)

	seq := c.Sequence()
	subscribe_packet, seq := GeneratePacket(seq, Header_type_request_with_reply, xml_subscribe)

	c.Write([]byte(subscribe_packet), seq)

}

func (c *Controller) Write(buf []byte, seq uint) {
	c.buffer <- buf
}

// 异步发送
//func (c *Controller) manage() {
//
//	for {
//		v := <-c.buffer
//		err := c.w.Write([]byte(v))
//		if err != nil {
//			c.Srv.diag.Error("Write data fail", err)
//			break
//		}
//
//		<-time.After(time.Duration(c.req_timeout)) //300毫秒发送一次信号
//	}
//}

func (c *Controller) Connect() error {
	c.StatusValue.Store(STATUS_OFFLINE)
	c.sequence = 0

	//c.Response = ResponseQueue{
	//	Results: map[uint]string{},
	//}

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

	c.updateStatus(STATUS_ONLINE)


	// 启动发送
	go c.manage()

	return nil
}

func (c *Controller) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(c.keep_period))
}

func (c *Controller) KeepAliveDeadLine() time.Time{
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *Controller) Close() error {

	closed := make(chan struct{})
	c.closing <- closed

	<-closed
	return c.w.Close()
}

// 客户端读取
func (c *Controller) Read(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, c.Srv.config().ReadBufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		c.updateKeepAliveDeadLine()

		c.updateKeepAliveCount(0)

		msg := string(buffer[0:n])

		//fmt.Printf("%s\n", string(buffer))

		// 处理应答
		header_str := msg[0:HEADER_LEN]
		header := CVI3Header{}
		header.Deserialize(header_str)

		//c.response <- header_str

		//c.Response.update(header.MID, header_str)

	}
}

// PSet程序设定
func (c *Controller) PSet(pset int, workorder_id int64, reseult_id int64, count int, user_id int64) (uint, error) {

	sdate, stime := utils.GetDateTime()
	xml_pset := fmt.Sprintf(Xml_pset, sdate, stime, c.cfg.SN, workorder_id, reseult_id, count, user_id, pset)

	seq := c.Sequence()
	pset_packet, seq := GeneratePacket(seq, Header_type_request_with_reply, xml_pset)

	//c.Response.Add(seq, "")
	c.Write([]byte(pset_packet), seq)

	return seq, nil
}
