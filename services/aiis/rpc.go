package aiis

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
	"sync"
)

type OnRPCNewClient func(stream *RPCAiis_RPCNodeServer)
type OnRPCRecv func(stream *RPCAiis_RPCNodeServer, payload string)

type GRPCServer struct {
	grpcServer   *grpc.Server
	RPCRecv      OnRPCRecv
	RPCNewClient OnRPCNewClient
	clients      map[string]*RPCAiis_RPCNodeServer
	mtx_clients  sync.Mutex
}

func (s *GRPCServer) Start(port int) error {
	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)
	RegisterRPCAiisServer(s.grpcServer, s)
	s.clients = make(map[string]*RPCAiis_RPCNodeServer)
	s.mtx_clients = sync.Mutex{}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	go s.grpcServer.Serve(lis)

	return nil
}

func (s *GRPCServer) Stop() error {
	s.grpcServer.GracefulStop()

	return nil
}

func (s *GRPCServer) addClient(addr string, stream *RPCAiis_RPCNodeServer) {
	defer s.mtx_clients.Unlock()
	s.mtx_clients.Lock()

	s.clients[addr] = stream
}

func (s *GRPCServer) removeClient(addr string) {
	defer s.mtx_clients.Unlock()
	s.mtx_clients.Lock()

	delete(s.clients, addr)
}

func (s *GRPCServer) RPCNode(stream RPCAiis_RPCNodeServer) error {

	if s.RPCNewClient != nil {
		s.RPCNewClient(&stream)
	}

	pr, _ := peer.FromContext(stream.Context())
	s.addClient(pr.Addr.String(), &stream)

	for {
		in, err := stream.Recv()
		if err != nil {
			s.removeClient(pr.Addr.String())
			return err
		}

		//fmt.Printf("%s\n", in.Payload)

		if s.RPCRecv != nil {
			s.RPCRecv(&stream, in.Payload)
		}

	}

	return nil
}

func (s *GRPCServer) RPCSend(stream *RPCAiis_RPCNodeServer, payload string) error {
	if stream != nil {
		return (*stream).Send(&Payload{
			Payload: payload,
		})
	}

	return errors.New("rpc disconnected")
}

func (s *GRPCServer) RPCSendAll(payload string) error {
	defer s.mtx_clients.Unlock()
	s.mtx_clients.Lock()

	for _, stream := range s.clients {
		s.RPCSend(stream, payload)
	}

	return nil
}
