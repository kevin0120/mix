package audi_vw

import (
	"errors"
	"fmt"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/openprotocol"
	"github.com/masami10/rush/services/storage"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/socket_writer"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

const (
	MINSEQUENCE             uint32 = 1
	MAXSEQUENCE             uint32 = 9999
	DAIL_TIMEOUT                   = time.Duration(5 * time.Second)
	MAX_KEEP_ALIVE_CHECK           = 3
	MAX_REPLY_TIMEOUT_COUNT        = 10
)

type TighteningController struct {
	device.BaseDevice
	w           *socket_writer.SocketWriter
	Srv         *Service
	StatusValue atomic.Value

	keepAliveCount    atomic.Int32
	response          chan string
	sequence          uint32 // 1~9999
	buffer            chan []byte
	Response          ResponseQueue
	muxSequence       sync.Mutex
	keepPeriod        time.Duration
	toolInfoPeriod    time.Duration
	writeTimeout      time.Duration
	recvFlag          bool
	keepaliveDeadLine atomic.Value
	closing           chan chan struct{}
	cfg               tightening_device.TighteningDeviceConfig
	protocol          string
}

func (c *TighteningController) Inputs() string {
	return ""
}

func (c *TighteningController) KeepAliveCount() int32 {
	return c.keepAliveCount.Load()
}

func (c *TighteningController) updateKeepAliveCount(i int32) {
	c.keepAliveCount.Swap(i)
}

func (c *TighteningController) addKeepAliveCount() {
	c.keepAliveCount.Inc()
}

func (c *TighteningController) Sequence() uint32 {

	c.muxSequence.Lock()
	defer c.muxSequence.Unlock()

	seq := c.sequence

	if seq == MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence++
	}

	return seq
}

func (c *TighteningController) setSequence(i uint32) {
	c.muxSequence.Lock()
	defer c.muxSequence.Unlock()
	if i >= MAXSEQUENCE {
		c.sequence = MINSEQUENCE
	} else {
		c.sequence = i
	}

}

func (c *TighteningController) Status() string {

	return c.StatusValue.Load().(string)
}

func (c *TighteningController) LoadController(controller *storage.Controllers) {}

func (c *TighteningController) updateStatus(status string) {

	if status != c.Status() {

		c.StatusValue.Store(status)

		if status == device.BaseDeviceStatusOffline {
			c.Close()

			// 断线重连
			go c.Connect()
		}

		// 将最新状态推送给hmi
		//s := wsnotify.WSStatus{
		//	SN:     c.cfg.SN,
		//	Status: string(status),
		//}

		//msg, _ := json.Marshal(s)
		//c.Srv.notifyService.WSSendControllerStatus(string(msg))

		c.Srv.diag.Debug(fmt.Sprintf("CVI3:%s %s\n", c.cfg.SN, status))

	}
}

func NewController(c Config) TighteningController {

	cont := TighteningController{
		buffer:         make(chan []byte, 1024),
		response:       make(chan string),
		closing:        make(chan chan struct{}),
		sequence:       MINSEQUENCE,
		muxSequence:    sync.Mutex{},
		keepPeriod:     time.Duration(c.KeepAlivePeriod),
		toolInfoPeriod: time.Duration(c.GetToolInfoPeriod),
		writeTimeout:   time.Duration(c.ReqTimeout),
		protocol:       tightening_device.TIGHTENING_AUDIVW,
	}

	cont.StatusValue.Store(device.BaseDeviceStatusOffline)

	return cont
}

func (c *TighteningController) SerialNumber() string {
	return c.cfg.SN
}

func (c *TighteningController) Protocol() string {
	return c.protocol
}

func (c *TighteningController) Start() error {

	c.Srv.DB.ResetTightning(c.cfg.SN)

	//c.w = socket_writer.NewSocketWriter(fmt.Sprintf("tcp://%s:%d", c.cfg.RemoteIP, c.cfg.Port), c)

	// 启动心跳检测
	//go c.keep_alive_check()

	c.Connect()

	// 订阅数据
	//c.subscribe()
	return nil
}

func (c *TighteningController) manage() {
	nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(c.keepPeriod):
			if c.Status() == device.BaseDeviceStatusOffline {
				continue
			}
			if c.KeepAliveCount() >= MAX_KEEP_ALIVE_CHECK {
				go c.updateStatus(device.BaseDeviceStatusOffline)
				c.updateKeepAliveCount(0)
				continue
			}
			if c.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.sendKeepalive()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.addKeepAliveCount()
			}
		case <-time.After(c.toolInfoPeriod):
			if c.Status() == device.BaseDeviceStatusOffline {
				continue
			}
			c.getToolInfo()

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
			nextWriteThreshold = time.Now().Add(c.writeTimeout)
		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *TighteningController) getToolInfo() {
	seq := c.Sequence()
	totalCount := fmt.Sprintf(Xml_get_total_count, c.Srv.config().Version)
	p, seq := GeneratePacket(seq, Header_type_request_with_reply, totalCount)
	c.Write([]byte(p), seq)
}

func (c *TighteningController) sendKeepalive() {
	if c.Status() == device.BaseDeviceStatusOffline {
		return
	}

	seq := c.Sequence()
	heartBeat := fmt.Sprintf(Xml_heart_beat, c.Srv.config().Version)
	keepAlivePacket, seq := GeneratePacket(seq, Header_type_keep_alive, heartBeat)
	c.Write([]byte(keepAlivePacket), seq)
}

//// 心跳检测
//func (c *TighteningController) keep_alive_check() {
//
//	for i := 0; i < MAX_KEEP_ALIVE_CHECK; i++ {
//		if c.recvFlag == true {
//			c.updateStatus(BaseDeviceStatusOnline)
//			c.recvFlag = false
//			time.Sleep(c.keepPeriod)
//
//			break
//		} else {
//			if i == (MAX_KEEP_ALIVE_CHECK - 1) {
//				c.updateStatus(BaseDeviceStatusOffline)
//			}
//		}
//
//		time.Sleep(c.keepPeriod)
//	}
//
//}

// 订阅数据
func (c *TighteningController) subscribe() {

	seq := c.Sequence()
	subscribe := fmt.Sprintf(Xml_subscribe, c.Srv.config().Version)
	subscribePacket, seq := GeneratePacket(seq, Header_type_request_with_reply, subscribe)

	c.Write([]byte(subscribePacket), seq)
}

func (c *TighteningController) Write(buf []byte, seq uint32) {
	c.buffer <- buf
}

// 异步发送
//func (c *TighteningController) manage() {
//
//	for {
//		v := <-c.buffer
//		err := c.w.Write([]byte(v))
//		if err != nil {
//			c.ProtocolService.diag.Error("Write data fail", err)
//			break
//		}
//
//		<-time.After(time.Duration(c.writeTimeout)) //300毫秒发送一次信号
//	}
//}

func (c *TighteningController) Connect() error {
	c.StatusValue.Store(device.BaseDeviceStatusOffline)
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

	c.updateStatus(device.BaseDeviceStatusOnline)

	// 启动发送
	go c.manage()

	return nil
}

func (c *TighteningController) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(c.keepPeriod))
}

func (c *TighteningController) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *TighteningController) Close() error {

	closed := make(chan struct{})
	c.closing <- closed

	<-closed
	return c.w.Close()
}

// 客户端读取
func (c *TighteningController) Read(conn net.Conn) {
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
		header := AUDIVWHeader{}
		header.Deserialize(headerStr)

		c.Response.update(header.MID, headerStr)

		if header.COD == Header_code_count_incorrect {
			c.setSequence(MINSEQUENCE)
		}

	}
}

// 拧紧抢使能
func (c *TighteningController) ToolControl(enable bool, channel int) error {
	tool_channel := ""
	if channel != 1 {
		tool_channel = fmt.Sprintf("<KNR>%d</KNR>", channel)
	}

	val_enable := 0
	if enable {
		val_enable = 1
	}

	xmlEnable := fmt.Sprintf(Xml_enable, c.Srv.config().Version, val_enable, tool_channel)

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
		time.Sleep(time.Duration(c.writeTimeout))
	}

	if header_str == "" {
		// 控制器请求失败
		return errors.New(ERR_CVI3_REPLY_TIMEOUT)
	}

	//fmt.Printf("reply_header:%s\n", header_str)
	header := AUDIVWHeader{}
	header.Deserialize(header_str)

	if !header.Check() {
		// 控制器请求失败
		return errors.New(fmt.Sprintf("%s:%s", ERR_CVI3_REPLY, request_errors[header.COD]))
	}

	return nil
}

// PSet程序设定
func (c *TighteningController) PSet(pset int, workorder_id int64, reseult_id int64, count int, user_id int64, channel int) (uint32, error) {

	//sdate, stime := utils.GetDateTime()

	tool_channel := ""
	if channel != 1 {
		tool_channel = fmt.Sprintf("<KNR>%d</KNR>", channel)
	}

	xmlPset := fmt.Sprintf(Xml_pset, c.Srv.config().Version, c.cfg.SN, workorder_id, reseult_id, count, user_id, pset, tool_channel, tool_channel, tool_channel)

	seq := c.Sequence()
	psetPacket, seq := GeneratePacket(seq, Header_type_request_with_reply, xmlPset)

	//c.response.Add(seq, "")
	c.Write([]byte(psetPacket), seq)

	c.Response.Add(seq, "")

	defer c.Response.remove(seq)

	var header_str string
	for i := 0; i < MAX_REPLY_TIMEOUT_COUNT; i++ {
		header_str = c.Response.get(seq)
		if header_str != "" {
			break
		}
		time.Sleep(time.Duration(c.writeTimeout))
	}

	if header_str == "" {
		// 控制器请求失败
		return seq, errors.New(ERR_CVI3_REPLY_TIMEOUT)
	}

	//fmt.Printf("reply_header:%s\n", header_str)
	header := AUDIVWHeader{}
	header.Deserialize(header_str)

	if !header.Check() {
		// 控制器请求失败
		return seq, errors.New(fmt.Sprintf("%s:%s", ERR_CVI3_REPLY, request_errors[header.COD]))
	}

	return seq, nil
}

func (c *TighteningController) audiVW2OPToolInfo(ti toolInfoCNT) openprotocol.ToolInfo {
	var info openprotocol.ToolInfo

	//var t tightening_device.ToolConfig

	var toolExist = false

	//for _, t = range c.cfg.Tools {
	//	//if t.ToolChannel == int(ti.MSL_MSG.KNR) {
	//	//	toolExist = true
	//	//	break
	//	//}
	//}

	if !toolExist {
		//c.Srv.diag.Error("audiVW2OPToolInfo", errors.New(fmt.Sprintf(" tool serial number:%s", t.SerialNO)))
	}

	//info.ToolSN = t.SerialNO
	info.CountSinLastService = int(ti.MSL_MSG.CSR)
	info.TotalTighteningCount = int(ti.MSL_MSG.CLT)

	return info
}

func (c *TighteningController) TryCreateMaintenance(ti toolInfoCNT) error {
	info := c.audiVW2OPToolInfo(ti)
	return c.Srv.Odoo.TryCreateMaintenance(info)
}

func (c *TighteningController) Tools() map[string]string {
	rt := map[string]string{}

	return rt
}

func (c *TighteningController) DeviceType() string {
	return tightening_device.TIGHTENING_DEVICE_TYPE_CONTROLLER
}

func (c *TighteningController) Children() map[string]device.IBaseDevice {
	return map[string]device.IBaseDevice{}
}

func (s *TighteningController) Data() interface{} {
	return nil
}

//fixme: 配置文件为空 后续强制转换都为错误
func (s *TighteningController) Config() interface{} {
	return nil
}
