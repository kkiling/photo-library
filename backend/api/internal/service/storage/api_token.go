package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
	"github.com/samber/lo"
)

func mapApiToken(item photo_library.ApiToken) model.ApiToken {
	return model.ApiToken{
		Base: model.Base{
			CreateAt: item.CreatedAt,
			UpdateAt: item.UpdatedAt,
		},
		ID:        item.ID,
		PersonID:  item.PersonID,
		Caption:   item.Caption,
		Token:     item.Token,
		Type:      model.ApiTokenType(item.Type),
		ExpiredAt: toTimePtr(item.ExpiredAt),
	}
}

func (r *Adapter) GetApiTokens(ctx context.Context, personID uuid.UUID) ([]model.ApiToken, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetApiTokens(ctx, personID)
	if err != nil {
		return nil, printError(err)
	}

	return lo.Map(res, func(item photo_library.ApiToken, index int) model.ApiToken {
		return mapApiToken(item)
	}), nil
}

func (r *Adapter) SaveApiToken(ctx context.Context, apiToken model.ApiToken) error {
	queries := r.getQueries(ctx)

	err := queries.SaveApiToken(ctx, photo_library.SaveApiTokenParams{
		ID:        apiToken.ID,
		PersonID:  apiToken.PersonID,
		Caption:   apiToken.Caption,
		Token:     apiToken.Token,
		CreatedAt: apiToken.CreateAt,
		UpdatedAt: apiToken.UpdateAt,
		ExpiredAt: toTimestamptz(apiToken.ExpiredAt),
		Type:      photo_library.ApiTokenType(apiToken.Type),
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) DeleteApiToken(ctx context.Context, personID, tokenID uuid.UUID) error {
	queries := r.getQueries(ctx)

	_, err := queries.DeleteApiToken(ctx, photo_library.DeleteApiTokenParams{
		ID:       tokenID,
		PersonID: personID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return serviceerr.ErrNotFound
		}
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetApiToken(ctx context.Context, token string) (model.ApiToken, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetApiToken(ctx, token)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ApiToken{}, serviceerr.ErrNotFound
		}
		return model.ApiToken{}, printError(err)
	}

	return mapApiToken(res), nil
}
