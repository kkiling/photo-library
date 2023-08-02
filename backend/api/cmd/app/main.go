package main

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/log"
	"github.com/kkiling/photo-library/backend/api/internal/server"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"google.golang.org/grpc"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	photoService := service.NewPhotoService()
	descriptor := server.Descriptor{
		GatewayRegistrar: pbv1.RegisterPhotosServiceHandlerFromEndpoint,
		OnRegisterGrpcServer: func(grpcServer *grpc.Server) {
			pbv1.RegisterPhotosServiceServer(grpcServer, photoService)
		},
	}

	interceptors := []grpc.UnaryServerInterceptor{
		photoService.NewTestTestInterceptor(),
	}

	var cfg = server.Config{
		Host:                    "127.0.0.1",
		GrpcPort:                8080,
		HttpPort:                8181,
		MaxSendMessageLength:    2147483647,
		MaxReceiveMessageLength: 63554432,
	}
	srv := server.NewServer(log.NewLogger(), cfg, interceptors...)

	if err := srv.Register(ctx, descriptor); err != nil {
		panic(err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	srv.Stop()

	fmt.Println("EndStop")
}
