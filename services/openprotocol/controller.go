package openprotocol

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/masami10/rush/services/controller"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/socket_writer"
	"net"
	"sync/atomic"
	"time"
)

const (
	DAIL_TIMEOUT         = time.Duration(5 * time.Second)
	MAX_KEEP_ALIVE_CHECK = 3
)

type handlerPkg struct {
	Header OpenProtocolHeader
	Body   string
}

type Controller struct {
	w                 *socket_writer.SocketWriter
	cfg               controller.Config
	StatusValue       atomic.Value
	keepAliveCount    int32
	keep_period       time.Duration
	req_timeout       time.Duration
	Response          ResponseQueue
	Srv               *Service
	buffer            chan []byte
	closing           chan chan struct{}
	handlerBuf        chan handlerPkg
	keepaliveDeadLine atomic.Value
	protocol          string
}

func NewController(c Config) Controller {

	cont := Controller{
		buffer:      make(chan []byte, 1024),
		closing:     make(chan chan struct{}),
		keep_period: time.Duration(c.KeepAlivePeriod),
		req_timeout: time.Duration(c.ReqTimeout),
		Response:    ResponseQueue{},
		handlerBuf:  make(chan handlerPkg, 1024),
		protocol:    controller.OPENPROTOCOL,
	}

	cont.StatusValue.Store(controller.STATUS_OFFLINE)

	return cont
}

func (c *Controller) handlerProcess() {
	for {
		select {
		case pkg := <-c.handlerBuf:
			c.HandleMsg(&pkg)
		}
	}
}

func (c *Controller) HandleMsg(pkg *handlerPkg) {
	c.Srv.diag.Debug(fmt.Sprintf("%s%s\n", pkg.Header.Serialize(), pkg.Body))

	switch pkg.Header.MID {
	case MID_0061_LAST_RESULT:
		// 处理结果

	case MID_7410_LAST_CURVE:
		// 处理波形
	}
}

func (c *Controller) Start() {

	c.w = socket_writer.NewSocketWriter(fmt.Sprintf("tcp://%s:%d", c.cfg.RemoteIP, c.cfg.Port), c)

	go c.handlerProcess()

	c.Connect()
}

func (c *Controller) Protocol() string {
	return c.protocol
}

func (c *Controller) Connect() error {
	c.StatusValue.Store(controller.STATUS_OFFLINE)

	for {
		err := c.w.Connect(DAIL_TIMEOUT)
		if err != nil {
			c.Srv.diag.Error("connect err", err)
		} else {
			break
		}

		time.Sleep(time.Duration(c.req_timeout))
	}

	c.updateStatus(controller.STATUS_ONLINE)

	c.startComm()

	c.PSetSubscribe()
	//c.CurveSubscribe()
	c.SelectorSubscribe()
	c.ResultSubcribe()
	c.JobInfoSubscribe()
	//c.DataSubscribeCurve()
	//c.IdentifierSubcribe()
	//c.PSet()
	// 启动发送
	c.lastResult()
	go c.manage()

	return nil
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

func (c *Controller) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(c.keep_period))
}

func (c *Controller) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *Controller) Status() string {

	return c.StatusValue.Load().(string)
}

func (c *Controller) sendKeepalive() {
	if c.Status() == controller.STATUS_OFFLINE {
		return
	}

	keep_alive := GeneratePackage(MID_9999_ALIVE, DEFAULT_REV, "", DEFAULT_MSG_END)
	c.Write([]byte(keep_alive))
}

func (c *Controller) startComm() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	start := GeneratePackage(MID_0001_START, "003", "", DEFAULT_MSG_END)

	//c.Response.Add(MID_0001_START, "")

	c.Write([]byte(start))

	return nil
}

func (c *Controller) Write(buf []byte) {
	c.buffer <- buf
}

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

func (c *Controller) Read(conn net.Conn) {
	defer conn.Close()

	len_header := LEN_HEADER
	rest := 0
	body := ""
	header_rest := 0

	var header_buffer string
	var header OpenProtocolHeader

	buffer := make([]byte, c.Srv.config().ReadBufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		c.updateKeepAliveCount(0)

		msg := string(buffer[0:n])

		off := 0 //循环前偏移为0
		for off < n {
			if rest == 0 {
				len_msg := n - off
				if len_msg < len_header-header_rest {
					//长度不够
					if header_rest == 0 {
						header_rest = len_header - len_msg
					} else {
						header_rest -= len_msg
					}
					header_buffer += msg[off : off+len_msg]
					break
				} else {
					//完整
					if header_rest == 0 {
						header_buffer = msg[off : off+len_header]
						off += len_header
					} else {
						header_buffer += msg[off : off+header_rest]
						off += header_rest
						header_rest = 0
					}
				}
				//fmt.Printf("header rest:%d, offset:%d, n %d, header : %s\n", header_rest, off, n, header_buffer)
				header.Deserialize(header_buffer)
				header_buffer = ""
				if n-off > header.LEN {
					//粘包
					body = msg[off : off+header.LEN]

					//c.Parse(body)
					pkg := handlerPkg{
						Header: header,
						Body:   body,
					}
					c.handlerBuf <- pkg
					off += header.LEN + 1
					rest = 0 //同样解析头

				} else {
					body = msg[off:n]
					rest = header.LEN - (n - off)
					break
				}
			} else {
				if n-off > rest {
					//粘包
					body += string(buffer[off : off+rest]) //已经是完整的包
					//p.Parse(body)
					pkg := handlerPkg{
						Header: header,
						Body:   body,
					}
					c.handlerBuf <- pkg
					off += rest + 1
					rest = 0 //进入解析头
				} else {
					body += string(buffer[off:n])
					rest -= n - off
					break
				}
			}
		}

	}
}

func (c *Controller) Close() error {

	closed := make(chan struct{})
	c.closing <- closed

	<-closed
	return c.w.Close()
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

func (c *Controller) lastResult() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_last_result := GeneratePackage(MID_0064_OLD_TIGHTING, "006", fmt.Sprintf("%010d", 0), DEFAULT_MSG_END)

	c.Write([]byte(s_last_result))

	return nil
}

func (c *Controller) pset(pset int) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_pset := GeneratePackage(MID_0018_PSET, "001", fmt.Sprintf("%03d", pset), DEFAULT_MSG_END)

	c.Write([]byte(s_pset))

	return nil
}

// 0: set 1: reset
func (c *Controller) JobOff(off string) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_off := GeneratePackage(MID_0130_JOB_OFF, "001", off, DEFAULT_MSG_END)

	c.Write([]byte(s_off))

	return nil
}

func (c *Controller) jobSelect(job int) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	s_job := GeneratePackage(MID_0038_JOB_SELECT, "002", fmt.Sprintf("%04d", job), DEFAULT_MSG_END)

	c.Write([]byte(s_job))

	return nil
}

func (c *Controller) IdentifierSet(str string) error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	ide := GeneratePackage(MID_0150_IDENTIFIER_SET, "001", str, DEFAULT_MSG_END)

	//c.Response.Add(MID_0018_PSET, "")

	c.Write([]byte(ide))

	return nil
}

func (c *Controller) DataSubscribeCurve() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	cs := GeneratePackage(MID_0008_DATA_SUB, "001", "0900001350                             01001", DEFAULT_MSG_END)

	c.Write([]byte(cs))

	return nil
}

func (c *Controller) PSetSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0014_PSET_SUBSCRIBE, "000", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) SelectorSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0250_SELECTOR_SUBSCRIBE, "001", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) JobInfoSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0034_JOB_INFO_SUBSCRIBE, "003", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) ResultSubcribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_0060_LAST_RESULT_SUBSCRIBE, "998", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) CurveSubscribe() error {
	if c.Status() == controller.STATUS_OFFLINE {
		return errors.New("status offline")
	}

	pset := GeneratePackage(MID_7408_LAST_CURVE_SUBSCRIBE, "000", "", DEFAULT_MSG_END)

	c.Write([]byte(pset))

	return nil
}

func (c *Controller) PSet(pset int, workorder_id int64, result_id int64, count int, user_id int64, channel int) (uint32, error) {
	// 设定结果标识
	// 结果id-拧接次数-用户id
	err := c.IdentifierSet(fmt.Sprintf("%d-%d-%d", user_id, count, result_id))
	if err != nil {
		return 0, err
	}

	// 设定pset
	err = c.pset(pset)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (c *Controller) JobSet(result_ids []int64, user_id int64, job int) error {
	ids := ""
	for _, v := range result_ids {
		ids += fmt.Sprintf("%d,", v)
	}
	err := c.IdentifierSet(fmt.Sprintf("%d-%s", user_id, ids))
	if err != nil {
		return err
	}

	err = c.jobSelect(job)
	if err != nil {
		return err
	}

	return nil
}
