package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

type contextKey string

const txKey contextKey = "pgx_tx"

type transactor struct {
	pool *pgxpool.Pool
}

func newTransactor(pool *pgxpool.Pool) transactor {
	return transactor{pool: pool}
}

func (t *transactor) getQueries(ctx context.Context) *photo_library.Queries {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return photo_library.New(tx)
	}
	return photo_library.New(t.pool)
}

func (t *transactor) runTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	if _, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return txFunc(ctx)
	}

	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("con.BeginTx: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey, tx)
	err = txFunc(txCtx)

	if err != nil {
		rollBackErr := tx.Rollback(ctx)
		if rollBackErr != nil {
			return fmt.Errorf("tx.Rollback: %w", rollBackErr)
		}
		return printError(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", printError(err))
	}

	return nil
}

// {"level":"error","ts":"2024-04-21T11:05:33.4534391+03:00","log":"processing_photos",
//"file":"processing/service.go:243","msg":"564621ae-fba1-401e-948d-ac3bdceea89d: processor.Processing:

//status PHOTO_GROUP: s.mergePhotoGroups: s.storage.RunTransaction: s.storage.GetGroupPhotoIDs: SQL Error: UPDATE или DELETE в таблице \"photo_groups\" нарушает ограничение внешнего ключа \"photo_groups_photos_group_id_fkey\" таблицы \"photo_groups_photos\", Detail: На ключ (id)=(257dcc3a-8378-413d-b2b0-97d434684c55)
//всё ещё есть ссылки в таблице \"photo_groups_photos\"., Where: , Code: 23503, SQLState: 23503\n"}
