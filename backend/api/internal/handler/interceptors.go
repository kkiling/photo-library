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

type CustomDescriptor struct {
	method interface{}
}

func NewCustomDescriptor(method interface{}) *CustomDescriptor {
	return &CustomDescriptor{
		method: method,
	}
}

func (c *CustomDescriptor) Method() interface{} {
	return c.method
}

func getCustomDescriptor(descriptors methoddescriptor.DescriptorsMap, fullName string) *CustomDescriptor {
	ds, ok := descriptors.GetByFullName(fullName)
	if !ok {
		return nil
	}

	if result, ok := ds.(*CustomDescriptor); !ok {
		panic("cannot convert method descriptor to customDescriptor")
	} else {
		return result
	}
}

func NewAuthInterceptor(logger log.Logger, descriptors methoddescriptor.DescriptorsMap) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ds := getCustomDescriptor(descriptors, info.FullMethod)
		if ds == nil {
			return nil, server.ErrUnauthenticated(methoddescriptor.ErrMethodDescriptorNotFound)
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
