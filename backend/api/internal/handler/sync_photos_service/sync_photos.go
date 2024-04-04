package syncphotosservice

import (
	"context"
	"net/http"

	"google.golang.org/grpc"

	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/internal/handler/mapper"
	"github.com/kkiling/photo-library/backend/api/internal/service/syncphotos"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type SyncPhotosService interface {
	UploadPhoto(ctx context.Context, form *syncphotos.SyncPhotoRequest) (*syncphotos.SyncPhotoResponse, error)
}

type HandlerSyncPhotosService struct {
	desc.UnimplementedSyncPhotosServiceServer
	logger    log.Logger
	syncPhoto SyncPhotosService
}

func NewHandlerSyncPhotosService(logger log.Logger, syncPhoto SyncPhotosService) *HandlerSyncPhotosService {
	return &HandlerSyncPhotosService{
		logger:    logger,
		syncPhoto: syncPhoto,
	}
}

func (p *HandlerSyncPhotosService) UploadPhoto(ctx context.Context, request *desc.UploadPhotoRequest) (*desc.UploadPhotoResponse, error) {
	response, err := p.syncPhoto.UploadPhoto(ctx, mapper.UploadPhotoRequest(request))

	if err != nil {
		return nil, handler.HandleError(err, "p.syncPhoto.UploadPhoto")
	}

	return mapper.UploadPhotoResponse(response), nil
}

func (p *HandlerSyncPhotosService) RegistrationServerHandlers(mux *http.ServeMux) {

}

func (p *HandlerSyncPhotosService) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterSyncPhotosServiceHandlerFromEndpoint
}

func (p *HandlerSyncPhotosService) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterSyncPhotosServiceServer(server, p)
}

func (p *HandlerSyncPhotosService) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		handler.NewCustomDescriptor((*HandlerSyncPhotosService).UploadPhoto),
	}
}
