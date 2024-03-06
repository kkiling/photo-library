package handler

import (
	"context"
	"fmt"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
)

type PhotosServiceServer struct {
	desc.UnimplementedPhotosServiceServer
	server        *server.Server
	logger        log.Logger
	serverConfig  server.Config
	photosService PhotosService
}

func NewPhotosServiceServer(logger log.Logger,
	photosService PhotosService,
	serverConfig server.Config) *PhotosServiceServer {
	return &PhotosServiceServer{
		logger:        logger,
		serverConfig:  serverConfig,
		photosService: photosService,
	}
}

func (p *PhotosServiceServer) crateServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	descriptors, err := method_descriptor.NewMethodDescriptorMap(
		[]method_descriptor.Descriptor{
			&customDescriptor{
				method: (*PhotosServiceServer).GetPhotoGroups,
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("method_descriptor.NewMethodDescriptorMap: %w", err)
	}

	return []grpc.UnaryServerInterceptor{
		NewPanicRecoverInterceptor(p.logger),
		NewLoggerInterceptor(p.logger),
		NewAuthInterceptor(p.logger, descriptors),
	}, nil
}

func (p *PhotosServiceServer) Start(ctx context.Context) error {
	interceptors, err := p.crateServerInterceptors()
	if err != nil {
		return fmt.Errorf("crateServerInterceptors: %w", err)
	}

	p.server = server.NewServer(p.logger, p.serverConfig, interceptors...)
	serverDs := server.Descriptor{
		GatewayRegistrar: desc.RegisterPhotosServiceHandlerFromEndpoint,
		OnRegisterGrpcServer: func(grpcServer *grpc.Server) {
			desc.RegisterPhotosServiceServer(grpcServer, p)
		},
	}

	if err := p.server.Register(ctx, serverDs); err != nil {
		return fmt.Errorf("server.Register: %w", err)
	}

	if err := p.server.Start("photos_service"); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

func (p *PhotosServiceServer) Stop() {
	p.server.Stop()
}

type PhotosService interface {
	// UploadPhoto(ctx context.Context, form *model.SyncPhotoRequest) (*model.SyncPhotoResponse, error)
}

func (p *PhotosServiceServer) GetPhotoGroups(context.Context, *desc.GetPhotoGroupsRequest) (*desc.GetPhotoGroupsResponse, error) {
	return &desc.GetPhotoGroupsResponse{
		Items:      []*desc.PhotoGroup{},
		TotalItems: 0,
	}, nil
}
