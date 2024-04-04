package pgrepo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
)

func (r *PhotoRepository) GetTagCategory(ctx context.Context, categoryID uuid.UUID) (*entity.TagCategory, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, type, color
		FROM tags_category
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
		FROM tags_category
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
		INSERT INTO tags_category (id, type, color)
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
		FROM tags
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

func (r *PhotoRepository) GetTags(ctx context.Context, photoID uuid.UUID) ([]entity.Tag, error) {
	conn := r.getConn(ctx)

	builder := sq.
		Select("id", "category_id", "photo_id", "name").
		From("tags").
		Where(sq.Eq{"photo_id": photoID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("builder.ToSql: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result = make([]entity.Tag, 0)
	for rows.Next() {
		var tag entity.Tag

		errScan := rows.Scan(&tag.ID, &tag.CategoryID, &tag.PhotoID, &tag.Name)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, printError(err)
		}

		result = append(result, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}

func (r *PhotoRepository) SaveTag(ctx context.Context, tag entity.Tag) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO tags (id, category_id, photo_id, name)
		VALUES ($1, $2, $3, $4)
	`

	_, err := conn.Exec(ctx, query, tag.ID, tag.CategoryID, tag.PhotoID, tag.Name)
	if err != nil {
		return printError(err)
	}

	return nil
}
