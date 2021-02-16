package example1

import (
	"google.golang.org/grpc"
	"log"
)

func CreateClient(address string) GreeterClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("didn't connect: %v", err)
	}
	return NewGreeterClient(conn)
}
