package pgrepo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
)

func (r *PhotoRepository) GetPhotoByHash(ctx context.Context, hash string) (*entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_path, hash, update_at, upload_at, extension
		FROM photos
		WHERE hash = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, hash)

	var photo entity.Photo
	err := row.Scan(&photo.ID, &photo.FilePath, &photo.Hash, &photo.UpdateAt, &photo.UploadAt, &photo.Extension)
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
		INSERT INTO photos (id, file_path, hash, update_at, upload_at, extension)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := conn.Exec(ctx, query, photo.ID, photo.FilePath, photo.Hash, photo.UpdateAt, photo.UploadAt, photo.Extension)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetPhotoById(ctx context.Context, id uuid.UUID) (*entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_path, hash, update_at, upload_at, extension
		FROM photos
		WHERE id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, id)

	var photo entity.Photo
	err := row.Scan(&photo.ID, &photo.FilePath, &photo.Hash, &photo.UpdateAt, &photo.UploadAt, &photo.Extension)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &photo, nil
}

func (r *PhotoRepository) GetPaginatedPhotos(ctx context.Context, offset int64, limit int) ([]entity.Photo, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT id, file_path, hash, update_at, upload_at, extension
		FROM photos
		OFFSET $1
		LIMIT $2
	`

	rows, err := conn.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result []entity.Photo
	for rows.Next() {
		var photo entity.Photo

		errScan := rows.Scan(&photo.ID, &photo.FilePath, &photo.Hash, &photo.UpdateAt, &photo.UploadAt, &photo.Extension)
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

func (r *PhotoRepository) GetPhotosCount(ctx context.Context) (int64, error) {
	conn := r.getConn(ctx)

	var counter int64

	const query = `
		SELECT count(*)
		FROM photos
	`

	err := conn.QueryRow(ctx, query).Scan(&counter)
	if err != nil {
		return 0, printError(err)
	}

	return counter, nil
}
