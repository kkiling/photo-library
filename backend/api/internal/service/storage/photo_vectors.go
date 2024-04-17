package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) GetPhotoVector(ctx context.Context, photoID uuid.UUID) (model.PhotoVector, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoVector(ctx, photoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PhotoVector{}, serviceerr.ErrNotFound
		}
		return model.PhotoVector{}, printError(err)
	}

	return model.PhotoVector{
		PhotoID: res.PhotoID,
		Vector:  res.Vector,
		Norm:    res.Norm,
	}, nil
}

func (r *Adapter) GetPhotoVectors(ctx context.Context, pagination model.Pagination) ([]model.PhotoVector, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoVectors(ctx, photo_library.GetPhotoVectorsParams{
		Offset: pagination.GetOffset(),
		Limit:  pagination.GetLimit(),
	})
	if err != nil {
		return nil, printError(err)
	}

	return lo.Map(res, func(item photo_library.PhotoVector, index int) model.PhotoVector {
		return model.PhotoVector{
			PhotoID: item.PhotoID,
			Vector:  item.Vector,
			Norm:    item.Norm,
		}
	}), nil
}

func (r *Adapter) SavePhotoVector(ctx context.Context, vector model.PhotoVector) error {
	queries := r.getQueries(ctx)

	err := queries.SavePhotoVector(ctx, photo_library.SavePhotoVectorParams{
		PhotoID: vector.PhotoID,
		Vector:  vector.Vector,
		Norm:    vector.Norm,
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) DeletePhotoVector(ctx context.Context, photoID uuid.UUID) error {
	queries := r.getQueries(ctx)

	_, err := queries.DeletePhotoVector(ctx, photoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return serviceerr.ErrNotFound
		}
		return printError(err)
	}

	return nil
}
