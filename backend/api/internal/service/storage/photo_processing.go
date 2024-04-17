package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) AddPhotoProcessing(ctx context.Context, processing model.PhotoProcessing) error {
	queries := r.getQueries(ctx)

	err := queries.AddPhotoProcessing(ctx, photo_library.AddPhotoProcessingParams{
		PhotoID:     processing.PhotoID,
		ProcessedAt: processing.ProcessedAt,
		Type:        photo_library.ProcessingType(processing.ProcessingType),
		Success:     processing.Success,
	})
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetPhotoProcessingTypes(ctx context.Context, photoID uuid.UUID) ([]model.PhotoProcessing, error) {
	queries := r.getQueries(ctx)

	processing, err := queries.GetPhotoProcessing(ctx, photoID)

	if err != nil {
		return nil, printError(err)
	}

	result := lo.Map(processing, func(p photo_library.PhotoProcessing, _ int) model.PhotoProcessing {
		return model.PhotoProcessing{
			PhotoID:        p.PhotoID,
			ProcessedAt:    p.ProcessedAt,
			ProcessingType: model.ProcessingType(p.Type),
			Success:        p.Success,
		}
	})

	return result, nil
}

func (r *Adapter) GetUnprocessedPhotoIDs(ctx context.Context, processingTypes []model.ProcessingType, limit int) ([]uuid.UUID, error) {
	queries := r.getQueries(ctx)
	var err error
	var photoIDs []uuid.UUID
	for _, processingType := range processingTypes {
		photoIDs, err = queries.GetUnprocessedPhotos(ctx, photo_library.GetUnprocessedPhotosParams{
			Type:  photo_library.ProcessingType(processingType),
			Limit: int32(limit),
		})

		if err != nil {
			return nil, printError(err)
		}

		if len(photoIDs) > 0 {
			return photoIDs, nil
		}
	}

	return nil, nil
}
