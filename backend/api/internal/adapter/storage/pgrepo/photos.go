package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PhotoRepository) GetPhotoByHash(ctx context.Context, hash string) (*entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_name, hash, update_at, extension
		FROM photos
		WHERE hash = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, hash)

	var photo entity.Photo
	err := row.Scan(&photo.ID, &photo.FileName, &photo.Hash, &photo.UpdateAt, &photo.Extension)
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
		INSERT INTO photos (id, file_name, hash, update_at, extension, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := conn.Exec(ctx, query, photo.ID, photo.FileName, photo.Hash, photo.UpdateAt, photo.Extension, entity.NewPhotoStatus)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) MakeNotValidPhoto(ctx context.Context, photoID uuid.UUID, error string) error {
	conn := r.getConn(ctx)

	builder := sq.
		Update("photos").
		Set("status", entity.NotValidStatus).
		Set("error", error).
		Where(sq.Eq{"id": photoID}).
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

func (r *PhotoRepository) GetPhotoById(ctx context.Context, id uuid.UUID) (*entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_name, hash, update_at, extension
		FROM photos
		WHERE id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, id)

	var photo entity.Photo
	err := row.Scan(&photo.ID, &photo.FileName, &photo.Hash, &photo.UpdateAt, &photo.Extension)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &photo, nil
}

func (r *PhotoRepository) GetPhotoByFilename(ctx context.Context, fileName string) (*entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_name, hash, update_at, extension
		FROM photos
		WHERE file_name = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, fileName)

	var photo entity.Photo
	err := row.Scan(&photo.ID, &photo.FileName, &photo.Hash, &photo.UpdateAt, &photo.Extension)
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
		Select("id", "file_name", "hash", "update_at", "extension").
		From("photos").
		Where(sq.Eq{"status": entity.NewPhotoStatus}).
		Offset(params.Offset).
		Limit(params.Limit).
		PlaceholderFormat(sq.Dollar)

	if filter != nil {
		/*if len(filter.ProcessingStatusIn) > 0 {
			builder = builder.Where(sq.Eq{"processing_status": filter.ProcessingStatusIn})
		}*/
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

		errScan := rows.Scan(&photo.ID, &photo.FileName, &photo.Hash, &photo.UpdateAt, &photo.Extension)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, printError(err)
		}

		result = append(result, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}

func (r *PhotoRepository) GetPhotosCount(ctx context.Context, filter *entity.PhotoFilter) (uint64, error) {
	conn := r.getConn(ctx)

	var counter uint64

	builder := sq.
		Select("count(1)").
		From("photos").
		PlaceholderFormat(sq.Dollar)

	if filter != nil {
		/*if len(filter.ProcessingStatusIn) > 0 {
			builder = builder.Where(sq.Eq{"processing_status": filter.ProcessingStatusIn})
		}*/
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

func (r *PhotoRepository) AddPhotosProcessingStatus(ctx context.Context, photoID uuid.UUID, status string, success bool) error {
	conn := r.getConn(ctx)

	builder := sq.Insert("photo_processing_statuses").
		SetMap(map[string]interface{}{
			"photo_id":     photoID,
			"processed_at": time.Now(),
			"status":       status,
			"success":      success,
		}).
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

func (r *PhotoRepository) GetUnprocessedPhotoIDs(ctx context.Context, lastProcessingStatus string, limit uint64) ([]uuid.UUID, error) {
	conn := r.getConn(ctx)

	builder := sq.
		Select("p.id").
		From("photos p").
		LeftJoin(fmt.Sprintf("photo_processing_statuses ps ON p.id = ps.photo_id AND ps.status = '%s'", lastProcessingStatus)).
		Where("ps.photo_id IS NULL").
		Where(sq.Eq{"p.status": entity.NewPhotoStatus}).
		Limit(limit).
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

	var result = make([]uuid.UUID, 0, limit)
	for rows.Next() {
		var photoID uuid.UUID

		errScan := rows.Scan(&photoID)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, printError(err)
		}

		result = append(result, photoID)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}

func (r *PhotoRepository) GetPhotoProcessingStatuses(ctx context.Context, photoID uuid.UUID) ([]string, error) {
	conn := r.getConn(ctx)

	builder := sq.
		Select("status").
		From("photo_processing_statuses").
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

	var result = make([]string, 0)
	for rows.Next() {
		var status string

		errScan := rows.Scan(&status)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, printError(err)
		}

		result = append(result, status)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}
