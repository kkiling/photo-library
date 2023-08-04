package handler

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/cfg"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
)

type customDescriptor struct {
	method  interface{}
	useAuth bool
}

func (c *customDescriptor) Method() interface{} {
	return c.method
}

type PhotosServiceServer struct {
	server      *server.Server
	logger      log.Logger
	cfgProvider config.Provider
}

func NewPhotosServiceServer(logger log.Logger, cfgProvider config.Provider) *PhotosServiceServer {
	return &PhotosServiceServer{
		logger:      logger,
		cfgProvider: cfgProvider,
	}
}

func (p *PhotosServiceServer) crateServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	descriptors, err := method_descriptor.NewMethodDescriptorMap(
		[]method_descriptor.Descriptor{
			&customDescriptor{
				method:  (*PhotosServiceServer).UploadPhoto,
				useAuth: false,
			},
			&customDescriptor{
				method:  (*PhotosServiceServer).CheckHashPhoto,
				useAuth: false,
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("method_descriptor.NewMethodDescriptorMap: %w", err)
	}

	return []grpc.UnaryServerInterceptor{
		NewAuthInterceptor(p.logger, descriptors),
		NewPanicRecoverInterceptor(),
	}, nil
}

func (p *PhotosServiceServer) Start(ctx context.Context) error {
	interceptors, err := p.crateServerInterceptors()
	if err != nil {
		return fmt.Errorf("crateServerInterceptors: %w", err)
	}

	var serverConfig server.Config
	err = p.cfgProvider.PopulateByKey(cfg.ServerConfigName, &serverConfig)
	if err != nil {
		return fmt.Errorf("PopulateByKey: %w", err)
	}

	p.server = server.NewServer(p.logger, serverConfig, interceptors...)
	serverDs := server.Descriptor{
		GatewayRegistrar: pbv1.RegisterPhotosServiceHandlerFromEndpoint,
		OnRegisterGrpcServer: func(grpcServer *grpc.Server) {
			pbv1.RegisterPhotosServiceServer(grpcServer, p)
		},
	}

	if err := p.server.Register(ctx, serverDs); err != nil {
		return fmt.Errorf("server.Register: %w", err)
	}

	if err := p.server.Start(); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

func (p *PhotosServiceServer) Stop() {
	p.server.Stop()
}

func (p *PhotosServiceServer) CheckHashPhoto(ctx context.Context, request *pbv1.CheckHashPhotoRequest) (*pbv1.CheckHashPhotoResponse, error) {
	p.logger.Info("CheckHashPhoto")
	return &pbv1.CheckHashPhotoResponse{AlreadyLoaded: true}, nil
}

func (p *PhotosServiceServer) UploadPhoto(ctx context.Context, request *pbv1.UploadPhotoRequest) (*pbv1.UploadPhotoResponse, error) {
	return &pbv1.UploadPhotoResponse{Success: true}, nil

}
