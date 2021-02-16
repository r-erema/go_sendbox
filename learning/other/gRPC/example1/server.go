package example1

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	_ = ctx
	log.Printf("Received: %v", in.GetName())
	return &HelloReply{Message: fmt.Sprintf("Hello %s", in.GetName())}, nil
}

func CreateAndRunServer(listener net.Listener) *grpc.Server {

	s := grpc.NewServer()
	RegisterGreeterServer(s, &server{})

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return s
}
