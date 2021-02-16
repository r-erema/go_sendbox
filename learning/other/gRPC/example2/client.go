package example2

import (
	"google.golang.org/grpc"
	"log"
)

func CreateClient(address string) RouteGuideClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("didn't connect: %v", err)
	}
	return NewRouteGuideClient(conn)
}
