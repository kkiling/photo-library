package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) SaveCoefficientSimilarPhotos(ctx context.Context, coefficient model.CoefficientSimilarPhoto) error {
	queries := r.getQueries(ctx)

	err := queries.SaveCoefficientSimilarPhoto(ctx, photo_library.SaveCoefficientSimilarPhotoParams{
		PhotoId1:    coefficient.PhotoID1,
		PhotoId2:    coefficient.PhotoID2,
		Coefficient: coefficient.Coefficient,
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) FindCoefficientSimilarPhoto(ctx context.Context, photoID uuid.UUID) ([]model.CoefficientSimilarPhoto, error) {
	queries := r.getQueries(ctx)

	res, err := queries.FindCoefficientSimilarPhoto(ctx, photoID)
	if err != nil {
		return nil, printError(err)
	}

	return lo.Map(res, func(item photo_library.CoefficientsSimilarPhoto, index int) model.CoefficientSimilarPhoto {
		return model.CoefficientSimilarPhoto{
			PhotoID1:    item.PhotoId1,
			PhotoID2:    item.PhotoId2,
			Coefficient: item.Coefficient,
		}
	}), nil
}

func (r *Adapter) DeleteCoefficientSimilarPhoto(ctx context.Context, photoID uuid.UUID) error {
	queries := r.getQueries(ctx)
	err := queries.DeleteCoefficientSimilarPhoto(ctx, photoID)
	if err != nil {
		return printError(err)
	}
	return nil
}
