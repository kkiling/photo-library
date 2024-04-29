package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
	"time"
)

func (r *Adapter) GetLastActiveRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (model.RefreshToken, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetLastActiveRefreshToken(ctx, refreshTokenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RefreshToken{}, serviceerr.ErrNotFound
		}
		return model.RefreshToken{}, printError(err)
	}

	return model.RefreshToken{
		Base: model.Base{
			CreateAt: res.CreatedAt,
			UpdateAt: res.UpdatedAt,
		},
		ID:       res.ID,
		PersonID: res.PersonID,
		Status:   model.RefreshTokenStatus(res.Status),
	}, nil
}

func (r *Adapter) SaveRefreshToken(ctx context.Context, refreshToken model.RefreshToken) error {
	queries := r.getQueries(ctx)

	err := queries.SaveRefreshToken(ctx, photo_library.SaveRefreshTokenParams{
		ID:        refreshToken.ID,
		PersonID:  refreshToken.PersonID,
		CreatedAt: refreshToken.CreateAt,
		UpdatedAt: refreshToken.UpdateAt,
		Status:    photo_library.RefreshTokenStatus(refreshToken.Status),
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) UpdateRefreshTokenStatus(ctx context.Context, refreshTokenID uuid.UUID, status model.RefreshTokenStatus) error {
	queries := r.getQueries(ctx)

	err := queries.UpdateRefreshTokenStatus(ctx, photo_library.UpdateRefreshTokenStatusParams{
		ID:        refreshTokenID,
		UpdatedAt: time.Now(),
		Status:    photo_library.RefreshTokenStatus(status),
	})

	if err != nil {
		return printError(err)
	}

	return nil
}
