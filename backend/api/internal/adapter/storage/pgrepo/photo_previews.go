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

func (r *PhotoRepository) CreatePhotoPreview(ctx context.Context, preview *entity.PhotoPreview) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO photo_previews (id, photo_id, file_name, width_pixel, height_pixel, size_pixel)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := conn.Exec(ctx, query, preview.ID, preview.PhotoID, preview.FileName, preview.WidthPixel, preview.HeightPixel, preview.SizePixel)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetPhotoPreviews(ctx context.Context, photoID uuid.UUID) ([]entity.PhotoPreview, error) {
	conn := r.getConn(ctx)

	builder := sq.
		Select("id", "photo_id", "file_name", "width_pixel", "height_pixel", "size_pixel").
		From("photo_previews").
		Where(sq.Eq{"photo_id": photoID}).
		OrderBy("size_pixel").
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

	var result = make([]entity.PhotoPreview, 0)
	for rows.Next() {
		var preview entity.PhotoPreview

		errScan := rows.Scan(&preview.ID, &preview.PhotoID, &preview.FileName, &preview.WidthPixel, &preview.HeightPixel, &preview.SizePixel)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, printError(err)
		}

		result = append(result, preview)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}
