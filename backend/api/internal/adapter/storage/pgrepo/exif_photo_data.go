package pgrepo

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
)

func (r *PhotoRepository) SaveExif(ctx context.Context, data *entity.ExifPhotoData) error {
	conn := r.getConn(ctx)

	query, args, err := sq.
		Insert("exif_photo_data").
		SetMap(map[string]interface{}{
			"photo_id": data.PhotoID,
			"data":     data.Data,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return printError(err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetExif(ctx context.Context, photoId uuid.UUID) (*entity.ExifPhotoData, error) {
	conn := r.getConn(ctx)

	const query = `SELECT data FROM exif_photo_data WHERE photo_id = $1`
	var row = conn.QueryRow(ctx, query, photoId)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //fmt.Errorf("exif data not found")
		}
		return nil, printError(err)
	}

	return &entity.ExifPhotoData{
		PhotoID: photoId,
		Data:    data,
	}, nil
}
