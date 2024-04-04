package tagsservice

import (
	"context"
	"net/http"

	"google.golang.org/grpc"

	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
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
	notFoundErr := serviceerr.NotFoundf("group not found")

	return &desc.GetTagsCategoryResponse{
		Page:    request.Page,
		PerPage: request.PerPage,
	}, handler.HandleError(notFoundErr, "GetTagsCategory")
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
