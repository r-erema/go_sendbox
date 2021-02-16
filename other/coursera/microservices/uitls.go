package microservices

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"regexp"
	"time"
)

func checkAccess(consumer, path string, rules ACL) bool {
	if routePatterns, ok := rules[consumer]; ok {
		for _, pattern := range routePatterns {
			if ok, _ = regexp.MatchString(pattern, path); ok {
				return true
			}
		}
	}
	return false
}

func checkAccessByContext(ctx context.Context, rules ACL) bool {
	if consumer, ok := getConsumerFromContext(ctx); ok {
		if method, ok := grpc.Method(ctx); ok && checkAccess(consumer, method, rules) {
			return true
		}
	}
	return false
}

func createLogEventByContext(ctx context.Context) *Event {
	consumer, consumerOk := getConsumerFromContext(ctx)
	method, methodOk := grpc.Method(ctx)
	userPeer, peerOk := peer.FromContext(ctx)
	if consumerOk && methodOk && peerOk {
		return &Event{
			Timestamp: time.Now().UnixNano(),
			Consumer:  consumer,
			Method:    method,
			Host:      userPeer.Addr.String(),
		}
	}
	return nil
}

func getConsumerFromContext(ctx context.Context) (string, bool) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if consumer, ok := md["consumer"]; ok {
			return consumer[0], true
		}
	}
	return "", false
}
