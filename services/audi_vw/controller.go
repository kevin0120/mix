package audi_vw

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/socket_writer"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MINSEQUENCE             uint32 = 1
	MAXSEQUENCE             uint32 = 9999
	DAIL_TIMEOUT                   = time.Duration(5 * time.Second)
	MAX_KEEP_ALIVE_CHECK           = 3
	MAX_REPLY_TIMEOUT_COUNT        = 10
)

type Controller struct {
	w           *socket_writer.SocketWriter
	Srv         *Service
	StatusValue atomic.Value

	keepAliveCount    int32
	response          chan string
	sequence          uint32 // 1~9999
	buffer            chan []byte
	Response          ResponseQueue
	mux_seq           sync.Mutex
	keep_period       time.Duration
	req_timeout       time.Duration
	recv_flag         bool
	keepaliveDeadLine atomic.Value
	closing           chan chan struct{}
	cfg               controller.ControllerConfig
	protocol          string
}

func (c *Controller) Inputs() string {
	return ""
}

func (c *Controller) KeepAliveCount() int32 {
	return atomic.LoadInt32(&c.keepAliveCount)
}

func (c *Controller) updateKeepAliveCount(i int32) {
	atomic.SwapInt32(&c.keepAliveCount, i)
}

func (c *Controller) addKeepAliveCount() {
	atomic.AddInt32(&c.keepAliveCount, 1)
}

func (c *Controller) Sequence() uint32 {

	c.mux_seq.Lock()
	defer c.mux_seq.Unlock()

	seq := c.sequence

	if seq == MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence++
	}

	return seq
}

func (c *Controller) setSequence(i uint32) {
	c.mux_seq.Lock()
	defer c.mux_seq.Unlock()
	if i >= MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence = i
	}

}

func (c *Controller) Status() string {

	return c.StatusValue.Load().(string)
}

func (c *Controller) LoadController(controller *storage.Controllers) {}

func (c *Controller) updateStatus(status string) {

	if status != c.Status() {

		c.StatusValue.Store(status)

		if status == controller.STATUS_OFFLINE {
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

		c.Srv.diag.Debug(fmt.Sprintf("CVI3:%s %s\n", c.cfg.SN, status))

	}
}

func NewController(c Config) Controller {

	cont := Controller{
		buffer:      make(chan []byte, 1024),
		response:    make(chan string),
		closing:     make(chan chan struct{}),
		sequence:    MINSEQUENCE,
		mux_seq:     sync.Mutex{},
		keep_period: time.Duration(c.KeepAlivePeriod),
		req_timeout: time.Duration(c.ReqTimeout),
		protocol:    controller.AUDIPROTOCOL,
	}

	cont.StatusValue.Store(controller.STATUS_OFFLINE)

	return cont
}

func (c *Controller) Protocol() string {
	return c.protocol
}

func (c *Controller) Start() {

	c.Srv.DB.ResetTightning(c.cfg.SN)

	c.w = socket_writer.NewSocketWriter(fmt.Sprintf("tcp://%s:%d", c.cfg.RemoteIP, c.cfg.Port), c)

	// 启动心跳检测
	//go c.keep_alive_check()

	c.Connect()

	// 订阅数据
	//c.subscribe()
}

func (c *Controller) manage() {
	nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(c.keep_period):
			if c.Status() == controller.STATUS_OFFLINE {
				continue
			}
			if c.KeepAliveCount() >= MAX_KEEP_ALIVE_CHECK {
				go c.updateStatus(controller.STATUS_OFFLINE)
				c.updateKeepAliveCount(0)
				continue
			}
			if c.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.sendKeepalive()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.addKeepAliveCount()
			}
		case v := <-c.buffer:
			for nextWriteThreshold.After(time.Now()) {
				time.Sleep(time.Microsecond * 100)
			}
			err := c.w.Write([]byte(v))
			if err != nil {
				c.Srv.diag.Error("Write data fail", err)
			} else {
				c.updateKeepAliveDeadLine()
			}
			nextWriteThreshold = time.Now().Add(c.req_timeout)
		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *Controller) sendKeepalive() {
	if c.Status() == controller.STATUS_OFFLINE {
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

	seq := c.Sequence()
	subscribePacket, seq := GeneratePacket(seq, Header_type_request_with_reply, Xml_subscribe)

	c.Write([]byte(subscribePacket), seq)

}

func (c *Controller) Write(buf []byte, seq uint32) {
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
	c.StatusValue.Store(controller.STATUS_OFFLINE)
	c.setSequence(MINSEQUENCE)

	c.Response = ResponseQueue{
		Results: map[uint32]string{},
	}

	c.Srv.diag.Debug(fmt.Sprintf("CVI3:%s connecting ...\n", c.cfg.SN))

	for {
		err := c.w.Connect(DAIL_TIMEOUT)
		if err != nil {
			c.Srv.diag.Error("connect err", err)
		} else {

			// 订阅数据
			c.subscribe()

			break
		}

		time.Sleep(time.Duration(c.Srv.config().KeepAlivePeriod * 3))
	}

	c.updateStatus(controller.STATUS_ONLINE)

	// 启动发送
	go c.manage()

	return nil
}

func (c *Controller) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(c.keep_period))
}

func (c *Controller) KeepAliveDeadLine() time.Time {
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

		//c.updateKeepAliveDeadLine()

		c.updateKeepAliveCount(0)

		msg := string(buffer[0:n])

		//fmt.Printf("%s\n", string(buffer))

		// 处理应答
		headerStr := msg[0:HEADER_LEN]
		header := CVI3Header{}
		header.Deserialize(headerStr)

		c.Response.update(header.MID, headerStr)
		//if c.Response.HasResponse(header.MID) {
		//	c.Response.update(header.MID, headerStr)
		//	c.response <- headerStr
		//}
	}
}

// 拧紧抢使能
func (c *Controller) ToolControl(enable bool, channel int) error {
	tool_channel := ""
	if channel != controller.DEFAULT_TOOL_CHANNEL {
		tool_channel = fmt.Sprintf("<KNR>%d</KNR>", channel)
	}

	val_enable := 0
	if enable {
		val_enable = 1
	}

	xmlEnable := fmt.Sprintf(Xml_enable, val_enable, tool_channel)

	seq := c.Sequence()
	psetPacket, seq := GeneratePacket(seq, Header_type_request_with_reply, xmlEnable)

	c.Write([]byte(psetPacket), seq)

	c.Response.Add(seq, "")

	defer c.Response.remove(seq)

	var header_str string
	for i := 0; i < MAX_REPLY_TIMEOUT_COUNT; i++ {
		header_str = c.Response.get(seq)
		if header_str != "" {
			break
		}
		time.Sleep(time.Duration(c.req_timeout))
	}

	if header_str == "" {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REPLY_TIMEOUT)
	}

	//fmt.Printf("reply_header:%s\n", header_str)
	header := CVI3Header{}
	header.Deserialize(header_str)

	if !header.Check() {
		// 控制器请求失败
		return errors.New(fmt.Sprintf("%s:%s", ERR_CVI3_REPLY, request_errors[header.COD]))
	}

	return nil
}

// PSet程序设定
func (c *Controller) PSet(pset int, workorder_id int64, reseult_id int64, count int, user_id int64, channel int) (uint32, error) {

	//sdate, stime := utils.GetDateTime()

	tool_channel := ""
	if channel != controller.DEFAULT_TOOL_CHANNEL {
		tool_channel = fmt.Sprintf("<KNR>%d</KNR>", channel)
	}

	xmlPset := fmt.Sprintf(Xml_pset, c.cfg.SN, workorder_id, reseult_id, count, user_id, pset, tool_channel, tool_channel, tool_channel)

	seq := c.Sequence()
	psetPacket, seq := GeneratePacket(seq, Header_type_request_with_reply, xmlPset)

	//c.Response.Add(seq, "")
	c.Write([]byte(psetPacket), seq)

	c.Response.Add(seq, "")

	defer c.Response.remove(seq)

	var header_str string
	for i := 0; i < MAX_REPLY_TIMEOUT_COUNT; i++ {
		header_str = c.Response.get(seq)
		if header_str != "" {
			break
		}
		time.Sleep(time.Duration(c.req_timeout))
	}

	if header_str == "" {
		// 控制器请求失败
		return seq, errors.New(ERR_CVI3_REPLY_TIMEOUT)
	}

	//fmt.Printf("reply_header:%s\n", header_str)
	header := CVI3Header{}
	header.Deserialize(header_str)

	if !header.Check() {
		// 控制器请求失败
		return seq, errors.New(fmt.Sprintf("%s:%s", ERR_CVI3_REPLY, request_errors[header.COD]))
	}

	return seq, nil
}
