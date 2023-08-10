package pgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const txKey contextKey = "pgx_tx"

type Transactor struct {
	pool *pgxpool.Pool
}

type ConnectContext interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewTransactor(pool *pgxpool.Pool) Transactor {
	return Transactor{pool: pool}
}

func (t *Transactor) getConn(ctx context.Context) ConnectContext {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return t.pool
}

func (t *Transactor) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {

	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("con.BeginTx: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey, tx)
	err = txFunc(txCtx)

	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("tx.Rollback: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
