package pmon

import (
	"fmt"
	"github.com/masami10/aiis/services/pmon/udp_driver"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"
)

type DispatchPkg struct {
	p  PmonPackage
	ch string
}

type Dispatcher interface {
	Dispatch(PmonPackage, string) //数据字节流，通道名字
}

//为了定位某个通道
type channelInfo struct {
	name string //通道名字
	sNoT string
	sNoR string
}

type Connection struct {
	U           *udp_driver.UDPDriver
	started     bool
	workers     int
	wg          sync.WaitGroup
	name        string
	channels    []channelInfo
	closing     chan struct{}
	DispatchBuf chan DispatchPkg
	Dispatcher
	service *Service
}

func NewConnection(addr string, name string, deadline time.Duration, workers int) *Connection {
	u := udp_driver.NewUDPDriver(addr, deadline)
	c := &Connection{
		U:           u,
		started:     false,
		name:        name,
		workers:     workers,
		DispatchBuf: make(chan DispatchPkg, workers),
		closing:     make(chan struct{}),
	}
	u.SetConnection(c) //注入服务为了进行分发

	return c
}

func (c *Connection) DispatchProcess() {
	for {
		select {
		case pkg := <-c.DispatchBuf:
			c.Dispatcher.Dispatch(pkg.p, pkg.ch)

		case <-c.closing:
			c.wg.Done()
			return
		}
	}
}

//连接中打开相关的通道
func (c *Connection) Open() error {
	if !c.started {
		err := c.U.Open()
		if err != nil {
			return errors.Wrap(err, "Open Connection fail")
		}

		for i := 0; i < c.workers; i++ {
			c.wg.Add(1)

			go c.DispatchProcess()
		}

		c.started = true

	}
	return nil
}

func (c *Connection) Close() error {
	if c.started {
		err := c.U.Close()
		if err != nil {
			return err
		}
		for i := 0; i < c.workers; i++ {
			c.closing <- struct{}{}
		}

		c.wg.Wait()

		c.started = false
	}
	return nil
}

func (c *Connection) SetDispatcher(d Dispatcher) error {
	c.Dispatcher = d
	return nil
}

func (c *Connection) SetService(srv *Service) {
	c.service = srv
}

func (c *Connection) AppendChannel(name string, sNoT string, sNoR string) {
	ch := channelInfo{
		name: name,
		sNoT: sNoT,
		sNoR: sNoR,
	}
	c.channels = append(c.channels, ch)
}

func (c *Connection) Write(buf []byte, deadline time.Duration) error {
	err := c.U.Write(buf, deadline)
	return err
}

func (c *Connection) dispatch(buf []byte) error {
	//将接受的数据进行通道分发
	rNoT := string(buf[5:9])
	rNoR := string(buf[9:13])
	for _, ch := range c.channels {
		if ch.sNoR == rNoT && ch.sNoT == rNoR {

			c.service.diag.Debug(fmt.Sprintf("channel:%s, recv msg: %s", ch.name, string(buf)))

			p := PMONParseMsg(buf)

			pkg := DispatchPkg{
				p:  p,
				ch: ch.name,
			}

			c.DispatchBuf <- pkg
		}
	}
	return nil
}

func (c *Connection) Parse(buf []byte) error {
	err := ValidateChecksum(buf)
	if err != nil {
		log.Printf("Validate fail %s ", err)
	}
	if IsUDPResponse(buf) {
		c.U.NoReadDeadline() //设定为read永远阻塞
		return nil
	}
	msgId := GetMsgId(buf)
	c.Write([]byte(UdpResponse(msgId)), time.Duration(time.Millisecond*5)) //收到消息先发送udpResponse
	c.dispatch(buf)                                                        //连接将数据分发到不同的通道
	return nil
}
