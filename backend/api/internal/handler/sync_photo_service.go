package handler

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/cfg"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
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
}

func NewSyncPhotosServiceServer(logger log.Logger, cfgProvider config.Provider) *SyncPhotosServiceServer {
	return &SyncPhotosServiceServer{
		logger:      logger,
		cfgProvider: cfgProvider,
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

func (p *SyncPhotosServiceServer) UploadPhoto(ctx context.Context, request *desc.UploadPhotoRequest) (*desc.UploadPhotoResponse, error) {
	/*// Извлекаем расширение файла
	ext := filepath.Ext(request.Paths[0])
	// Создаем UUID
	uuid := uuid.New()

	// Формируем новое имя файла
	newFilename := fmt.Sprintf("/Users/kkiling/Desktop/photos/%s%s", uuid.String(), ext)

	// Создаем новый файл с новым именем
	newFile, err := os.Create(newFilename)
	defer newFile.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new file: %v", err)
	}

	// Записываем данные в новый файл
	if _, err := newFile.Write(request.Body); err != nil {
		return nil, fmt.Errorf("Failed to write to new file: %v", err)
	}
	*/
	return &desc.UploadPhotoResponse{
		Success: true,
		Hash:    request.Hash,
		UploadedAt: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}, nil
}
