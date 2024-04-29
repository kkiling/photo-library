package storage

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) SavePersonAuth(ctx context.Context, auth model.Auth) error {
	queries := r.getQueries(ctx)

	err := queries.SavePersonAuth(ctx, photo_library.SavePersonAuthParams{
		PersonID:     auth.PersonID,
		CreatedAt:    auth.CreateAt,
		UpdatedAt:    auth.UpdateAt,
		Email:        auth.Email,
		PasswordHash: auth.PasswordHash,
		Status:       photo_library.AuthStatus(auth.Status),
		Role:         photo_library.AuthRole(auth.Role),
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) EmailExists(ctx context.Context, email string) (bool, error) {
	queries := r.getQueries(ctx)

	res, err := queries.EmailExists(ctx, email)
	if err != nil {
		return false, printError(err)
	}

	return res > 0, nil
}

func (r *Adapter) GetAuth(ctx context.Context, personID uuid.UUID) (model.Auth, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetAuth(ctx, personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Auth{}, serviceerr.ErrNotFound
		}
		return model.Auth{}, printError(err)
	}

	return model.Auth{
		Base: model.Base{
			CreateAt: res.CreatedAt,
			UpdateAt: res.UpdatedAt,
		},
		PersonID:     res.PersonID,
		Email:        res.Email,
		PasswordHash: res.PasswordHash,
		Status:       model.AuthStatus(res.Status),
		Role:         model.AuthRole(res.Role),
	}, nil
}

func (r *Adapter) GetAuthByEmail(ctx context.Context, email string) (model.Auth, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetAuthByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Auth{}, serviceerr.ErrNotFound
		}
		return model.Auth{}, printError(err)
	}

	return model.Auth{
		Base: model.Base{
			CreateAt: res.CreatedAt,
			UpdateAt: res.UpdatedAt,
		},
		PersonID:     res.PersonID,
		Email:        res.Email,
		PasswordHash: res.PasswordHash,
		Status:       model.AuthStatus(res.Status),
		Role:         model.AuthRole(res.Role),
	}, nil
}

func (r *Adapter) UpdatePersonAuth(ctx context.Context, personID uuid.UUID, update model.UpdateAuth) error {
	tx := r.getTX(ctx)

	builder := sq.Update("auth").
		Where(sq.Eq{"person_id": personID}).
		Set("updated_at", update.UpdateAt)

	if update.PasswordHash.NeedUpdate {
		builder = builder.Set("password_hash", update.PasswordHash.Value)
	}
	if update.Status.NeedUpdate {
		builder = builder.Set("status", photo_library.AuthStatus(update.Status.Value))
	}
	if update.Role.NeedUpdate {
		builder = builder.Set("role", photo_library.AuthRole(update.Role.Value))
	}

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return printError(err)
	}

	return nil
}
