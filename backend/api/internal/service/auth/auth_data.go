package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *Service) createAuthData(ctx context.Context, personAuth model.Auth) (model.AuthDataDTO, error) {
	refreshToken := model.RefreshToken{
		Base:           model.NewBase(),
		RefreshTokenID: uuid.New(),
		PersonID:       personAuth.PersonID,
		Status:         model.RefreshTokenStatusActive,
	}

	if err := s.storage.SaveRefreshToken(ctx, refreshToken); err != nil {
		return model.AuthDataDTO{}, serviceerr.MakeErr(err, "s.storage.SaveOrCreateRefreshSession")
	}

	session := model.Session{
		PersonID: personAuth.PersonID,
		Role:     personAuth.Role,
	}

	access, err := s.sessionService.CreateTokenBySession(session)
	if err != nil {
		return model.AuthDataDTO{}, serviceerr.MakeErr(err, "s.sessionService.CreateTokenBySession")
	}

	refreshSession := model.RefreshSession{
		RefreshTokenID: refreshToken.RefreshTokenID,
		PersonID:       refreshToken.PersonID,
	}

	refresh, err := s.sessionService.CreateTokenByRefresh(refreshSession)
	if err != nil {
		return model.AuthDataDTO{}, serviceerr.MakeErr(err, "s.sessionService.CreateTokenByRefresh")
	}

	return model.AuthDataDTO{
		PersonID:               personAuth.PersonID,
		Email:                  personAuth.Email,
		AccessToken:            access.Token,
		AccessTokenExpiration:  access.ExpiresAt,
		RefreshToken:           refresh.Token,
		RefreshTokenExpiration: refresh.ExpiresAt,
	}, nil
}
