package handler

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/cfg"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type customDescriptor struct {
	method interface{}
}

func (c *customDescriptor) Method() interface{} {
	return c.method
}

type SyncPhotosServiceServer struct {
	desc.UnimplementedSyncPhotosServiceServer
	server      *server.Server
	logger      log.Logger
	cfgProvider config.Provider
	syncPhoto   SyncPhotosService
}

func NewSyncPhotosServiceServer(logger log.Logger,
	syncPhoto SyncPhotosService,
	cfgProvider config.Provider) *SyncPhotosServiceServer {
	return &SyncPhotosServiceServer{
		logger:      logger,
		cfgProvider: cfgProvider,
		syncPhoto:   syncPhoto,
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
		NewAuthInterceptor(p.logger, descriptors),
		NewPanicRecoverInterceptor(),
	}, nil
}

func (p *SyncPhotosServiceServer) Start(ctx context.Context) error {
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
		GatewayRegistrar: desc.RegisterSyncPhotosServiceHandlerFromEndpoint,
		OnRegisterGrpcServer: func(grpcServer *grpc.Server) {
			desc.RegisterSyncPhotosServiceServer(grpcServer, p)
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

func (p *SyncPhotosServiceServer) Stop() {
	p.server.Stop()
}

type SyncPhotosService interface {
	UploadPhoto(ctx context.Context, form model.SyncPhotoRequest) (model.SyncPhotoResponse, error)
}

func (p *SyncPhotosServiceServer) UploadPhoto(ctx context.Context, request *desc.UploadPhotoRequest) (*desc.UploadPhotoResponse, error) {

	response, err := p.syncPhoto.UploadPhoto(ctx, model.SyncPhotoRequest{
		Paths:    request.Paths,
		Hash:     request.Hash,
		Body:     request.Body,
		UpdateAt: request.UpdateAt.AsTime(),
		ClientId: request.ClientId,
	})

	if err != nil {
		return nil, err
	}

	return &desc.UploadPhotoResponse{
		HasBeenUploadedBefore: response.HasBeenUploadedBefore,
		Hash:                  response.Hash,
		UploadedAt: &timestamppb.Timestamp{
			Seconds: response.UploadAt.Unix(),
		},
	}, nil
}
