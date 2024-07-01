package utils

import (
	"context"
	"errors"
	"github.com/kkiling/photo-library/backend/api/internal/ctxutils"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/samber/lo"
)

func GetSession(ctx context.Context, roles ...model.AuthRole) (model.Session, error) {
	session, err := ctxutils.Get[model.Session](ctx, ctxutils.Session)
	switch {
	case err == nil:
	case errors.Is(err, ctxutils.ErrNotFound):
		return model.Session{}, serviceerr.PermissionDeniedf("permission denied")
	default:
		return model.Session{}, serviceerr.MakeErr(err, "ctxutils.Get")
	}

	if len(roles) > 0 {
		if !lo.Contains(roles, session.Role) {
			return model.Session{}, serviceerr.PermissionDeniedf("permission denied")
		}
	}

	return session, nil
}
