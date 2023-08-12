package pgrepo

import (
	"context"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
)

func (r *PhotoRepository) SaveUploadPhotoData(ctx context.Context, data entity.UploadPhotoData) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO upload_photo_data (photo_id, paths, upload_at, client_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := conn.Exec(ctx, query, data.PhotoID, data.Paths, data.UploadAt, data.ClientId)
	if err != nil {
		return err
	}

	return nil
}
