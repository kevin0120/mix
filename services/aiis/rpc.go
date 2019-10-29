package aiis

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
	"sync"
)

type Diagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type OnRPCNewClient func(stream RPCAiis_RPCNodeServer)
type OnRPCRecv func(stream RPCAiis_RPCNodeServer, payload string)

type GRPCServer struct {
	server             *grpc.Server
	OnAIISRPCRecv      OnRPCRecv
	OnAIISRPCNewClient OnRPCNewClient
	clients            map[string]RPCAiis_RPCNodeServer
	mtxClients         sync.Mutex
	diag               Diagnostic
}

func NewAiisGrpcServer(d Diagnostic) *GRPCServer {
	var opts []grpc.ServerOption
	s := &GRPCServer{
		server:  grpc.NewServer(opts...),
		clients: make(map[string]RPCAiis_RPCNodeServer),
		diag:    d,
	}

	return s
}

func (s *GRPCServer) Start(port int) error {

	RegisterRPCAiisServer(s.server, s)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	go s.server.Serve(lis)

	return nil
}

func (s *GRPCServer) Stop() error {
	s.server.GracefulStop()

	return nil
}

func (s *GRPCServer) addClient(addr string, stream RPCAiis_RPCNodeServer) {
	s.mtxClients.Lock()
	defer s.mtxClients.Unlock()

	s.clients[addr] = stream
}

func (s *GRPCServer) removeClient(addr string) {
	s.mtxClients.Lock()
	defer s.mtxClients.Unlock()

	delete(s.clients, addr)
}

func (s *GRPCServer) grpcHandleMsg(stream RPCAiis_RPCNodeServer, in *Payload) {
	switch in.Payload {
	case RPC_PING:
		msg := &Payload{
			Payload: RPC_PONG,
		}
		if err := stream.Send(msg); err != nil {
			s.diag.Error("grpcHandleMsg Send PONG Message", err)
		}
	default:
		if s.OnAIISRPCRecv == nil {
			err := errors.New("OnAIISRPCRecv Is Not Defined!!!\n")
			s.diag.Error("grpcHandleMsg", err)
		}
		s.OnAIISRPCRecv(stream, in.Payload)
	}
}

func (s *GRPCServer) RPCNode(stream RPCAiis_RPCNodeServer) error {

	if s.OnAIISRPCNewClient != nil {
		s.OnAIISRPCNewClient(stream)
	}

	pr, _ := peer.FromContext(stream.Context())
	s.addClient(pr.Addr.String(), stream)

	for {
		in, err := stream.Recv()
		if err != nil {
			s.removeClient(pr.Addr.String())
			return err
		}
		s.grpcHandleMsg(stream, in)
	}

}

func (s *GRPCServer) RPCSend(stream RPCAiis_RPCNodeServer, payload string) error {
	if stream != nil {
		return stream.Send(&Payload{
			Payload: payload,
		})
	}
	ctx := stream.Context()

	pr, _ := peer.FromContext(ctx)

	return errors.New(fmt.Sprintf("AIIS GRPC: %s Disconnected", pr.Addr.String()))
}

func (s *GRPCServer) RPCSendAll(payload string) error {
	s.mtxClients.Lock()
	defer s.mtxClients.Unlock()

	for _, stream := range s.clients {
		if err := s.RPCSend(stream, payload); err != nil {
			s.diag.Error("RPCSendAll Error", err)
		}
	}

	return nil
}
