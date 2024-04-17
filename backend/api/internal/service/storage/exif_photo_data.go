package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) GetExif(ctx context.Context, photoID uuid.UUID) (model.ExifPhotoData, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetExif(ctx, photoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ExifPhotoData{}, serviceerr.ErrNotFound
		}
		return model.ExifPhotoData{}, printError(err)
	}

	data := make(map[string]interface{})
	if err = json.Unmarshal(res.Data, &data); err != nil {
		return model.ExifPhotoData{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return model.ExifPhotoData{
		PhotoID: res.PhotoID,
		Data:    data,
	}, nil
}

func (r *Adapter) SaveExif(ctx context.Context, exif model.ExifPhotoData) error {
	queries := r.getQueries(ctx)

	data, err := json.Marshal(exif.Data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = queries.SaveExif(ctx, photo_library.SaveExifParams{
		PhotoID: exif.PhotoID,
		Data:    data,
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) DeleteExif(ctx context.Context, photoID uuid.UUID) error {
	queries := r.getQueries(ctx)

	_, err := queries.DeleteExif(ctx, photoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return serviceerr.ErrNotFound
		}
		return printError(err)
	}

	return nil
}
