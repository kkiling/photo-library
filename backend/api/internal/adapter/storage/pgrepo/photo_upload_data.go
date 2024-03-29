package pgrepo

import (
	"context"
	"errors"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PhotoRepository) SavePhotoUploadData(ctx context.Context, data entity.PhotoUploadData) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO photo_upload_data (photo_id, paths, upload_at, client_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := conn.Exec(ctx, query, data.PhotoID, data.Paths, data.UploadAt, data.ClientId)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetPhotoUploadData(ctx context.Context, photoID uuid.UUID) (*entity.PhotoUploadData, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT photo_id, paths, upload_at, client_id
		FROM photo_upload_data
		WHERE photo_id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, photoID)

	var uploadData entity.PhotoUploadData
	err := row.Scan(&uploadData.PhotoID, &uploadData.Paths, &uploadData.UploadAt, &uploadData.ClientId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &uploadData, nil
}
