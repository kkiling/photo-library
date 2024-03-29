package photosservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/internal/handler/mapper"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"path/filepath"
	"strconv"
)

type PhotosService interface {
	GetPhotoGroups(ctx context.Context, req *photos.GetPhotoGroupsRequest) (*photos.PaginatedPhotoGroups, error)
	GetPhotoContent(ctx context.Context, fileName string, previewSize *int) (*photos.PhotoContent, error)
	GetPhotoGroup(ctx context.Context, groupID uuid.UUID) (*photos.PhotoGroupData, error)
}

type HandlerPhotosService struct {
	// desc.UnimplementedPhotosServiceServer
	logger        log.Logger
	photosService PhotosService
}

func NewHandlerPhotosService(logger log.Logger, photosService PhotosService) *HandlerPhotosService {
	return &HandlerPhotosService{
		logger:        logger,
		photosService: photosService,
	}
}

func (p *HandlerPhotosService) GetPhotoGroup(ctx context.Context, req *desc.GetPhotoGroupRequest) (*desc.PhotoGroupData, error) {
	groupID, err := uuid.ParseBytes([]byte(req.GroupId))
	if err != nil {
		return nil, server.ErrInvalidArgument(err)
	}

	response, err := p.photosService.GetPhotoGroup(ctx, groupID)
	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoGroup")
	}

	return mapper.GetPhotoGroupResponse(response), nil
}

func (p *HandlerPhotosService) GetPhotoGroups(ctx context.Context, request *desc.GetPhotoGroupsRequest) (*desc.PaginatedPhotoGroups, error) {
	response, err := p.photosService.GetPhotoGroups(ctx, mapper.GetPhotoGroupsRequest(request))

	if err != nil {
		return nil, handler.HandleError(err, "p.photosService.GetPhotoGroups")
	}

	return mapper.GetPhotoGroupsResponse(response), nil
}

func (p *HandlerPhotosService) GetPhotoContent(w http.ResponseWriter, r *http.Request) {
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
		p.logger.Errorf(" w.Write: %v", err)
		http.Error(w, "w.Write(photoContent.PhotoBody)", http.StatusInternalServerError)
	}
}

func (p *HandlerPhotosService) SetMainPhotoGroup(context.Context, *desc.SetMainPhotoGroupRequest) (*emptypb.Empty, error) {
	return nil, server.NewResponseError(codes.Unimplemented, fmt.Errorf("SetMainPhotoGroup Unimplemented"))
}

func (p *HandlerPhotosService) RegistrationServerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/photos/", p.GetPhotoContent)
}

func (p *HandlerPhotosService) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterPhotosServiceHandlerFromEndpoint
}

func (p *HandlerPhotosService) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterPhotosServiceServer(server, p)
}

func (p *HandlerPhotosService) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		handler.NewCustomDescriptor((*HandlerPhotosService).GetPhotoGroup),
		handler.NewCustomDescriptor((*HandlerPhotosService).GetPhotoGroups),
		handler.NewCustomDescriptor((*HandlerPhotosService).SetMainPhotoGroup),
	}
}
