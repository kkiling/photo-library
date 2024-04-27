package auth_handler

import (
	"context"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p *AuthHandler) EmailAvailable(ctx context.Context, request *desc.EmailAvailableRequest) (*desc.EmailAvailableResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) CheckPersonsExists(ctx context.Context, empty *emptypb.Empty) (*desc.CheckPersonsExistsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) SendInvite(ctx context.Context, request *desc.SendInviteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) ActivateInvite(ctx context.Context, request *desc.ActivateInviteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) Registration(ctx context.Context, request *desc.RegistrationRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) ActivateRegistration(ctx context.Context, request *desc.ActivateRegistrationRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) Logout(ctx context.Context, request *desc.LoginRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) Login(ctx context.Context, request *desc.LoginRequest) (*desc.AuthData, error) {
	//TODO implement me
	panic("implement me")
}

func (p *AuthHandler) RefreshToken(ctx context.Context, request *desc.RefreshTokenRequest) (*desc.AuthData, error) {
	//TODO implement me
	panic("implement me")
}
