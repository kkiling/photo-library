package auth_handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/internal/service/form"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p *AuthHandler) EmailAvailable(ctx context.Context, request *desc.EmailAvailableRequest) (*desc.EmailAvailableResponse, error) {
	available, err := p.authService.EmailAvailable(ctx, form.EmailAvailableForm{
		Email: request.Email,
	})
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.EmailAvailable")
	}

	return &desc.EmailAvailableResponse{
		Available: available,
	}, nil
}

func (p *AuthHandler) CheckPersonsExists(ctx context.Context, empty *emptypb.Empty) (*desc.CheckPersonsExistsResponse, error) {
	exists, err := p.authService.CheckPersonsExists(ctx)
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.CheckPersonsExists")
	}

	return &desc.CheckPersonsExistsResponse{
		Exists: exists,
	}, nil
}

func (p *AuthHandler) AdminInitInvite(ctx context.Context, request *desc.AdminInitInviteRequest) (*emptypb.Empty, error) {
	err := p.authService.AdminInitInvite(ctx, form.AdminInitInviteForm{
		Email: request.Email,
	})
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.AdminInitInvite")
	}
	return &emptypb.Empty{}, nil
}

func (p *AuthHandler) SendInvite(ctx context.Context, request *desc.SendInviteRequest) (*emptypb.Empty, error) {
	role, err := mapToModelRole(request.Role)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("mapToModelRole: %w", err))
	}

	err = p.authService.SendInvite(ctx, form.SendInviteForm{
		Email: request.Email,
		Role:  role,
	})
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.SendInvite")
	}

	return &emptypb.Empty{}, nil
}

func (p *AuthHandler) ActivateInvite(ctx context.Context, request *desc.ActivateInviteRequest) (*emptypb.Empty, error) {
	err := p.authService.ActivateInvite(ctx, form.ActivateInviteForm{
		FirstName:   request.Firstname,
		Surname:     request.Surname,
		Patronymic:  request.Patronymic,
		CodeConfirm: request.CodeConfirm,
		Password:    request.Password,
	})
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.ActivateInvite")
	}
	return &emptypb.Empty{}, nil
}

func (p *AuthHandler) Registration(ctx context.Context, request *desc.RegistrationRequest) (*emptypb.Empty, error) {
	err := p.authService.Registration(ctx, form.RegisterForm{
		FirstName:  request.Firstname,
		Surname:    request.Surname,
		Patronymic: request.Patronymic,
		Email:      request.Email,
		Password:   request.Password,
	})
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.Registration")
	}
	return &emptypb.Empty{}, nil
}

func (p *AuthHandler) ActivateRegistration(ctx context.Context, request *desc.ActivateRegistrationRequest) (*emptypb.Empty, error) {
	err := p.authService.ActivateRegistration(ctx, form.ActivateRegisterForm{
		CodeConfirm: request.CodeConfirm,
	})
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.ActivateRegistration")
	}
	return &emptypb.Empty{}, nil
}

func (p *AuthHandler) Login(ctx context.Context, request *desc.LoginRequest) (*desc.AuthData, error) {
	res, err := p.authService.Login(ctx, form.LoginForm{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.Login")
	}
	return mapAuthData(res), nil
}

func (p *AuthHandler) RefreshToken(ctx context.Context, request *desc.RefreshTokenRequest) (*desc.AuthData, error) {
	res, err := p.authService.RefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.RefreshToken")
	}
	return mapAuthData(res), nil
}

func (p *AuthHandler) Logout(ctx context.Context, request *desc.LogoutRequest) (*emptypb.Empty, error) {
	err := p.authService.Logout(ctx, request.RefreshToken)
	if err != nil {
		return nil, handler.HandleError(err, "p.authService.Logout")
	}
	return &emptypb.Empty{}, nil
}

func (p *AuthHandler) GetApiTokens(ctx context.Context, request *desc.GetApiTokensRequest) (*desc.GetApiTokensResponse, error) {
	res, err := p.apiTokenService.GetApiTokens(ctx)
	if err != nil {
		return nil, handler.HandleError(err, "p.apiTokenService.GetApiTokens")
	}

	items, err := mapApiTokens(res)
	if err != nil {
		return nil, handler.HandleError(err, "mapApiTokens")
	}

	return &desc.GetApiTokensResponse{
		Items: items,
	}, nil
}

func (p *AuthHandler) CreateApiToken(ctx context.Context, request *desc.CreateApiTokenRequest) (*desc.CreateApiTokenResponse, error) {
	formReq, err := mapCreateApiToken(request)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("mapCreateApiToken: %w", err))
	}

	token, err := p.apiTokenService.CreateApiToken(ctx, formReq)
	if err != nil {
		return nil, handler.HandleError(err, "p.apiTokenService.CreateApiToken")
	}
	return &desc.CreateApiTokenResponse{
		Token: token,
	}, nil
}

func (p *AuthHandler) DeleteApiToken(ctx context.Context, request *desc.DeleteApiTokenRequest) (*desc.DeleteApiTokenResponse, error) {
	tokenID, err := uuid.Parse(request.TokenId)
	if err != nil {
		return nil, server.ErrInvalidArgument(fmt.Errorf("tokenID is invalid: %w", err))
	}

	if err = p.apiTokenService.DeleteApiToken(ctx, tokenID); err != nil {
		return nil, handler.HandleError(err, "p.apiTokenService.DeleteApiToken")
	}

	return &desc.DeleteApiTokenResponse{}, nil
}
