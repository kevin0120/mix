package aiis

import (
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

const (
	RPC_PING = "ping"
	RPC_PONG = "pong"

	PING_ITV         = 1 * time.Second
	KEEP_ALIVE_CHECK = 3
)

func NewGRPCClient(d Diagnostic, cfg Config) *GRPCClient {
	return &GRPCClient{
		diag:          d,
		cfg:           cfg,
		BaseTransport: BaseTransport{},
		recvBuf:       make(chan string, 1024),
	}
}

type GRPCClient struct {
	cfg               Config
	conn              *grpc.ClientConn
	stream            RPCAiis_RPCNodeClient
	opts              []grpc.DialOption
	rpcClient         RPCAiisClient
	diag              Diagnostic
	status            atomic.Value
	keepAliveCount    int32
	keepaliveDeadLine atomic.Value
	closing           chan chan struct{}
	recvBuf           chan string

	BaseTransport
}

func (c *GRPCClient) Start() error {
	c.status.Store(TRANSPORT_STATUS_OFFLINE)
	c.closing = make(chan chan struct{})
	c.updateKeepAliveDeadLine()

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
	c.keepaliveDeadLine.Store(time.Now().Add(time.Duration(c.cfg.KeepAlive)))
}

func (c *GRPCClient) keepAliveDeadLine() time.Time {
	return c.keepaliveDeadLine.Load().(time.Time)
}

func (c *GRPCClient) sendPing() {
	c.rpcSend(RPC_PING)
}

func (c *GRPCClient) updateStatus(status string) {

	if status != c.Status() {

		c.status.Store(status)

		if status == TRANSPORT_STATUS_OFFLINE {
			c.diag.Debug("grpc disconnected")

			c.conn.Close()
			c.stream.CloseSend()

			// 断线重连
			go c.connect()
		}

		c.handlerStatus(status)
	}
}

func (c *GRPCClient) manage() {
	//nextWriteThreshold := time.Now()
	for {
		select {
		case <-time.After(time.Duration(c.cfg.KeepAlive)):
			if c.Status() == TRANSPORT_STATUS_OFFLINE {
				continue
			}
			if c.KeepAliveCount() >= KEEP_ALIVE_CHECK {
				go c.updateStatus(TRANSPORT_STATUS_OFFLINE)
				c.updateKeepAliveCount(0)
				continue
			}
			if c.keepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				c.sendPing()
				c.updateKeepAliveDeadLine() //更新keepalivedeadline
				c.addKeepAliveCount()
			}

		case buf := <-c.recvBuf:
			c.handleRpcRecv(buf)

		case stopDone := <-c.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (c *GRPCClient) connect() {
	var err error = nil

	for {
		c.diag.Debug("grpc connecting ...\n")
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithBackoffConfig(grpc.BackoffConfig{
			MaxDelay: 1 * time.Second,
		}))

		c.conn, err = grpc.Dial(c.cfg.GRPCServer, opts...)
		if err == nil {
			c.diag.Debug("Grpc Connected\n")
			c.updateStatus(TRANSPORT_STATUS_ONLINE)
			break
		}

		time.Sleep(1 * time.Second)
	}

	c.rpcClient = NewRPCAiisClient(c.conn)
	c.stream, _ = c.rpcClient.RPCNode(context.Background())

	go c.recvProcess()
}

func (c *GRPCClient) recvProcess() {
	for {
		if c.stream == nil {
			continue
		}

		in, err := c.stream.Recv()
		if err != nil {
			c.diag.Debug("Rpc recvProcess Exit")
			return
		}

		c.updateKeepAliveCount(0)
		c.recvBuf <- in.Payload
	}
}

func (s *GRPCClient) handleRpcRecv(payload string) {

	rpcPayload := TransportPayload{}
	if err := json.Unmarshal([]byte(payload), &rpcPayload); err != nil {
		s.diag.Error("handleRpcRecv", err)
	}
	strData, _ := json.Marshal(rpcPayload.Data)

	switch rpcPayload.Method {
	case TRANSPORT_METHOD_RESULT_PATCH:
		rp := ResultPatch{}
		json.Unmarshal(strData, &rp)
		s.handlerResultPatch(rp)
		break

	case TRANSPORT_METHOD_SERVICE_STATUS:
		serviceStatus := ServiceStatus{}
		json.Unmarshal(strData, &serviceStatus)
		s.handlerServiceStatus(serviceStatus)
		break

		//case TYPE_ORDER_START:
		//	// TODO: 开工响应
		//	//s.notifyService.WSSendMes(wsnotify.WS_EVENT_MES,"123","mes允许开工")
		//	break
		//
		//case TYPE_ORDER_FINISH:
		//	// TODO: 完工响应
		//	//s.notifyService.WSSendMes(wsnotify.WS_EVENT_MES,"123","mes确认完工")
		//	break
	}
}

func (c *GRPCClient) rpcSend(payload string) error {

	if c.stream != nil {
		return c.stream.Send(&Payload{
			Payload: payload,
		})
	}

	return errors.New("rpc not connected")
}

func (c *GRPCClient) SendResult(result *AIISResult) error {
	str, _ := json.Marshal(TransportPayload{
		Method: TRANSPORT_METHOD_RESULT,
		Data:   result,
	})

	return c.rpcSend(string(str))
}
