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

func (r *Adapter) GetPhotoById(ctx context.Context, id uuid.UUID) (model.Photo, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Photo{}, serviceerr.ErrNotFound
		}
		return model.Photo{}, printError(err)
	}

	return model.Photo{
		ID:        res.ID,
		FileKey:   res.FileKey,
		Hash:      res.Hash,
		UpdateAt:  res.UpdatedAt,
		Extension: model.PhotoExtension(res.Extension),
	}, nil
}

func (r *Adapter) GetPhotoByFilename(ctx context.Context, fileName string) (model.Photo, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoByFileKey(ctx, fileName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Photo{}, serviceerr.ErrNotFound
		}
		return model.Photo{}, printError(err)
	}

	return model.Photo{
		ID:        res.ID,
		FileKey:   res.FileKey,
		Hash:      res.Hash,
		UpdateAt:  res.UpdatedAt,
		Extension: model.PhotoExtension(res.Extension),
	}, nil
}

func (r *Adapter) GetPhotoByHash(ctx context.Context, hash string) (model.Photo, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Photo{}, serviceerr.ErrNotFound
		}
		return model.Photo{}, printError(err)
	}

	return model.Photo{
		ID:        res.ID,
		FileKey:   res.FileKey,
		Hash:      res.Hash,
		UpdateAt:  res.UpdatedAt,
		Extension: model.PhotoExtension(res.Extension),
	}, nil
}

func (r *Adapter) SavePhoto(ctx context.Context, photo model.Photo) error {
	queries := r.getQueries(ctx)

	err := queries.SavePhoto(ctx, photo_library.SavePhotoParams{
		ID:        photo.ID,
		FileKey:   photo.FileKey,
		Hash:      photo.Hash,
		UpdatedAt: photo.UpdateAt,
		Extension: photo_library.PhotoExtension(photo.Extension),
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) MakeNotValidPhoto(ctx context.Context, photoID uuid.UUID, error string) error {
	queries := r.getQueries(ctx)

	err := queries.MakeNotValidPhoto(ctx, photo_library.MakeNotValidPhotoParams{
		Error: &error,
		ID:    photoID,
	})

	if err != nil {
		return printError(err)
	}

	return nil
}
