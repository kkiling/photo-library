package photo_metadata_handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/kkiling/photo-library/backend/api/internal/interceptor"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type PhotoMetaDataService interface {
	GetPhotoMetaData(ctx context.Context, photoID uuid.UUID) (*model.PhotoMetadata, error)
}

type PhotoMetadataHandler struct {
	logger        log.Logger
	photosService PhotoMetaDataService
}

func NewHandler(logger log.Logger, photosService PhotoMetaDataService) *PhotoMetadataHandler {
	return &PhotoMetadataHandler{
		logger:        logger,
		photosService: photosService,
	}
}

func (p *PhotoMetadataHandler) RegistrationServerHandlers(_ *http.ServeMux) {
}

func (p *PhotoMetadataHandler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterPhotoMetadataServiceHandlerFromEndpoint
}

func (p *PhotoMetadataHandler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterPhotoMetadataServiceServer(server, p)
}

func (p *PhotoMetadataHandler) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		interceptor.NewCustomDescriptor((*PhotoMetadataHandler).GetPhotoMetaData),
	}
}
