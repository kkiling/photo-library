package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/handler/mapper"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
	"net/http"
	"path/filepath"
	"strconv"
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
			&customDescriptor{
				method: (*PhotosServiceServer).GetPhotoGroup,
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

func (p *PhotosServiceServer) registrationServerHandlers(mux *http.ServeMux) {
	// Настройка HTTP-хендлера для изображений
	mux.HandleFunc("/photos/", func(w http.ResponseWriter, r *http.Request) {
		fileName := filepath.Base(r.URL.Path)

		queryParams := r.URL.Query()
		sizeStr := queryParams.Get("size")
		var previewSize *int
		if sizeStr != "" {
			sizeInt, err := strconv.Atoi(sizeStr)
			if err != nil {
				http.Error(w, "invalid size parameter", http.StatusBadRequest)
				return
			}
			previewSize = &sizeInt
		}

		photoContent, err := p.photosService.GetPhotoContent(r.Context(), fileName, previewSize)
		if err != nil {
			if errors.Is(err, serviceerr.ErrNotFound) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, fmt.Errorf("p.photosService.GetPhotoContent: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		var contentType string
		switch photoContent.Extension {
		case model.PhotoExtensionJpg, model.PhotoExtensionJpeg:
			contentType = "image/jpeg"
		case model.PhotoExtensionPng:
			contentType = "image/png"
		case model.PhotoExtensionBmb:
			contentType = "image/bmp"
		default:
			http.Error(w, "Unsupported image format", http.StatusBadRequest)
			return
		}

		// Установка заголовка Content-Type и отправка изображения
		w.Header().Set("Content-Type", contentType)
		_, err = w.Write(photoContent.PhotoBody)
		if err != nil {
			http.Error(w, "w.Write(photoContent.PhotoBody)", http.StatusInternalServerError)
		}
	})
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

	if err := p.server.Start("photos_service", p.registrationServerHandlers); err != nil {
		return fmt.Errorf("server.Start: %w", err)
	}

	return nil
}

func (p *PhotosServiceServer) Stop() {
	p.server.Stop()
}

type PhotosService interface {
	GetPhotoGroups(ctx context.Context, req *photos.GetPhotoGroupsRequest) (*photos.GetPhotoGroupsResponse, error)
	GetPhotoContent(ctx context.Context, fileName string, previewSize *int) (*photos.PhotoContent, error)
	GetPhotoGroup(ctx context.Context, groupID uuid.UUID) (*photos.PhotoGroupData, error)
}

func (p *PhotosServiceServer) GetPhotoGroup(ctx context.Context, req *desc.GetPhotoGroupRequest) (*desc.GetPhotoGroupResponse, error) {
	groupID, err := uuid.ParseBytes([]byte(req.GroupId))
	if err != nil {
		return nil, server.ErrInvalidArgument(err)
	}

	response, err := p.photosService.GetPhotoGroup(ctx, groupID)
	if err != nil {
		return nil, handleError(err, "p.photosService.GetPhotoGroup")
	}

	return mapper.GetPhotoGroupResponse(response), nil
}

func (p *PhotosServiceServer) GetPhotoGroups(ctx context.Context, request *desc.GetPhotoGroupsRequest) (*desc.GetPhotoGroupsResponse, error) {
	response, err := p.photosService.GetPhotoGroups(ctx, mapper.GetPhotoGroupsRequest(request))

	if err != nil {
		return nil, handleError(err, "p.photosService.GetPhotoGroups")
	}

	return mapper.GetPhotoGroupsResponse(response), nil
}
