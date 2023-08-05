package handler

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"github.com/prometheus/client_golang/prometheus"
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

	// Создайте метрику
	myCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_custom_counter",
		Help: "This is my counter",
	})

	// Регистрируйте метрику
	prometheus.MustRegister(myCounter)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger = logger.WithCtx(ctx, "middleware", "NewTestTestInterceptor")

		ds := getDescriptor(descriptors, info.FullMethod)
		if ds == nil {
			return nil, server.ErrUnauthenticated(method_descriptor.ErrMethodDescriptorNotFound)
		}

		if ds.useAuth {
			logger.Errorf("Call %s - no auth", info.FullMethod)
			return nil, server.ErrPermissionDenied(fmt.Errorf("no auth"))
		}

		// Обновите значение метрики
		myCounter.Inc()

		// logger.Infof("Call %s", info.FullMethod)
		return handler(ctx, req)
	}
}

func NewPanicRecoverInterceptor() grpc.UnaryServerInterceptor {
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
