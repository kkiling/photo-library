package pgrepo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
)

func (r *PhotoRepository) GetTagCategory(ctx context.Context, categoryID uuid.UUID) (*entity.TagCategory, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, type, color
		FROM tag_category
		WHERE id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, categoryID)

	var category entity.TagCategory
	err := row.Scan(&category.ID, &category.Type, &category.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &category, nil
}

func (r *PhotoRepository) GetTagCategoryByType(ctx context.Context, typeCategory string) (*entity.TagCategory, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, type, color
		FROM tag_category
		WHERE type = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, typeCategory)

	var category entity.TagCategory
	err := row.Scan(&category.ID, &category.Type, &category.Color)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &category, nil
}

func (r *PhotoRepository) SaveTagCategory(ctx context.Context, category entity.TagCategory) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO tag_category (id, type, color)
		VALUES ($1, $2, $3)
	`

	_, err := conn.Exec(ctx, query, category.ID, category.Type, category.Color)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetTagByName(ctx context.Context, photoID uuid.UUID, name string) (*entity.Tag, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, category_id, photo_id, name
		FROM tag
		WHERE photo_id = $1 AND name = $2
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, photoID, name)

	var tag entity.Tag
	err := row.Scan(&tag.ID, &tag.CategoryID, &tag.PhotoID, &tag.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &tag, nil
}

func (r *PhotoRepository) SaveTag(ctx context.Context, tag entity.Tag) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO tag (id, category_id, photo_id, name)
		VALUES ($1, $2, $3, $4)
	`

	_, err := conn.Exec(ctx, query, tag.ID, tag.CategoryID, tag.PhotoID, tag.Name)
	if err != nil {
		return printError(err)
	}

	return nil
}
