package utils

import (
	"context"
	"errors"
	"github.com/kkiling/photo-library/backend/api/internal/ctxutils"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func GetApiToken(ctx context.Context) (model.ApiToken, error) {
	session, err := ctxutils.Get[model.ApiToken](ctx, ctxutils.ApiToken)
	switch {
	case err == nil:
	case errors.Is(err, ctxutils.ErrNotFound):
		return model.ApiToken{}, serviceerr.PermissionDeniedf("permission denied")
	default:
		return model.ApiToken{}, serviceerr.MakeErr(err, "ctxutils.Get")
	}
	return session, nil
}
