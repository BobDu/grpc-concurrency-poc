package main

import (
	"fmt"
	"log"
	"net"

	api "server/hello"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type HelloServer struct {
	api.UnimplementedHelloServiceServer
}

func (s *HelloServer) SayHello(ctx context.Context, req *api.HelloReq) (*api.HelloResp, error) {
	return &api.HelloResp{Result: fmt.Sprintf("Hey, %s!", req.GetName())}, nil
}

func Serve() {
	addr := fmt.Sprintf(":%d", 50051)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Cannot listen to address %s", addr)
	}
	s := grpc.NewServer()
	api.RegisterHelloServiceServer(s, &HelloServer{})
	if err := s.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	Serve()
}
