package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
)

var ErrAlreadyLocked = fmt.Errorf("rocket Key already locked")

func (r *PhotoRepository) RocketLock(ctx context.Context, key string, ttl time.Duration) (*entity.RocketLockID, error) {
	if ttl < time.Millisecond {
		return nil, errors.New("ttl must be >= 1 ms")
	}

	conn := r.getConn(ctx)

	interval := fmt.Sprintf("'%d milliseconds'", ttl.Milliseconds())
	var query = fmt.Sprintf(`INSERT INTO rocket_locks (key, locked_until) VALUES ($1, now() + INTERVAL %s)
			  ON CONFLICT (key) DO UPDATE SET locked_until = (now() + INTERVAL %s) WHERE rocket_locks.locked_until < now()
			  RETURNING floor(extract(epoch from locked_until))`, interval, interval)

	row := conn.QueryRow(ctx, query, key)

	var ts uint64
	err := row.Scan(&ts)
	switch {
	case err == nil:
	case errors.Is(err, pgx.ErrNoRows):
		return nil, ErrAlreadyLocked
	default:
		return nil, printError(err)
	}

	return &entity.RocketLockID{Key: key, Ts: ts}, nil
}

func (r *PhotoRepository) RocketLockDelete(ctx context.Context, lockID *entity.RocketLockID) error {
	conn := r.getConn(ctx)

	const query = `DELETE FROM rocket_locks where Key = $1 AND floor(extract(epoch from locked_until)) = $2`

	_, err := conn.Exec(ctx, query, lockID.Key, lockID.Ts)
	if err != nil {
		return printError(err)
	}

	return nil
}
