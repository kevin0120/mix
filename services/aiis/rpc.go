package aiis

import (
	"github.com/kataras/iris/core/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

type OnRPCRecv func(payload string)

type GRPCClient struct {
	srv       *Service
	conn      *grpc.ClientConn
	stream    RPCAiis_RPCNodeClient
	opts      []grpc.DialOption
	rpcClient RPCAiisClient
	RPCRecv   OnRPCRecv
}

func (c *GRPCClient) Start() error {
	go c.Connect()

	return nil
}

func (c *GRPCClient) Stop() error {
	c.conn.Close()
	c.stream.CloseSend()

	return nil
}

func (c *GRPCClient) Connect() {
	var err error = nil

	for {
		c.srv.diag.Info("grpc connecting ...\n")
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithBackoffConfig(grpc.BackoffConfig{
			MaxDelay: 1 * time.Second,
		}))

		c.conn, err = grpc.Dial(c.srv.Config().GRPCServer, opts...)
		if err == nil {
			c.srv.diag.Info("grpc connected\n")
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
			go c.Connect()
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

	return errors.New("rpc disconnected")
}
