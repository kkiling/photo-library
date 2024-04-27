package photo_tags_handler

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

type PhotoTagsService interface {
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, tagName string) error
	DeletePhotoTag(ctx context.Context, tagID uuid.UUID) error
	GetPhotoTags(ctx context.Context, photoID uuid.UUID) ([]model.TagWithCategoryDTO, error)
	GetTagCategories(ctx context.Context) ([]model.TagCategory, error)
}

type PhotoTagsHandler struct {
	logger        log.Logger
	photosService PhotoTagsService
}

func NewHandler(logger log.Logger, photosService PhotoTagsService) *PhotoTagsHandler {
	return &PhotoTagsHandler{
		logger:        logger,
		photosService: photosService,
	}
}

func (p *PhotoTagsHandler) RegistrationServerHandlers(mux *http.ServeMux) {
}

func (p *PhotoTagsHandler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterPhotoTagsServiceHandlerFromEndpoint
}

func (p *PhotoTagsHandler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterPhotoTagsServiceServer(server, p)
}

func (p *PhotoTagsHandler) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		interceptor.NewCustomDescriptor((*PhotoTagsHandler).GetPhotoTags),
		interceptor.NewCustomDescriptor((*PhotoTagsHandler).AddPhotoTag),
		interceptor.NewCustomDescriptor((*PhotoTagsHandler).DeletePhotoTag),
		interceptor.NewCustomDescriptor((*PhotoTagsHandler).GetTagCategories),
	}
}
