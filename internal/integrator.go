package internal

import (
	pb "../api"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type WootServer struct {}

func (s *WootServer) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoReply, error) {
	return &pb.EchoReply{Body: in.GetBody()}, nil
}

func StartServer(port int) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterWootServer(s, &WootServer{})

	go func () {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return s
}
