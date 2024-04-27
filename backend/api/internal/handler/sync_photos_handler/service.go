package sync_photo_handler

import (
	"context"
	"net/http"

	"google.golang.org/grpc"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type SyncPhotosService interface {
	UploadPhoto(ctx context.Context, form *model.SyncPhotoRequest) (*model.SyncPhotoResponse, error)
}

type SyncPhotoHandler struct {
	desc.UnimplementedSyncPhotosServiceServer
	logger    log.Logger
	syncPhoto SyncPhotosService
}

func NewHandler(logger log.Logger, syncPhoto SyncPhotosService) *SyncPhotoHandler {
	return &SyncPhotoHandler{
		logger:    logger,
		syncPhoto: syncPhoto,
	}
}

func (p *SyncPhotoHandler) RegistrationServerHandlers(mux *http.ServeMux) {

}

func (p *SyncPhotoHandler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterSyncPhotosServiceHandlerFromEndpoint
}

func (p *SyncPhotoHandler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterSyncPhotosServiceServer(server, p)
}

func (p *SyncPhotoHandler) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		// server2.NewCustomDescriptor((*SyncPhotoHandler).UploadPhoto),
	}
}
