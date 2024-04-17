package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Adapter struct {
	transactor
}

func NewStorageAdapter(pool *pgxpool.Pool) *Adapter {
	return &Adapter{
		transactor: newTransactor(pool),
	}
}

func (r *Adapter) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	return r.transactor.runTransaction(ctx, txFunc)
}
