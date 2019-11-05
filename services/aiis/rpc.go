package aiis

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
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

	go c.connect()
	go c.manage()

	return nil
}

func (c *GRPCClient) Stop() error {
	if c.conn != nil {
		c.conn.Close()
	}
	if c.stream != nil {
		c.stream.CloseSend()
	}

	closed := make(chan struct{})
	c.closing <- closed
	<-closed

	return nil
}

func (c *GRPCClient) Status() string {

	return c.status.Load().(string)
}

func (c *GRPCClient) setStatus(status string) {

	c.status.Store(status)
}

func (c *GRPCClient) updateStatus() {

	prevStatus := c.Status()
	conn := c.conn
	if conn == nil {
		return
	}
	status := RPC_OFFLINE
	ss := conn.GetState()
	switch ss {
	case connectivity.Connecting:
		c.srv.diag.Debug("GRPC Connection Is Connecting")
	case connectivity.Idle:
	case connectivity.Ready:
		c.srv.diag.Debug(fmt.Sprintf("GRPC Connection Is %s", ss))
		status = RPC_ONLINE
	case connectivity.Shutdown:
	case connectivity.TransientFailure:
	default:
		c.srv.diag.Debug(fmt.Sprintf("GRPC Connection Is %s, Try To Connecting", ss))
		status = RPC_OFFLINE
	}

	if prevStatus != status {
		c.setStatus(status)
		if c.OnRPCStatus != nil {
			c.OnRPCStatus(status)
		}
	}
}

func (c *GRPCClient) manage() {

	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			c.updateStatus()

		case stopDone := <-c.closing:
			ticker.Stop()
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *GRPCClient) connect() {
	var err error = nil
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithBackoffConfig(grpc.BackoffConfig{
		MaxDelay: 1 * time.Second,
	}))
	kp := keepalive.ClientParameters{
		Time:    PING_ITV,
		Timeout: KEEP_ALIVE_CHECK * PING_ITV,
	}
	opts = append(opts, grpc.WithKeepaliveParams(kp))
	for ; ; {
		c.conn, err = grpc.Dial(c.srv.Config().GRPCServer, opts...)
		if err != nil {
			c.srv.diag.Error("GRPC connect Error", err)
		} else {
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
			c.srv.diag.Error("GRPC RecvProcess Error", err)
			return
		}

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
