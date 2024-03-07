package handler

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/service/syncphotos"
	"net/http"

	"github.com/kkiling/photo-library/backend/api/internal/handler/mapper"

	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
)

type customDescriptor struct {
	method interface{}
}

func (c *customDescriptor) Method() interface{} {
	return c.method
}

type SyncPhotosServiceServer struct {
	desc.UnimplementedSyncPhotosServiceServer
	server       *server.Server
	logger       log.Logger
	serverConfig server.Config
	syncPhoto    SyncPhotosService
}

func NewSyncPhotosServiceServer(logger log.Logger,
	syncPhoto SyncPhotosService,
	serverConfig server.Config) *SyncPhotosServiceServer {
	return &SyncPhotosServiceServer{
		logger:       logger,
		serverConfig: serverConfig,
		syncPhoto:    syncPhoto,
	}
}

func (p *SyncPhotosServiceServer) crateServerInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	descriptors, err := method_descriptor.NewMethodDescriptorMap(
		[]method_descriptor.Descriptor{
			&customDescriptor{
				method: (*SyncPhotosServiceServer).UploadPhoto,
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

func (p *SyncPhotosServiceServer) registrationServerHandlers(*http.ServeMux) {

}

func (p *SyncPhotosServiceServer) Start(ctx context.Context) error {
	interceptors, err := p.crateServerInterceptors()
	if err != nil {
		return fmt.Errorf("crateServerInterceptors: %w", err)
	}

	p.server = server.NewServer(p.logger, p.serverConfig, interceptors...)
	serverDs := server.Descriptor{
		GatewayRegistrar: desc.RegisterSyncPhotosServiceHandlerFromEndpoint,
		OnRegisterGrpcServer: func(grpcServer *grpc.Server) {
			desc.RegisterSyncPhotosServiceServer(grpcServer, p)
		},
	}

	if err := p.server.Register(ctx, serverDs); err != nil {
		return fmt.Errorf("server.Register: %w", err)
	}

	if err := p.server.Start("sync_photos_service", p.registrationServerHandlers); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

func (p *SyncPhotosServiceServer) Stop() {
	p.server.Stop()
}

type SyncPhotosService interface {
	UploadPhoto(ctx context.Context, form *syncphotos.SyncPhotoRequest) (*syncphotos.SyncPhotoResponse, error)
}

func (p *SyncPhotosServiceServer) UploadPhoto(ctx context.Context, request *desc.UploadPhotoRequest) (*desc.UploadPhotoResponse, error) {
	response, err := p.syncPhoto.UploadPhoto(ctx, mapper.UploadPhotoRequest(request))

	if err != nil {
		return nil, handleError(err, "p.syncPhoto.UploadPhoto")
	}

	return mapper.UploadPhotoResponse(response), nil
}
