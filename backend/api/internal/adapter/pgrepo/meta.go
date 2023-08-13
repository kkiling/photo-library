package pgrepo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
)

func (r *PhotoRepository) SaveOrUpdateMetaData(ctx context.Context, data *entity.MetaData) error {
	conn := r.getConn(ctx)

	const query = `
			INSERT INTO meta_data (photo_id, model_info, size_bytes, width_pixel, height_pixel, date_time, geo_latitude, geo_longitude)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (photo_id) 
			DO UPDATE SET 
				model_info = EXCLUDED.model_info, 
				size_bytes = EXCLUDED.size_bytes, 
				width_pixel = EXCLUDED.width_pixel, 
				height_pixel = EXCLUDED.height_pixel, 
				date_time = EXCLUDED.date_time, 
				geo_latitude = EXCLUDED.geo_latitude, 
				geo_longitude = EXCLUDED.geo_longitude;
	`

	_, err := conn.Exec(ctx, query, data.PhotoID, data.ModelInfo, data.SizeBytes, data.WidthPixel,
		data.HeightPixel, data.DateTime, data.GeoLatitude, data.GeoLongitude)

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetMetaData(ctx context.Context, photoID uuid.UUID) (*entity.MetaData, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT photo_id, model_info, size_bytes, width_pixel, height_pixel, date_time, geo_latitude, geo_longitude
		FROM meta_data
		WHERE photo_id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, photoID)

	var meta entity.MetaData
	err := row.Scan(&meta.PhotoID, &meta.ModelInfo, &meta.SizeBytes, &meta.WidthPixel,
		&meta.HeightPixel, &meta.DateTime, &meta.GeoLatitude, &meta.GeoLongitude)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, printError(err)
	}

	return &meta, nil
}
