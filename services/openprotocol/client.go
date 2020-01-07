package openprotocol

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/tightening_device"
	"github.com/masami10/rush/socket_writer"
	"github.com/masami10/rush/utils"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type IClientHandler interface {
	handleMsg(pkg *handlerPkg) error
	handleStatus(sn string, status string)
	GetVendorMid(mid string) (string, error)
}

func createClientContext(endpoint string, diag Diagnostic, handler IClientHandler, sn string) *clientContext {
	ctx := clientContext{
		buffer:          make(chan []byte, 1024),
		closing:         make(chan struct{}, 1),
		diag:            diag,
		clientHandler:   handler,
		sn:              sn,
		tempResultCurve: tightening_device.NewTighteningCurve(),
	}

	ctx.sockClient = socket_writer.NewSocketWriter(endpoint, &ctx)
	return &ctx
}

type clientContext struct {
	sn                string
	status            atomic.Value
	sockClient        *socket_writer.SocketWriter
	keepAliveCount    int32
	keepaliveDeadLine atomic.Value
	buffer            chan []byte
	handlerBuf        chan handlerPkg
	Response          ResponseQueue
	receiveBuf        chan []byte
	tempResultCurve   *tightening_device.TighteningCurve
	requestChannel    chan uint32
	sequence          *utils.Sequence
	handleRecvBuf     []byte
	writeOffset       int
	closing           chan struct{}

	diag          Diagnostic
	clientHandler IClientHandler
}

func (c *clientContext) start() {
	go c.manage()
	go c.handleRecv()
	go c.Connect()
}

func (c *clientContext) stop() error {
	c.closing <- struct{}{}
	return nil
}

func (c *clientContext) handlePackageOPPayload(src []byte, data []byte) error {
	msg := append(src, data...)

	//c.diag.Debug(fmt.Sprintf("%s op target buf: %s", c.deviceConf.SerialNumber, string(msg)))

	lenMsg := len(msg)

	// 如果头的长度不够
	if lenMsg < LEN_HEADER {
		return errors.New("Head Is Error")
	}

	header := OpenProtocolHeader{}
	header.Deserialize(string(msg[0:LEN_HEADER]))

	// 如果body的长度匹配
	if header.LEN == lenMsg-LEN_HEADER {
		pkg := handlerPkg{
			SN:     c.sn,
			Header: header,
			Body:   string(msg[LEN_HEADER : LEN_HEADER+header.LEN]),
		}

		c.handlerBuf <- pkg
	} else {
		return errors.New(fmt.Sprintf("Body Len Err: %s", string(msg)))
	}

	return nil
}

func (c *clientContext) handleRecv() {
	c.handleRecvBuf = make([]byte, 65535)
	c.receiveBuf = make(chan []byte, 65535)
	lenBuf := len(c.handleRecvBuf)

	for {
		select {
		case buf := <-c.receiveBuf:
			// 处理接收缓冲
			var readOffset = 0

			for {
				if readOffset >= len(buf) {
					break
				}
				index := bytes.IndexByte(buf[readOffset:], OP_TERMINAL)
				if index == -1 {
					// 没有结束字符,放入缓冲等待后续处理
					restBuf := buf[readOffset:]
					if c.writeOffset+len(restBuf) > lenBuf {
						c.diag.Error("full", errors.New("full"))
						break
					}

					copy(c.handleRecvBuf[c.writeOffset:c.writeOffset+len(restBuf)], restBuf)
					c.writeOffset += len(restBuf)
					break
				} else {
					// 找到结束字符，结合缓冲进行处理
					err := c.handlePackageOPPayload(c.handleRecvBuf[0:c.writeOffset], buf[readOffset:readOffset+index])
					if err != nil {
						//数据需要丢弃
						c.diag.Error("msg", err)
					}

					c.writeOffset = 0
					readOffset += index + 1
				}
			}
		}
	}
}

func (c *clientContext) manage() {

	nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(time.Duration(OpenProtocolDefaultKeepAlivePeriod)):
			if c.Status() == device.BaseDeviceStatusOffline {
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

			err := c.sockClient.Write(v)
			if err != nil {
				c.diag.Error("Write Data Fail", err)
			} else {
				c.updateKeepAliveDeadLine()
			}
			nextWriteThreshold = time.Now().Add(300 * time.Millisecond)

		case <-c.closing:
			c.diag.Debug("manage exit")
			return //退出manage协程

		case pkg := <-c.handlerBuf:
			err := c.clientHandler.handleMsg(&pkg)
			if err != nil {
				c.diag.Error("Open IProtocol handleMsg Fail", err)
			}
		}
	}
}

func (c *clientContext) Read(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			c.diag.Error("Controller Close Error", err)
		}
	}()

	buffer := make([]byte, 65535)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(OpenProtocolDefaultKeepAlivePeriod) * MAX_KEEP_ALIVE_CHECK).Add(1 * time.Second))
		n, err := conn.Read(buffer)
		if err != nil {
			c.diag.Error("read failed", err)
			c.handleStatus(device.BaseDeviceStatusOffline)
			c.clientHandler.handleStatus(c.sn, device.BaseDeviceStatusOffline)
			break
		}

		c.updateKeepAliveCount(0)
		c.receiveBuf <- buffer[0:n]
	}
}

func (c *clientContext) Status() string {
	return c.status.Load().(string)
}

func (c *clientContext) UpdateStatus(status string) {
	c.status.Store(status)
}

func (c *clientContext) handleStatus(status string) {

	if status != c.Status() {

		c.UpdateStatus(status)

		if status == device.BaseDeviceStatusOffline {

			// 断线重连
			go c.Connect()
		}
	}
}

func (c *clientContext) KeepAliveCount() int32 {
	return atomic.LoadInt32(&c.keepAliveCount)
}

func (c *clientContext) updateKeepAliveCount(i int32) {
	atomic.SwapInt32(&c.keepAliveCount, i)
}

func (c *clientContext) addKeepAliveCount() {
	atomic.AddInt32(&c.keepAliveCount, 1)
}

func (c *clientContext) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(time.Duration(OpenProtocolDefaultKeepAlivePeriod)))
}

func (c *clientContext) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *clientContext) sendKeepalive() {
	if c.Status() == device.BaseDeviceStatusOffline {
		return
	}

	keepAlive := GeneratePackage(MID_9999_ALIVE, DEFAULT_REV, "1", "", "", "")
	c.Write([]byte(keepAlive))
}

func (c *clientContext) Write(buf []byte) {
	c.diag.Debug(fmt.Sprintf("OpenProtocol Send %s: %s", c.sn, string(buf)))
	c.buffer <- buf
}

func (c *clientContext) Connect() {
	c.UpdateStatus(device.BaseDeviceStatusOffline)
	c.handlerBuf = make(chan handlerPkg, 1024)
	c.writeOffset = 0
	c.requestChannel = make(chan uint32, 1024)
	c.sequence = utils.CreateSequence()
	c.Response = ResponseQueue{
		Results: map[interface{}]interface{}{},
		mtx:     sync.Mutex{},
	}

	for {
		err := c.sockClient.Connect(DAIL_TIMEOUT)
		if err != nil {
			c.diag.Error("connect", err)
		} else {
			break
		}

		time.Sleep(time.Duration(OpenProtocolDefaultRequestTimeOut))
	}

	c.handleStatus(device.BaseDeviceStatusOnline)
	c.clientHandler.handleStatus(c.sn, device.BaseDeviceStatusOnline)
}

func (c *clientContext) ProcessRequest(mid string, noack string, station string, spindle string, data string) (interface{}, error) {
	rev, err := c.clientHandler.GetVendorMid(mid)
	if err != nil {
		return nil, err
	}

	if c.Status() == device.BaseDeviceStatusOffline {
		return nil, errors.New(device.BaseDeviceStatusOffline)
	}

	pkg := GeneratePackage(mid, rev, noack, station, spindle, data)

	seq := c.sequence.GetSequence()
	c.requestChannel <- seq
	c.Response.Add(seq, nil)

	c.Write([]byte(pkg))
	ctx, _ := context.WithTimeout(context.Background(), MAX_REPLY_TIME)
	reply := c.Response.Get(seq, ctx)

	if reply == nil {
		return nil, errors.New(tightening_device.TIGHTENING_ERR_TIMEOUT)
	}

	return reply, nil
}
