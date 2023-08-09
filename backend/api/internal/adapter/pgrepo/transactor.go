package pgrepo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type contextKey string

const txKey contextKey = "sql_tx"
const txIDKey contextKey = "tx_id"

type Transactor struct {
	conn *pgx.Conn
}

func NewTransactor(c *pgx.Conn) Transactor {
	return Transactor{conn: c}
}

func (t *Transactor) GetConn(ctx context.Context) *pgx.Conn {
	conn, ok := ctx.Value(txKey).(*pgx.Conn)
	if !ok {
		return t.conn
	}
	return conn
}

func (t *Transactor) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	ctxTx := context.WithValue(ctx, txIDKey, uuid.NewString())

	tx, err := t.conn.BeginTx(ctxTx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("con.BeginTx: %w", err)
	}

	err = txFunc(ctxTx)
	if err != nil {
		err = tx.Rollback(ctxTx)
		if err != nil {
			return fmt.Errorf("tx.Rollback: %w", err)
		}
	}

	if err = tx.Commit(ctxTx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
