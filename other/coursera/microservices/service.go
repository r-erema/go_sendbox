package microservices

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ACL map[string][]string

func StartMyMicroservice(ctx context.Context, listenAddr, ACLJson string) error {

	aclData := make(ACL)
	err := json.Unmarshal([]byte(ACLJson), &aclData)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	adminService := &adminServiceImpl{aclRules: aclData, logsCh: make(chan *Event), logsStorage: &LogsStorage{}}
	RegisterAdminServer(grpcServer, adminService)
	RegisterBizServer(grpcServer, &bizServiceImpl{aclData, adminService})

	go func(lis net.Listener) {
		if err := grpcServer.Serve(lis); err != nil {
			log.Print(err)
			return
		}
	}(listener)

	go func(context context.Context, server *grpc.Server) {
		for {
			select {
			case <-context.Done():
				server.GracefulStop()
				return
			}
		}
	}(ctx, grpcServer)
	return nil
}
