package handler

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/handler/descriptor"
	"github.com/kkiling/photo-library/backend/api/internal/handler/interseptors"
	pbv1 "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/grpc_server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/provider"
	"google.golang.org/grpc"
)

type PhotosServiceServer struct {
	server      *grpc_server.Server
	logger      log.Logger
	cfgProvider provider.Provider
}

func NewPhotosServiceServer(logger log.Logger, cfgProvider provider.Provider) *PhotosServiceServer {
	return &PhotosServiceServer{
		logger:      logger,
		cfgProvider: cfgProvider,
	}
}

func (p *PhotosServiceServer) crateServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	descriptors, err := descriptor.NewMethodDescriptorMap(
		[]descriptor.MethodDescriptor{
			{
				Method:  (*PhotosServiceServer).UploadPhoto,
				UseAuth: 10,
			},
			{
				Method:  (*PhotosServiceServer).CheckHashPhoto,
				UseAuth: 11,
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("descriptor.NewMethodDescriptorMap: %w", err)
	}

	return []grpc.UnaryServerInterceptor{
		interseptors.NewAuthInterceptor(p.logger, descriptors),
		interseptors.NewPanicRecoverInterceptor(),
	}, nil
}

func (p *PhotosServiceServer) Start(ctx context.Context) error {
	inters, err := p.crateServerInterceptors()
	if err != nil {
		return fmt.Errorf("crateServerInterceptors: %w", err)
	}

	// Создать конфиг
	var cfg grpc_server.Config
	err = p.cfgProvider.PopulateByKey("server", &cfg)
	if err != nil {
		return fmt.Errorf("PopulateByKey: %w", err)
	}

	p.server = grpc_server.NewServer(p.logger, cfg, inters...)
	serverDs := grpc_server.Descriptor{
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
