package udp_driver

import (
	"fmt"
	"strings"
	"time"
	"sync"
	"log"
	"net"
	"github.com/pkg/errors"
)


type Connection interface {
	Parse(buf []byte) error
}

type UDPDriver struct {
	Writer     *SocketWriter
	Listener   *SocketListener
	ReadTimeout time.Duration
	ReadDeadline time.Time
	Connection
	mux 		*sync.Mutex
	messageNum int //driver层存在的msg num，为了进行组包时候需要存在
	checkTimeout       chan struct{} //检查readtimeout
}

// addr 格式为 udp://REMOTE_IP:PORT
func NewUDPDriver(addr string, deadline time.Duration) *UDPDriver {
	sl := strings.Split(addr, ":")
	listenerAddr := fmt.Sprintf("udp://:%s", sl[2])
	l := NewSocketListener(listenerAddr)

	w := NewSocketWriter(addr, time.Duration(time.Duration(0)))

	return &UDPDriver{
		Writer:     w,
		Listener:   l,
		messageNum: 1, //start from 1
		mux: new(sync.Mutex),
		ReadTimeout: deadline,
		checkTimeout: make(chan struct{}),
	}
}

func (u *UDPDriver) GetMsgNum() int{
	defer u.mux.Unlock()
	u.mux.Lock()
	x := u.messageNum
	if x >= 9999 {
		u.messageNum = 1
	}else{
		u.messageNum += 1
	}
	return x
}

func (u *UDPDriver) Open() error{
	if u.Listener != nil {
		u.Listener.SetProtocol(u) //设定协议，必须在打开之前
		if err := u.Listener.Start(); err != nil {
			//将会打开接收协程,
			return err
		}
	}
	if u.Writer != nil {
		if err := u.Writer.Connect(); err != nil {
			return err
		}
	}

	//go u.manage()
	return nil
}

func (u *UDPDriver) Close() error  {
	if u.Listener != nil {
		err := u.Listener.Close()
		if err != nil {
			return errors.Wrap(err, "Close the UDP Driver Listener")
		}
		u.Listener = nil
	}
	if u.Writer != nil {
		if err := u.Writer.Close(); err != nil {
			return errors.Wrap(err, "Close the UDP Driver Writer")
		}
		u.Writer = nil
	}
	return nil
}

//func (u *UDPDriver) manage(){
//	for {
//		<- u.checkTimeout
//		time.Sleep(u.ReadTimeout)
//	}
//}

func (u *UDPDriver) xWrite(buf []byte , deadline time.Duration) error{
	u.Writer.SetWriteDeadline(time.Now().Add(deadline)) //write 设定tiemout时间
	return  u.Writer.Write(buf)
}

func (u *UDPDriver) Write(buf []byte , deadline time.Duration) error {
	mr := 3 //最大重试次数
	var err error = nil
	for i := 0; i < mr; i ++{
		err = u.xWrite(buf, deadline)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("D! Timeout in Write: %s", err)
				time.Sleep(2 * deadline) // 2倍timeout后重试
				continue
			} else if netErr != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
				log.Printf(" %s", err)
				return err
			}
		}
		//没有错误直接跳出循环
		break
	}
	if err != nil {
		//最后重试也没成功
		return err
	}
	u.updateReadDeadline() //写成功，更新read deadline
	return nil
}

func (u *UDPDriver) updateReadDeadline() error{
	return u.Listener.InterListener.setReadDeadLine(time.Now().Add( u.ReadTimeout )) //更新read deadline
}

func (u *UDPDriver) NoReadDeadline() error{
	return u.Listener.InterListener.setReadDeadLine(time.Time{}) //设定为永远block
}

func (u *UDPDriver) SetConnection( c Connection) {
	u.Connection = c
}

func (u *UDPDriver) Parse(buf []byte) error {
	err := u.Connection.Parse(buf)
	if err != nil {
		log.Printf("Validate fail %s ",err)
	}

	return nil
}

func (u *UDPDriver) Read(c net.Conn) {
	//无需实现
}
func (u *UDPDriver) NewConn(c net.Conn) {
	//无需事先
}
