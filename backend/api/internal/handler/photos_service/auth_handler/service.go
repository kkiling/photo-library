package auth_handler

import (
	"net/http"

	"google.golang.org/grpc"

	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type AuthService interface {
}

type AuthHandler struct {
	logger      log.Logger
	authService AuthService
}

func NewHandler(logger log.Logger, authService AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (p *AuthHandler) RegistrationServerHandlers(mux *http.ServeMux) {
}

func (p *AuthHandler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterAuthServiceHandlerFromEndpoint
}

func (p *AuthHandler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterAuthServiceServer(server, p)
}

func (p *AuthHandler) GetMethodDescriptors() []methoddescriptor.Descriptor {
	return []methoddescriptor.Descriptor{
		//interceptor.NewCustomDescriptor((*PhotoGroupsHandler).GetPhotoGroup),
		//interceptor.NewCustomDescriptor((*PhotoGroupsHandler).GetPhotoGroups),
		//interceptor.NewCustomDescriptor((*PhotoGroupsHandler).SetMainPhotoGroup),
	}
}
