package aiis

import (
	"github.com/kataras/iris/core/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

const (
	TYPE_RESULT      = "result_patch"
	TYPE_ODOO_STATUS = "odoo_status"
	TYPE_WORKORDER   = "workorder"

	RPC_PING = "ping"
	RPC_PONG = "pong"

	PING_ITV         = 1 * time.Second
	KEEP_ALIVE_CHECK = 3

	RPC_OFFLINE = "offline"
	RPC_ONLINE  = "online"
)

type RPCPayload struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type OnRPCRecv func(payload string)
type OnRPCStatus func(status string)

type GRPCClient struct {
	srv               *Service
	conn              *grpc.ClientConn
	stream            RPCAiis_RPCNodeClient
	opts              []grpc.DialOption
	rpcClient         RPCAiisClient
	RPCRecv           OnRPCRecv
	OnRPCStatus       OnRPCStatus
	status            atomic.Value
	keepAliveCount    int32
	keepaliveDeadLine atomic.Value
	closing           chan chan struct{}
}

func (c *GRPCClient) Start() error {
	c.status.Store(RPC_OFFLINE)
	c.closing = make(chan chan struct{})
	c.updateKeepAliveDeadLine()

	go c.Connect()
	go c.manage()

	return nil
}

func (c *GRPCClient) Stop() error {
	c.conn.Close()
	c.stream.CloseSend()

	closed := make(chan struct{})
	c.closing <- closed

	return nil
}

func (c *GRPCClient) Status() string {

	return c.status.Load().(string)
}

func (c *GRPCClient) KeepAliveCount() int32 {
	return atomic.LoadInt32(&c.keepAliveCount)
}

func (c *GRPCClient) updateKeepAliveCount(i int32) {
	atomic.SwapInt32(&c.keepAliveCount, i)
}

func (c *GRPCClient) addKeepAliveCount() {
	atomic.AddInt32(&c.keepAliveCount, 1)
}

func (c *GRPCClient) updateKeepAliveDeadLine() {
	c.keepaliveDeadLine.Store(time.Now().Add(time.Duration(c.srv.Config().KeepAlive)))
}

func (c *GRPCClient) KeepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *GRPCClient) sendPing() {
	c.RPCSend(RPC_PING)
}

func (c *GRPCClient) updateStatus(status string) {

	if status != c.Status() {

		c.status.Store(status)

		if status == RPC_OFFLINE {
			c.srv.diag.Debug("grpc disconnected")

			c.conn.Close()
			c.stream.CloseSend()

			// 断线重连
			go c.Connect()
		}

		if c.OnRPCStatus != nil {
			c.OnRPCStatus(status)
		}

		// 将最新状态推送给hmi
		//s := wsnotify.WSStatus{
		//	SN:     c.cfg.SN,
		//	Status: string(status),
		//}
		//
		//msg, _ := json.Marshal(s)
		//c.Srv.WS.WSSendControllerStatus(string(msg))
	}
}

func (c *GRPCClient) manage() {
	//nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(time.Duration(c.srv.Config().KeepAlive)):
			if c.Status() == RPC_OFFLINE {
				continue
			}
			if c.KeepAliveCount() >= KEEP_ALIVE_CHECK {
				go c.updateStatus(RPC_OFFLINE)
				c.updateKeepAliveCount(0)
				continue
			}
			if c.KeepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.sendPing()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.addKeepAliveCount()
			}

		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *GRPCClient) Connect() {
	var err error = nil

	for {
		c.srv.diag.Debug("grpc connecting ...\n")
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithBackoffConfig(grpc.BackoffConfig{
			MaxDelay: 1 * time.Second,
		}))

		c.conn, err = grpc.Dial(c.srv.Config().GRPCServer, opts...)
		if err == nil {
			c.srv.diag.Debug("grpc connected\n")
			c.updateStatus(RPC_ONLINE)
			break
		}

		time.Sleep(1 * time.Second)
	}

	c.rpcClient = NewRPCAiisClient(c.conn)
	c.stream, _ = c.rpcClient.RPCNode(context.Background())

	go c.RecvProcess()
}

func (c *GRPCClient) RecvProcess() {
	for {
		if c.stream == nil {
			continue
		}

		in, err := c.stream.Recv()
		if err != nil {
			c.srv.diag.Debug("rpc RecvProcess Exit")
			return
		}

		c.updateKeepAliveCount(0)

		if c.RPCRecv != nil {
			c.RPCRecv(in.Payload)
		}
	}
}

func (c *GRPCClient) RPCSend(payload string) error {

	if c.stream != nil {
		return c.stream.Send(&Payload{
			Payload: payload,
		})
	}

	return errors.New("rpc not connected")
}
