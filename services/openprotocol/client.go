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
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

const (
	DailTimeout = 5 * time.Second
	BufferSize  = 65535
)

type IClientHandler interface {
	handleMsg(pkg *handlerPkg) error
	HandleStatus(sn string, status string)
	GetVendorMid(mid string) (string, error)
}

func newClientContext(endpoint string, diag Diagnostic, handler IClientHandler, sn string, params *OpenProtocolParams) *clientContext {
	ctx := clientContext{
		sendBuffer:      make(chan []byte, BufferSize),
		closing:         make(chan struct{}, 1),
		diag:            diag,
		clientHandler:   handler,
		sn:              sn,
		tempResultCurve: tightening_device.NewTighteningCurve(),
		params:          params,
	}

	ctx.sockClient = socket_writer.NewSocketWriter(endpoint, &ctx)
	return &ctx
}

type clientContext struct {
	sn                string
	params            *OpenProtocolParams
	status            atomic.Value
	sockClient        *socket_writer.SocketWriter
	keepAliveCount    atomic.Int32
	keepaliveDeadLine atomic.Value
	sendBuffer        chan []byte
	handlerBuf        chan handlerPkg
	response          ResponseQueue
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
	go c.handleRecv()
	go c.manage()
	go c.connect()
}

func (c *clientContext) stop() {
	c.closing <- struct{}{}
}

func (c *clientContext) handlePackageOPPayload(src []byte, data []byte) error {
	msg := append(src, data...)

	lenMsg := len(msg)

	// 如果头的长度不够
	if lenMsg < LenHeader {
		return errors.New("Head Is Error ")
	}

	header := OpenProtocolHeader{}
	header.Deserialize(string(msg[0:LenHeader]))

	// 如果body的长度匹配
	if header.LEN == lenMsg-LenHeader {
		pkg := handlerPkg{
			SN:     c.sn,
			Header: header,
			Body:   string(msg[LenHeader : LenHeader+header.LEN]),
		}

		c.handlerBuf <- pkg
	} else {
		return errors.New(fmt.Sprintf("Body Len Err: %s", string(msg)))
	}

	return nil
}

func (c *clientContext) handleRecv() {
	c.handleRecvBuf = make([]byte, BufferSize)
	c.receiveBuf = make(chan []byte, BufferSize)
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
				index := bytes.IndexByte(buf[readOffset:], OpTerminal)
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
		case <-time.After(c.params.KeepAlivePeriod):
			if c.Status() == device.BaseDeviceStatusOffline {
				continue
			}

			if c.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.sendKeepalive()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.addKeepAliveCount()
			}
		case v := <-c.sendBuffer:
			for nextWriteThreshold.After(time.Now()) {
				time.Sleep(time.Microsecond * 100)
			}

			err := c.sockClient.Write(v)
			if err != nil {
				c.diag.Error("IOWrite Data Fail", err)
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
			c.diag.Error("Controller Close Error ", err)
		}
	}()

	buf := make([]byte, BufferSize)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(c.params.KeepAlivePeriod * time.Duration(c.params.MaxKeepAliveCheck))); err != nil {
			c.diag.Error("SetReadDeadline Failed ", err)
			break
		}

		n, err := conn.Read(buf)
		if err != nil {
			c.diag.Error("IORead Failed ", err)
			c.handleStatus(device.BaseDeviceStatusOffline)
			c.clientHandler.HandleStatus(c.sn, device.BaseDeviceStatusOffline)
			break
		}

		c.updateKeepAliveCount(0)
		c.receiveBuf <- buf[0:n]
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
			go c.connect()
		}
	}
}

func (c *clientContext) KeepAliveCount() int32 {
	return c.keepAliveCount.Load()
}

func (c *clientContext) updateKeepAliveCount(i int32) {
	c.keepAliveCount.Swap(i)
}

func (c *clientContext) addKeepAliveCount() {
	c.keepAliveCount.Inc()
}

func (c *clientContext) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(time.Duration(c.params.KeepAlivePeriod)))
}

func (c *clientContext) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *clientContext) sendKeepalive() {
	if c.Status() == device.BaseDeviceStatusOffline {
		return
	}

	keepAlive := GeneratePackage(MID_9999_ALIVE, DefaultRev, "1", "", "", "")
	c.Write([]byte(keepAlive))
}

func (c *clientContext) Write(buf []byte) {
	c.diag.Debug(fmt.Sprintf("OpenProtocol Send %s: %s", c.sn, string(buf)))
	c.sendBuffer <- buf
}

func (c *clientContext) connect() {
	c.UpdateStatus(device.BaseDeviceStatusOffline)
	c.handlerBuf = make(chan handlerPkg, 1024)
	c.writeOffset = 0
	c.requestChannel = make(chan uint32, 1024)
	c.sequence = utils.CreateSequence()
	c.response = ResponseQueue{
		Results: map[interface{}]interface{}{},
		mtx:     sync.Mutex{},
	}

	for {
		err := c.sockClient.Connect(DailTimeout)
		if err != nil {
			c.diag.Error("connect", err)
		} else {
			break
		}

		time.Sleep(1 * time.Second)
	}

	c.handleStatus(device.BaseDeviceStatusOnline)
	c.clientHandler.HandleStatus(c.sn, device.BaseDeviceStatusOnline)
	if err := c.startComm(); err != nil {
		c.diag.Error(fmt.Sprintf("Start Comm Failed: %s", c.sn), err)
	}
}

func (c *clientContext) startComm() error {
	reply, err := c.ProcessRequest(MID_0001_START, "", "", "", "")
	if err != nil {
		return err
	}

	if reply.(string) != request_errors["00"] && reply.(string) != request_errors["96"] {
		return errors.New(reply.(string))
	}

	return nil
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
	c.response.Add(seq, nil)

	c.Write([]byte(pkg))
	ctx, _ := context.WithTimeout(context.Background(), c.params.MaxReplyTime)
	reply := c.response.Get(seq, ctx)

	if reply == nil {
		return nil, errors.New(tightening_device.TIGHTENING_ERR_TIMEOUT)
	}

	return reply, nil
}
