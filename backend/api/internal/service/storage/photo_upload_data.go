package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) GetPhotoUploadData(ctx context.Context, photoID uuid.UUID) (model.PhotoUploadData, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoUploadData(ctx, photoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PhotoUploadData{}, serviceerr.ErrNotFound
		}
		return model.PhotoUploadData{}, printError(err)
	}

	return model.PhotoUploadData{
		PhotoID:    res.PhotoID,
		UploadAt:   res.UploadAt,
		Paths:      res.Paths,
		ClientInfo: res.ClientInfo,
		PersonID:   res.PersonID,
	}, nil
}

func (r *Adapter) SavePhotoUploadData(ctx context.Context, uploadData model.PhotoUploadData) error {
	queries := r.getQueries(ctx)

	err := queries.SavePhotoUploadData(ctx, photo_library.SavePhotoUploadDataParams{
		PhotoID:    uploadData.PhotoID,
		Paths:      uploadData.Paths,
		UploadAt:   uploadData.UploadAt,
		ClientInfo: uploadData.ClientInfo,
		PersonID:   uploadData.PersonID,
	})

	if err != nil {
		return printError(err)
	}

	return nil
}
