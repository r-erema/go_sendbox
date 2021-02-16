package microservices

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type bizServiceImpl struct {
	aclRules     ACL
	adminService *adminServiceImpl
}

func (bizService *bizServiceImpl) isAccessAllowed(ctx context.Context) bool {
	return checkAccessByContext(ctx, bizService.aclRules)
}

func (bizService *bizServiceImpl) Check(ctx context.Context, data *Nothing) (*Nothing, error) {

	bizService.adminService.SaveLogEvent(createLogEventByContext(ctx))

	if !bizService.isAccessAllowed(ctx) {
		return nil, status.Error(codes.Unauthenticated, "Access denied")
	}

	return &Nothing{}, nil
}
func (bizService *bizServiceImpl) Add(ctx context.Context, data *Nothing) (*Nothing, error) {

	bizService.adminService.SaveLogEvent(createLogEventByContext(ctx))

	_ = data
	if !bizService.isAccessAllowed(ctx) {
		return nil, status.Error(codes.Unauthenticated, "Access denied")
	}
	return &Nothing{}, nil
}

func (bizService *bizServiceImpl) Test(ctx context.Context, data *Nothing) (*Nothing, error) {

	bizService.adminService.SaveLogEvent(createLogEventByContext(ctx))

	_ = data
	if !bizService.isAccessAllowed(ctx) {
		return nil, status.Error(codes.Unauthenticated, "Access denied")
	}
	return &Nothing{}, nil
}
