package ctxutils

import (
	"context"

	"github.com/pkg/errors"
)

const (
	Session  = "session"
	ApiToken = "apiToken"
)

var (
	ErrNotFound = errors.Errorf("context not contains key")
)

func Set(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func Get[T any](ctx context.Context, key string) (T, error) {
	v, ok := ctx.Value(key).(T)
	if !ok {
		return v, ErrNotFound
	}

	return v, nil
}
