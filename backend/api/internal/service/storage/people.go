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

func (r *Adapter) GetPeopleCount(ctx context.Context) (int64, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPeopleCount(ctx)
	if err != nil {
		return 0, printError(err)
	}

	return res, nil
}

func (r *Adapter) SavePerson(ctx context.Context, person model.Person) error {
	queries := r.getQueries(ctx)

	err := queries.SavePerson(ctx, photo_library.SavePersonParams{
		ID:         person.ID,
		CreatedAt:  person.CreateAt,
		UpdatedAt:  person.UpdateAt,
		Firstname:  person.FirstName,
		Surname:    person.Surname,
		Patronymic: person.Patronymic,
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetPerson(ctx context.Context, personID uuid.UUID) (model.Person, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPerson(ctx, personID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Person{}, serviceerr.ErrNotFound
		}
		return model.Person{}, printError(err)
	}

	return model.Person{
		Base: model.Base{
			CreateAt: res.CreatedAt,
			UpdateAt: res.UpdatedAt,
		},
		ID:         res.ID,
		FirstName:  res.Firstname,
		Surname:    res.Surname,
		Patronymic: res.Patronymic,
	}, nil
}

func (r *Adapter) UpdatePerson(ctx context.Context, personID uuid.UUID, update model.UpdatePerson) error {
	tx := r.getTX(ctx)

	builder := sq.Update("people").
		Where(sq.Eq{"id": personID}).
		Set("updated_at", update.UpdateAt)

	if update.FirstName.NeedUpdate {
		builder = builder.Set("firstname", update.FirstName.Value)
	}
	if update.Surname.NeedUpdate {
		builder = builder.Set("surname", update.Surname.Value)
	}
	if update.Patronymic.NeedUpdate {
		builder = builder.Set("patronymic", update.Patronymic.Value)
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
