package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *Service) createAuthData(ctx context.Context, personAuth model.Auth) (model.AuthDataDTO, error) {
	refreshToken := model.RefreshToken{
		Base:     model.NewBase(),
		ID:       uuid.New(),
		PersonID: personAuth.PersonID,
		Status:   model.RefreshTokenStatusActive,
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
		RefreshTokenID: refreshToken.ID,
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

func (s *Service) sendInvite(ctx context.Context, email string, role model.AuthRole) error {
	if emailExists, err := s.storage.EmailExists(ctx, email); err != nil {
		return serviceerr.MakeErr(err, "s.storage.EmailExists")
	} else if emailExists {
		return serviceerr.Conflictf("email already exists")
	}

	newPerson := model.Person{
		Base: model.NewBase(),
		ID:   uuid.New(),
	}

	newAuth := model.Auth{
		Base:         model.NewBase(),
		PersonID:     newPerson.ID,
		Email:        email,
		PasswordHash: []byte{},
		Status:       model.AuthStatusSentInvite,
		Role:         role,
	}

	err := s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		if saveErr := s.storage.SavePerson(ctxTx, newPerson); saveErr != nil {
			return fmt.Errorf("s.storage.CreatePerson: %w", saveErr)
		}
		if saveErr := s.storage.SavePersonAuth(ctxTx, newAuth); saveErr != nil {
			return fmt.Errorf("s.storage.AddPersonAuth: %w", saveErr)
		}
		return nil
	})

	if err != nil {
		return serviceerr.MakeErr(err, " s.storage.RunTransaction")
	}

	err = s.confirmCodeService.SendConfirmCode(ctx, newPerson.ID, model.ConfirmCodeTypeActivateInvite)
	if err != nil {
		return serviceerr.MakeErr(err, "s.confirmCodeService.SendConfirmCode")
	}

	return nil
}
