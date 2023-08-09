package service

import "context"

type Transactor interface {
	RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error
}
