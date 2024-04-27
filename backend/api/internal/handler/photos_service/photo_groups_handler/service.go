package photo_groups_handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/kkiling/photo-library/backend/api/internal/interceptor"
	"github.com/kkiling/photo-library/backend/api/internal/service/form"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type PhotoGroupsService interface {
	GetPhotoContent(ctx context.Context, fileKey string) (model.PhotoContentDTO, error)
	GetPhotoGroups(ctx context.Context, req form.GetPhotoGroups) (model.PaginatedPhotoGroupsDTO, error)
	GetPhotoGroup(ctx context.Context, groupID uuid.UUID) (model.PhotoGroupDTO, error)
	SetMainPhotoGroup(ctx context.Context, groupID, photoID uuid.UUID) error
}

type PhotoGroupsHandler struct {
	logger        log.Logger
	photosService PhotoGroupsService
}

func NewHandler(logger log.Logger, photosService PhotoGroupsService) *PhotoGroupsHandler {
	return &PhotoGroupsHandler{
		logger:        logger,
		photosService: photosService,
	}
}

func (p *PhotoGroupsHandler) RegistrationServerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/photos/", p.GetPhotoContent)
}

func (p *PhotoGroupsHandler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterPhotoGroupsServiceHandlerFromEndpoint
}

func (p *PhotoGroupsHandler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterPhotoGroupsServiceServer(server, p)
}

func (p *PhotoGroupsHandler) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		interceptor.NewCustomDescriptor((*PhotoGroupsHandler).GetPhotoGroup),
		interceptor.NewCustomDescriptor((*PhotoGroupsHandler).GetPhotoGroups),
		interceptor.NewCustomDescriptor((*PhotoGroupsHandler).SetMainPhotoGroup),
	}
}
