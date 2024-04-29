package auth_handler

import (
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/service/form"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func mapToModelRole(role desc.AuthRole) (model.AuthRole, error) {
	switch role {
	case desc.AuthRole_AUTH_ROLE_ADMIN:
		return model.AuthRoleAdmin, nil
	case desc.AuthRole_AUTH_ROLE_USER:
		return model.AuthRoleUser, nil
	default:
		return "", fmt.Errorf("invalid auth role: %s", role)
	}
}

func mapApiTokenType(apiTokenType model.ApiTokenType) (desc.ApiTokenType, error) {
	switch apiTokenType {
	case model.ApiTokenSyncPhoto:
		return desc.ApiTokenType_API_TOKEN_TYPE_SYNC_PHOTO, nil
	default:
		return desc.ApiTokenType_API_TOKEN_TYPE_UNKNOWN, fmt.Errorf("invalid api token type: %s", apiTokenType)
	}
}

func mapAuthData(res model.AuthDataDTO) *desc.AuthData {
	return &desc.AuthData{
		PersonId:               res.PersonID.String(),
		Email:                  res.Email,
		AccessToken:            res.AccessToken,
		AccessTokenExpiration:  timestamppb.New(res.AccessTokenExpiration),
		RefreshToken:           res.RefreshToken,
		RefreshTokenExpiration: timestamppb.New(res.RefreshTokenExpiration),
	}
}

func mapApiTokens(res []model.ApiTokenDTO) ([]*desc.ApiToken, error) {
	var result = make([]*desc.ApiToken, 0, len(res))

	for _, item := range res {
		apiTokenType := desc.ApiTokenType_API_TOKEN_TYPE_UNKNOWN
		switch item.Type {
		case model.ApiTokenSyncPhoto:
			apiTokenType = desc.ApiTokenType_API_TOKEN_TYPE_SYNC_PHOTO
		default:
			return nil, fmt.Errorf("invalid api token type: %s", apiTokenType)
		}

		var expiredAt *timestamppb.Timestamp
		if item.ExpiredAt != nil {
			expiredAt = timestamppb.New(*item.ExpiredAt)
		}
		token := desc.ApiToken{
			Id:        item.ID.String(),
			Caption:   item.Caption,
			Type:      apiTokenType,
			ExpiredAt: expiredAt,
		}

		result = append(result, &token)
	}

	return result, nil
}

func mapCreateApiToken(request *desc.CreateApiTokenRequest) (form.CreateApiToken, error) {
	var timeDuration *time.Duration
	if request.TimeDuration != nil {
		tt, err := time.ParseDuration(request.GetTimeDuration())
		if err != nil {
			return form.CreateApiToken{}, fmt.Errorf("invalid time duration: %s", request.GetTimeDuration())
		}
		timeDuration = &tt
	}

	var apiTokenType model.ApiTokenType
	switch request.Type {
	case desc.ApiTokenType_API_TOKEN_TYPE_SYNC_PHOTO:
		apiTokenType = model.ApiTokenSyncPhoto
	default:
		return form.CreateApiToken{}, fmt.Errorf("invalid api token type: %s", apiTokenType)
	}

	return form.CreateApiToken{
		Caption:      request.Caption,
		Type:         apiTokenType,
		TimeDuration: timeDuration,
	}, nil
}
