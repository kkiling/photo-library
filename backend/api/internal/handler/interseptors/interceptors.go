package interseptors

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/handler/descriptor"
	"github.com/kkiling/photo-library/backend/api/pkg/common/grpc_server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"google.golang.org/grpc"
)

func NewAuthInterceptor(logger log.Logger, descriptors descriptor.MethodDescriptorMap) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger = logger.WithCtx(ctx, "middleware", "NewTestTestInterceptor")

		ds, ok := descriptors.GetByFullName(info.FullMethod)
		if !ok {
			return nil, grpc_server.ErrUnauthenticated(descriptor.ErrMethodDescriptorNotFound)
		}

		if ds.UseAuth == 10 {
			logger.Errorf("%s - no auth", info.FullMethod)
			return nil, grpc_server.ErrPermissionDenied(fmt.Errorf("no auth"))
		}

		logger.Info("%s", info.FullMethod)
		return handler(ctx, req)
	}
}

func NewPanicRecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = grpc_server.ErrInternal(err)
			}
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}
