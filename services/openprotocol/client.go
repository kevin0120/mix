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
	MaxProc     = 3
)

type IClientHandler interface {
	handleMsg(pkg *handlerPkg) error
	HandleStatus(sn string, status string)
	GetVendorMid(mid string) (string, error)
	UpdateToolStatus(sn string, status string)
}

func newClientContext(endpoint string, diag Diagnostic, handler IClientHandler, sn string, params *OpenProtocolParams) *clientContext {
	ctx := clientContext{
		closing:         make(chan struct{}, MaxProc),
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

	closing           chan struct{}
	closinghandleRecv chan struct{}

	diag          Diagnostic
	clientHandler IClientHandler
}

func (c *clientContext) start() {
	c.initProcs()
	go c.connect()
}

func (c *clientContext) stop() {
	for i := 0; i < MaxProc; i++ {
		c.closing <- struct{}{}
	}
}

func (c *clientContext) initProcs() {
	go c.procWrite()
	go c.procHandle()
	go c.procAlive()
}

func (c *clientContext) handlePackageOPPayload(src []byte) error {
	msg := src
	lenMsg := len(msg)

	// 如果头的长度不够
	if lenMsg < LenHeader {
		return errors.New(fmt.Sprintf("OpenProtocl SN:%s Head Is Error: %s", c.sn, string(msg)))
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
		return errors.New(fmt.Sprintf("OpenProtocl SN:%s Body Len Err: %s", c.sn, string(msg)))
	}

	return nil
}

func (c *clientContext) procHandleRecv() {
	c.receiveBuf = make(chan []byte, BufferSize)
	handleRecvBuf := make([]byte, BufferSize)
	lenBuf := len(handleRecvBuf)
	var writeOffset = 0

	for {
		select {
		case buf := <-c.receiveBuf:
			// 处理接收缓冲
			var readOffset = 0
			var index = 0

			for {
				if readOffset >= len(buf) {
					break
				}
				index = bytes.IndexByte(buf[readOffset:], OpTerminal)
				if index == -1 {
					// 没有结束字符,放入缓冲等待后续处理
					restBuf := buf[readOffset:]
					if writeOffset+len(restBuf) > lenBuf {
						c.diag.Error("full", errors.New("full"))
						break
					}

					copy(handleRecvBuf[writeOffset:writeOffset+len(restBuf)], restBuf)
					writeOffset += len(restBuf)
					break
				} else {
					// 找到结束字符，结合缓冲进行处理
					targetBuf := append(handleRecvBuf[0:writeOffset], buf[readOffset:readOffset+index]...)
					if len(buf) == 1 && index == 0 {
						targetBuf = handleRecvBuf[0:writeOffset]
					}

					err := c.handlePackageOPPayload(targetBuf)
					if err != nil {
						//数据需要丢弃
						c.diag.Error("handlePackageOPPayload Error", err)
						c.diag.Debug(fmt.Sprintf("procHandleRecv Raw Msg:%s", string(buf)))
						c.diag.Debug(fmt.Sprintf("procHandleRecv Rest Msg:%s writeOffset:%d readOffset:%d index:%d", string(handleRecvBuf), writeOffset, readOffset, index))
					}

					writeOffset = 0
					readOffset += index + 1
				}
			}

		case <-c.closinghandleRecv:
			c.diag.Debug("procHandleRecv Exit")
			return
		}

	}
}

func (c *clientContext) procWrite() {
	c.sendBuffer = make(chan []byte, BufferSize)

	for {
		select {
		case v := <-c.sendBuffer:
			err := c.sockClient.Write(v)
			if err != nil {
				c.diag.Error("IOWrite Data Fail", err)
			} else {
				c.diag.Debug(fmt.Sprintf("OpenProtocol Send %s: %s", c.sn, string(v)))
				c.updateKeepAliveDeadLine()
			}

			time.Sleep(100 * time.Millisecond)

		case <-c.closing:
			c.diag.Debug("procWrite Exit")
			return
		}
	}
}

func (c *clientContext) procHandle() {
	c.handlerBuf = make(chan handlerPkg, BufferSize)

	for {
		select {
		case pkg := <-c.handlerBuf:
			err := c.clientHandler.handleMsg(&pkg)
			if err != nil {
				c.diag.Error("Open IProtocol handleMsg Fail", err)
			}

		case <-c.closing:
			c.diag.Debug("procHandle Exit")
			return
		}
	}
}

func (c *clientContext) procAlive() {

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

		case <-c.closing:
			c.diag.Debug("procAlive Exit")
			return

		}
	}
}

func (c *clientContext) Read(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			c.diag.Error("Controller Close Error ", err)
		}

		c.closinghandleRecv <- struct{}{}
	}()

	buf := make([]byte, BufferSize)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(c.params.KeepAlivePeriod * time.Duration(c.params.MaxKeepAliveCheck)).Add(1 * time.Second)); err != nil {
			c.diag.Error("SetReadDeadline Failed ", err)
			break
		}

		n, err := conn.Read(buf)
		if err != nil {
			c.diag.Error("Failed ", err)
			c.handleStatus(device.BaseDeviceStatusOffline)
			c.clientHandler.HandleStatus(c.sn, device.BaseDeviceStatusOffline)
			c.clientHandler.UpdateToolStatus(c.sn, device.BaseDeviceStatusOffline)
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
	c.keepaliveDeadLine.Store(time.Now().Add(c.params.KeepAlivePeriod))
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
	c.sendBuffer <- buf
}

func (c *clientContext) resetConn() {
	c.UpdateStatus(device.BaseDeviceStatusOffline)

	c.requestChannel = make(chan uint32, 1024)
	c.sequence = utils.CreateSequence()
	c.response = ResponseQueue{
		Results: map[interface{}]interface{}{},
		mtx:     sync.Mutex{},
	}
}

func (c *clientContext) connect() {
	c.diag.Debug(fmt.Sprintf("OpenProtocol SN:%s Connecting ...", c.sn))
	c.updateKeepAliveDeadLine()
	c.resetConn()

	for {
		err := c.sockClient.Connect(DailTimeout)
		if err != nil {
			c.diag.Error("connect", err)
		} else {
			c.diag.Debug(fmt.Sprintf("OpenProtocol SN:%s Connected", c.sn))
			break
		}

		time.Sleep(1 * time.Second)
	}

	c.handleStatus(device.BaseDeviceStatusOnline)
	c.clientHandler.HandleStatus(c.sn, device.BaseDeviceStatusOnline)
	c.clientHandler.UpdateToolStatus(c.sn, device.BaseDeviceStatusOnline)

	c.closinghandleRecv = make(chan struct{}, 1)
	go c.procHandleRecv()

	time.Sleep(100 * time.Millisecond)
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
