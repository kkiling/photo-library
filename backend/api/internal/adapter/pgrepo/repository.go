package pgrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
)

type PhotoRepository struct {
	Transactor
}

func NewPhotoRepository(pool *pgxpool.Pool) *PhotoRepository {
	return &PhotoRepository{
		Transactor: NewTransactor(pool),
	}
}

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
		return nil, err
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
		return err
	}

	return nil
}

func (r *PhotoRepository) SaveUploadPhotoData(ctx context.Context, data entity.UploadPhotoData) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO upload_photo_data (id, photo_id, paths, upload_at, client_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := conn.Exec(ctx, query, data.ID, data.PhotoID, data.Paths, data.UploadAt, data.ClientId)
	if err != nil {
		return err
	}

	return nil
}
