package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
)

func getDescriptor(descriptors method_descriptor.MethodDescriptorMap, fullName string) *customDescriptor {
	ds, ok := descriptors.GetByFullName(fullName)
	if !ok {
		return nil
	}

	if result, ok := ds.(*customDescriptor); !ok {
		panic("cannot convert method descriptor to customDescriptor")
	} else {
		return result
	}
}

func NewAuthInterceptor(logger log.Logger, descriptors method_descriptor.MethodDescriptorMap) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// logger = logger.WithCtx(ctx, "middleware", "NewAuthInterceptor")

		ds := getDescriptor(descriptors, info.FullMethod)
		if ds == nil {
			return nil, server.ErrUnauthenticated(method_descriptor.ErrMethodDescriptorNotFound)
		}

		// TODO: логика авторизации
		// logger.Infof("auth method: %s", info.FullMethod)

		return handler(ctx, req)
	}
}

func NewPanicRecoverInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = server.ErrInternal(err)
			}
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

func NewLoggerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, err
		}
		switch status.Code(err) {
		case codes.Internal:
			logger.Errorf(err.Error())
		default:
			logger.Warnf(err.Error())
		}
		return resp, err
	}
}
