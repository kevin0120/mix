package aiis

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type OnRPCNewClient func(stream *RPCAiis_RPCNodeServer)
type OnRPCRecv func(stream *RPCAiis_RPCNodeServer, payload string)

type GRPCServer struct {
	grpcServer   *grpc.Server
	RPCRecv      OnRPCRecv
	RPCNewClient OnRPCNewClient
}

func (s *GRPCServer) Start(port int) error {
	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)
	RegisterRPCAiisServer(s.grpcServer, s)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
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

func (s *GRPCServer) RPCNode(stream RPCAiis_RPCNodeServer) error {

	if s.RPCNewClient != nil {
		s.RPCNewClient(&stream)
	}

	for {
		in, err := stream.Recv()
		if err != nil {
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
