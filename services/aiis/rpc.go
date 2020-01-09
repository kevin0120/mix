package aiis

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

const (
	RpcPing        = "ping"
	KeepAliveCheck = 3
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
	keepAliveCount    atomic.Int32
	keepaliveDeadLine atomic.Value
	closing           chan chan struct{}
	recvBuf           chan string

	BaseTransport
}

func (s *GRPCClient) Start() error {
	s.status.Store(TransportStatusOffline)
	s.closing = make(chan chan struct{})
	s.updateKeepAliveDeadLine()

	go s.connect()
	go s.manage()

	return nil
}

func (s *GRPCClient) Stop() error {
	s.grpcClose()

	closed := make(chan struct{})
	s.closing <- closed

	return nil
}

func (s *GRPCClient) grpcClose() {
	if s.conn == nil || s.stream == nil {
		return
	}

	if err := s.conn.Close(); err != nil {
		s.diag.Error("Close GRPC Conn Failed", err)
	}

	if err := s.stream.CloseSend(); err != nil {
		s.diag.Error("Close GRPC Stream Send Failed", err)
	}
}

func (s *GRPCClient) Status() string {

	return s.status.Load().(string)
}

func (s *GRPCClient) KeepAliveCount() int32 {
	return s.keepAliveCount.Load()
}

func (s *GRPCClient) updateKeepAliveCount(i int32) {
	s.keepAliveCount.Swap(i)
}

func (s *GRPCClient) addKeepAliveCount() {
	s.keepAliveCount.Inc()
}

func (s *GRPCClient) updateKeepAliveDeadLine() {
	s.keepaliveDeadLine.Store(time.Now().Add(time.Duration(s.cfg.KeepAlive)))
}

func (s *GRPCClient) keepAliveDeadLine() time.Time {
	return s.keepaliveDeadLine.Load().(time.Time)
}

func (s *GRPCClient) sendPing() {
	if err := s.rpcSend(RpcPing); err != nil {
		s.diag.Error("GRPC Send Failed", err)
	}
}

func (s *GRPCClient) updateStatus(status string) {

	if status != s.Status() {

		s.status.Store(status)

		if status == TransportStatusOffline {
			s.diag.Debug("GRPC Disconnected")
			s.grpcClose()

			// 断线重连
			go s.connect()
		}

		s.handlerStatus(status)
	}
}

func (s *GRPCClient) manage() {
	for {
		select {
		case <-time.After(time.Duration(s.cfg.KeepAlive)):
			if s.Status() == TransportStatusOffline {
				continue
			}
			if s.KeepAliveCount() >= KeepAliveCheck {
				go s.updateStatus(TransportStatusOffline)
				s.updateKeepAliveCount(0)
				continue
			}
			if s.keepAliveDeadLine().Before(time.Now()) {
				//到达了deadline
				s.sendPing()
				s.updateKeepAliveDeadLine() //更新keepalivedeadline
				s.addKeepAliveCount()
			}

		case buf := <-s.recvBuf:
			s.handleRpcRecv(buf)

		case stopDone := <-s.closing:
			close(stopDone)
			return //退出manage协程
		}
	}
}

func (s *GRPCClient) connect() {
	var err error = nil

	for {
		s.diag.Debug("grpc connecting ...\n")
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithBackoffConfig(grpc.BackoffConfig{
			MaxDelay: 1 * time.Second,
		}))

		s.conn, err = grpc.Dial(s.cfg.GRPCServer, opts...)
		if err == nil {
			s.diag.Debug("Grpc Connected\n")
			s.updateStatus(TransportStatusOnline)
			break
		}

		time.Sleep(1 * time.Second)
	}

	s.rpcClient = NewRPCAiisClient(s.conn)
	s.stream, _ = s.rpcClient.RPCNode(context.Background())

	go s.recvProcess()
}

func (s *GRPCClient) recvProcess() {
	for {
		if s.stream == nil {
			continue
		}

		in, err := s.stream.Recv()
		if err != nil {
			s.diag.Debug("Rpc recvProcess Exit")
			return
		}

		s.updateKeepAliveCount(0)
		s.recvBuf <- in.Payload
	}
}

func (s *GRPCClient) handleRpcRecv(payload string) {

	rpcPayload := TransportPayload{}
	if err := json.Unmarshal([]byte(payload), &rpcPayload); err != nil {
		s.diag.Error("handleRpcRecv", err)
	}
	strData, _ := json.Marshal(rpcPayload.Data)

	switch rpcPayload.Method {
	case TransportMethodResultPatch:
		rp := ResultPatch{}
		if err := json.Unmarshal(strData, &rp); err != nil {
			s.diag.Error("Unmarshal ResultPatch Failed", err)
			return
		}
		s.handlerResultPatch(rp)

	case TransportMethodServiceStatus:
		serviceStatus := ServiceStatus{}
		if err := json.Unmarshal(strData, &serviceStatus); err != nil {
			s.diag.Error("Unmarshal ServiceStatus Failed", err)
			return
		}
		s.handlerServiceStatus(serviceStatus)

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

func (s *GRPCClient) rpcSend(payload string) error {

	if s.stream != nil {
		return s.stream.Send(&Payload{
			Payload: payload,
		})
	}

	return errors.New("rpc not connected")
}

func (s *GRPCClient) SendResult(result *PublishResult) error {
	str, _ := json.Marshal(TransportPayload{
		Method: TransportMethodResult,
		Data:   result,
	})

	return s.rpcSend(string(str))
}
