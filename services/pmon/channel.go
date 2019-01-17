package pmon

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"sync"
	"sync/atomic"
		"unicode/utf8"
)

type PmonChannelStatus string

const (
	STATUSNORMAL PmonChannelStatus = "normal"
	STATUSCLOSE                    = "close"
)

type Channel struct {
	cChannel
	conn           *Connection
	mux            *sync.Mutex
	BlockCount     int
	recvBuf        chan PmonPackage
	status         PmonChannelStatus
	closed         chan struct{}
	hasSendSD      atomic.Value
	e              PMONEventHandler
	Service        *Service
	RestartPoint   string
	MtxRestarPoint sync.Mutex
	diag           Diagnostic
	ud             interface{} //user data,回调时注入
}

func NewChannel(c cChannel, d Diagnostic) *Channel {
	cc := &Channel{
		cChannel:       c,
		mux:            new(sync.Mutex),
		recvBuf:        make(chan PmonPackage, 10),
		closed:         make(chan struct{}, 2),
		status:         STATUSCLOSE,
		e:              nil,
		ud:             nil,
		BlockCount:     1,
		diag:           d,
		MtxRestarPoint: sync.Mutex{},
	}

	cc.SethasSD(false)

	return cc
}

func (c *Channel) RefreshRestartPoint(restartPoint string) {
	c.MtxRestarPoint.Lock()
	defer c.MtxRestarPoint.Unlock()

	if c.RestartPointLength > 0 {
		c.RestartPoint = restartPoint
	}
}

func (c *Channel) GetRestartPoint() string {
	c.MtxRestarPoint.Lock()
	defer c.MtxRestarPoint.Unlock()

	return c.RestartPoint
}

func (c *Channel) ResetBlockCount() {
	defer c.mux.Unlock()
	c.mux.Lock()

	c.BlockCount = 1
}

func (c *Channel) GetBlockCount() int {
	defer c.mux.Unlock()
	c.mux.Lock()
	x := c.BlockCount
	if x >= 9999 {
		c.BlockCount = 1
	} else {
		c.BlockCount += 1
	}
	return x
}

func (ch *Channel) Start() error {
	err := ch.conn.Open() //尝试打开连接
	if err != nil {
		return errors.Wrapf(err, "Open connection fail,by channel %s", ch.Ch)
	}

	go ch.manage()
	return nil
}

func (ch *Channel) Stop() error {
	err := ch.conn.Close() //尝试关闭连接
	if err != nil {
		return errors.Wrapf(err, "Open connection fail,by channel %s", ch.Ch)
	}
	go ch.manage()
	return nil
}

func (ch *Channel) GetStatus() PmonChannelStatus {
	defer ch.mux.Unlock()
	ch.mux.Lock()
	x := ch.status
	return x
}

func (ch *Channel) RegistryHandler(e PMONEventHandler, ud interface{}) {
	ch.e = e
	ch.ud = ud
}

func (ch *Channel) SetStatus(s PmonChannelStatus) {
	defer ch.mux.Unlock()
	ch.mux.Lock()
	ch.status = s
}

func (ch *Channel) SethasSD(t bool) {
	ch.hasSendSD.Store(t)
}

func (ch *Channel) gethasSD() bool {
	return ch.hasSendSD.Load().(bool)
}

func (ch *Channel) Write(buf []byte, msgType PMONSMGTYPE) error {
	//if msgType == PMONMSGSD && ch.GetStatus() == STATUSCLOSE {
	//	i := 0
	//	msg, err := ch.PMONGenerateMsg(PMONMSGSO, "")
	//	if err != nil {
	//		ch.Service.diag.Error(fmt.Sprintf("Generation %s msg fail", msgType), err)
	//		return errors.Wrap(err, "Channel.Write")
	//	}
	//
	//	ch.conn.Write([]byte(msg[0]), ch.WriteTimeout) //发送此SO忽略错误
	//	ch.diag.Debug(fmt.Sprintf("send msg:%s", msg[0]))
	//	time.Sleep(150 * time.Millisecond)             //sleep 150 ms
	//	for ch.GetStatus() == STATUSCLOSE && i < 6 {
	//		time.Sleep(100 * time.Millisecond) //sleep 100 ms
	//		i++
	//	}
	//	if ch.GetStatus() == STATUSCLOSE {
	//		return fmt.Errorf("channel %s is closed can not send SD and send SO %d times", ch.Ch, i)
	//	}
	//}
	ch.diag.Debug(fmt.Sprintf("send msg:%s", string(buf)))
	err := ch.conn.Write(buf, ch.WriteTimeout)
	if err != nil && msgType == PMONMSGSD {
		ch.SethasSD(true)
	}
	return err
}

func (ch *Channel) SetConnection(c *Connection) {
	ch.conn = c
}

func (ch *Channel) manage() {
	//buf := make([]byte, 64 * 1024)
	//off := 0
	//block_idx := 0
	for {
		select {
		case <-ch.closed:
			if ch.gethasSD() {
				//有发送过sd,但没有收到ad，但是关闭了通道
				res, _ := ch.generateAO(ch.conn.U.GetMsgNum())
				ch.Write([]byte(res), PMONMSGAO) //发送AO,忽略错误
				bc := fmt.Sprintf("%04d", ch.GetBlockCount())
				res, _ = ch.generateAD(ch.conn.U.GetMsgNum(), bc)
				ch.Write([]byte(res), PMONMSGAD) //发送AD忽略错误
				res, _ = ch.generateAC(ch.conn.U.GetMsgNum())
				ch.Write([]byte(res), PMONMSGAC) //发送AC忽略错误
				ch.SethasSD(false)               //走完此流程后没有SD
			}
		case data := <-ch.recvBuf:
			switch data.t {
			case PMONMSGSO:
				var res string
				if ch.Type == PMONMASTER {
					//如果是master 需要回一个SO
					res, _ = ch.generateSO(ch.conn.U.GetMsgNum())
					ch.Write([]byte(res), PMONMSGSO)
				} else {
					res, _ = ch.generateAO(ch.conn.U.GetMsgNum())
					ch.ResetBlockCount()
					ch.SetStatus(STATUSNORMAL)
					ch.Write([]byte(res), PMONMSGAO)
				}
			case PMONMSGAO:
				ch.ResetBlockCount()
				ch.SetStatus(STATUSNORMAL)
			case PMONMSGSC:
				var res string
				if ch.Type == PMONMASTER {
					//如果是master 需要回一个SC
					res, _ = ch.generateSC(ch.conn.U.GetMsgNum())
					ch.Write([]byte(res), PMONMSGSC)
				} else {
					res, _ = ch.generateAC(ch.conn.U.GetMsgNum())
					ch.SetStatus(STATUSCLOSE)
					ch.Write([]byte(res), PMONMSGAC)
				}
			case PMONMSGAC:
				ch.SetStatus(STATUSCLOSE)
				ch.closed <- struct{}{}
			case PMONMSGSD:
				BlockCounter := string(data.data[:4])
				dataLen, _ := strconv.Atoi(string(data.data[4:8]))
				res, _ := ch.generateAD(ch.conn.U.GetMsgNum(), BlockCounter)
				ch.Write([]byte(res), PMONMSGAD) //发送AD
				//if ch.Segment == NOSEGMENT {
				//	buf = data.data[8: 8 +dataLen]
				//}else {
				//	//需要进行偏移赋值
				//	d := data.data[8: 8 +dataLen]
				//	l := len(d)
				//	buf[off: off + l] = d
				//	off += l
				//}
				//block_idx += 1
				//if ch.e != nil {
				//	_off := off
				//	off = 0 //准备接受下一个包 msg
				//	ch.e(1, string(buf[:_off]), ch.ud)
				//}
				d := data.data[8 : 8+dataLen]
				var r []rune
				for i := 0; i < dataLen; {
					s, size := utf8.DecodeRune(d[i:])
					r = append(r, s)
					i += size
				}
				if ch.e != nil {
					ch.e(nil, r, ch.ud)
				}
			case PMONMSGAD:
				ch.SethasSD(false)
				ret := make([]string, 2) //最大程度为2
				ackInfo := string(data.data[:2])
				recipientInfo := string(data.data[2:4])
				if ackInfo != "00" {
					ret[0] = ackInfo
				}
				if recipientInfo != "00" {
					ret[1] = recipientInfo
				}
				if ch.e != nil {
					if len(ret) != 0 {
						ch.e(fmt.Errorf("AD msg return fail %s", ret), nil, ch.ud)
					}
				}
				//block_idx = 0
			}
		}

	}
}
