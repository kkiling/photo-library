package storage

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) SaveConfirmCode(ctx context.Context, code model.ConfirmCode) error {
	queries := r.getQueries(ctx)

	err := queries.SaveConfirmCode(ctx, photo_library.SaveConfirmCodeParams{
		Code:      code.Code,
		PersonID:  code.PersonID,
		CreatedAt: code.CreateAt,
		UpdatedAt: code.UpdateAt,
		Active:    code.Active,
		Type:      photo_library.CodeType(code.Type),
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetActiveConfirmCode(ctx context.Context, code string, confirmType model.ConfirmCodeType) (model.ConfirmCode, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetActiveConfirmCode(ctx, photo_library.GetActiveConfirmCodeParams{
		Code: code,
		Type: photo_library.CodeType(confirmType),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ConfirmCode{}, serviceerr.ErrNotFound
		}
		return model.ConfirmCode{}, printError(err)
	}

	return model.ConfirmCode{
		Base: model.Base{
			CreateAt: res.CreatedAt,
			UpdateAt: res.UpdatedAt,
		},
		Code:     res.Code,
		PersonID: res.PersonID,
		Type:     model.ConfirmCodeType(res.Type),
		Active:   res.Active,
	}, nil
}

func (r *Adapter) UpdateConfirmCode(ctx context.Context,
	personID uuid.UUID, confirmCodeType model.ConfirmCodeType, update model.UpdateConfirmCode) error {
	tx := r.getTX(ctx)

	builder := sq.Update("codes").
		Where(sq.Eq{"person_id": personID}).
		Where(sq.Eq{"type": photo_library.CodeType(confirmCodeType)}).
		Set("updated_at", update.UpdateAt)

	if update.Active.NeedUpdate {
		builder = builder.Set("active", update.Active.Value)
	}

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return printError(err)
	}

	return nil
}
