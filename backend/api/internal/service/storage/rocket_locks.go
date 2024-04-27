package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) RocketLock(ctx context.Context, key string, ttl time.Duration) (model.RocketLockID, error) {
	queries := r.getQueries(ctx)
	res, err := queries.RocketLock(ctx, photo_library.RocketLockParams{
		Key: key,
		Interval: pgtype.Interval{
			Microseconds: ttl.Microseconds(),
			Valid:        true,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RocketLockID{}, serviceerr.ErrAlreadyLocked
		}
		return model.RocketLockID{}, printError(err)
	}

	return model.RocketLockID{
		Key: key,
		Ts:  uint64(res),
	}, nil
}

func (r *Adapter) RocketLockDelete(ctx context.Context, lock model.RocketLockID) error {
	queries := r.getQueries(ctx)
	err := queries.RocketLockDelete(ctx, photo_library.RocketLockDeleteParams{
		Key: lock.Key,
		Ts:  lo.ToPtr(int64(lock.Ts)),
	})

	if err != nil {
		return printError(err)
	}
	return nil
}
