package auth_handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/interceptor"
	"github.com/kkiling/photo-library/backend/api/internal/service/form"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"net/http"

	"google.golang.org/grpc"

	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	methoddescriptor "github.com/kkiling/photo-library/backend/api/pkg/common/server/method_descriptor"
)

type AuthService interface {
	CheckPersonsExists(ctx context.Context) (bool, error)
	AdminInitInvite(ctx context.Context, form form.AdminInitInviteForm) error
	SendInvite(ctx context.Context, form form.SendInviteForm) error
	ActivateInvite(ctx context.Context, form form.ActivateInviteForm) error
	Login(ctx context.Context, form form.LoginForm) (model.AuthDataDTO, error)
	Registration(ctx context.Context, form form.RegisterForm) error
	ActivateRegistration(ctx context.Context, form form.ActivateRegisterForm) error
	EmailAvailable(ctx context.Context, form form.EmailAvailableForm) (bool, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token string) (model.AuthDataDTO, error)
}

type ApiTokenService interface {
	GetApiTokens(ctx context.Context) ([]model.ApiTokenDTO, error)
	CreateApiToken(ctx context.Context, form form.CreateApiToken) (string, error)
	DeleteApiToken(ctx context.Context, tokenID uuid.UUID) error
}

type AuthHandler struct {
	logger          log.Logger
	authService     AuthService
	apiTokenService ApiTokenService
}

func NewHandler(logger log.Logger, authService AuthService, apiTokenService ApiTokenService) *AuthHandler {
	return &AuthHandler{
		logger:          logger,
		authService:     authService,
		apiTokenService: apiTokenService,
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
		interceptor.NewCustomDescriptor((*AuthHandler).EmailAvailable),
		interceptor.NewCustomDescriptor((*AuthHandler).CheckPersonsExists),
		interceptor.NewCustomDescriptor((*AuthHandler).SendInvite, model.AuthRoleAdmin),
		interceptor.NewCustomDescriptor((*AuthHandler).ActivateInvite),
		interceptor.NewCustomDescriptor((*AuthHandler).Registration),
		interceptor.NewCustomDescriptor((*AuthHandler).ActivateRegistration),
		interceptor.NewCustomDescriptor((*AuthHandler).Login),
		interceptor.NewCustomDescriptor((*AuthHandler).RefreshToken),
		interceptor.NewCustomDescriptor((*AuthHandler).AdminInitInvite),
		interceptor.NewCustomDescriptor((*AuthHandler).Logout),
		interceptor.NewCustomDescriptor((*AuthHandler).GetApiTokens, model.AuthRoleAll...),
		interceptor.NewCustomDescriptor((*AuthHandler).CreateApiToken, model.AuthRoleAll...),
		interceptor.NewCustomDescriptor((*AuthHandler).DeleteApiToken, model.AuthRoleAll...),
	}
}
