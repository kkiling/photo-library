package pgrepo

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
)

func (r *PhotoRepository) GetPhotoByHash(ctx context.Context, hash string) (*entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_name, hash, update_at, upload_at, extension, processing_status
		FROM photos
		WHERE hash = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, hash)

	var photo entity.Photo
	err := row.Scan(&photo.ID, &photo.FileName, &photo.Hash, &photo.UpdateAt, &photo.UploadAt, &photo.Extension, &photo.ProcessingStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &photo, nil
}

func (r *PhotoRepository) SavePhoto(ctx context.Context, photo entity.Photo) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO photos (id, file_name, hash, update_at, upload_at, extension, processing_status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := conn.Exec(ctx, query, photo.ID, photo.FileName, photo.Hash, photo.UpdateAt, photo.UploadAt, photo.Extension, photo.ProcessingStatus)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetPhotoById(ctx context.Context, id uuid.UUID) (*entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_name, hash, update_at, upload_at, extension, processing_status
		FROM photos
		WHERE id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, id)

	var photo entity.Photo
	err := row.Scan(&photo.ID, &photo.FileName, &photo.Hash, &photo.UpdateAt, &photo.UploadAt, &photo.Extension, &photo.ProcessingStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &photo, nil
}

func (r *PhotoRepository) GetPaginatedPhotos(ctx context.Context, params entity.PhotoSelectParams, filter *entity.PhotoFilter) ([]entity.Photo, error) {
	conn := r.getConn(ctx)

	builder := sq.
		Select("id", "file_name", "hash", "update_at", "upload_at", "extension", "processing_status").
		From("photos").
		Offset(uint64(params.Offset)).
		Limit(uint64(params.Limit)).
		PlaceholderFormat(sq.Dollar)

	if filter != nil {
		if len(filter.ProcessingStatusIn) > 0 {
			builder = builder.Where(sq.Eq{"processing_status": filter.ProcessingStatusIn})
		}
	}

	if params.SortOrder != entity.PhotoSortOrderNone {
		// Добавляем сортировку
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("builder.ToSql: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result = make([]entity.Photo, 0, params.Limit)
	for rows.Next() {
		var photo entity.Photo

		errScan := rows.Scan(&photo.ID, &photo.FileName, &photo.Hash, &photo.UpdateAt, &photo.UploadAt, &photo.Extension, &photo.ProcessingStatus)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, errScan
		}

		if errScan != nil {
			return nil, printError(err)
		}
		result = append(result, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}

func (r *PhotoRepository) GetPhotosCount(ctx context.Context, filter *entity.PhotoFilter) (int64, error) {
	conn := r.getConn(ctx)

	var counter int64

	builder := sq.
		Select("count(1)").
		From("photos").
		PlaceholderFormat(sq.Dollar)

	if filter != nil {
		if len(filter.ProcessingStatusIn) > 0 {
			builder = builder.Where(sq.Eq{"processing_status": filter.ProcessingStatusIn})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("builder.ToSql: %w", err)
	}

	err = conn.QueryRow(ctx, query, args...).Scan(&counter)
	if err != nil {
		return 0, printError(err)
	}

	return counter, nil
}
func (r *PhotoRepository) UpdatePhotosProcessingStatus(ctx context.Context, id uuid.UUID, newProcessingStatus string) error {
	conn := r.getConn(ctx)

	builder := sq.
		Update("photos").
		Set("processing_status", newProcessingStatus).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("builder.ToSql: %w", err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return printError(err)
	}

	return nil
}
