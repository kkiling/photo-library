package tagsservice

import (
	"context"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
	"google.golang.org/grpc"
	"net/http"
)

type HandlerTagsService struct {
	desc.UnimplementedTagsServiceServer
	logger log.Logger
}

func NewHandlerTagsService(logger log.Logger) *HandlerTagsService {
	return &HandlerTagsService{
		logger: logger,
	}
}

func (p *HandlerTagsService) GetTagsCategory(ctx context.Context, request *desc.GetTagsCategoryRequest) (*desc.GetTagsCategoryResponse, error) {
	return &desc.GetTagsCategoryResponse{
		Page:    request.Page,
		PerPage: request.PerPage,
	}, nil
}

func (p *HandlerTagsService) RegistrationServerHandlers(mux *http.ServeMux) {

}

func (p *HandlerTagsService) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterTagsServiceHandlerFromEndpoint
}

func (p *HandlerTagsService) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterTagsServiceServer(server, p)
}

func (p *HandlerTagsService) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		handler.NewCustomDescriptor((*HandlerTagsService).GetTagsCategory),
	}
}
