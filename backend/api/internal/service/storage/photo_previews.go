package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) SavePhotoPreview(ctx context.Context, preview model.PhotoPreview) error {
	queries := r.getQueries(ctx)

	err := queries.SavePhotoPreview(ctx, photo_library.SavePhotoPreviewParams{
		ID:          preview.ID,
		PhotoID:     preview.PhotoID,
		FileKey:     preview.FileKey,
		SizePixel:   preview.SizePixel,
		WidthPixel:  preview.WidthPixel,
		HeightPixel: preview.HeightPixel,
		Original:    preview.Original,
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetPhotoPreviews(ctx context.Context, photoID uuid.UUID) ([]model.PhotoPreview, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoPreviews(ctx, photoID)

	if err != nil {
		return nil, printError(err)
	}

	return lo.Map(res, func(row photo_library.PhotoPreview, index int) model.PhotoPreview {
		return model.PhotoPreview{
			ID:          row.ID,
			PhotoID:     photoID,
			FileKey:     row.FileKey,
			WidthPixel:  row.WidthPixel,
			HeightPixel: row.HeightPixel,
			SizePixel:   row.SizePixel,
			Original:    row.Original,
		}
	}), nil
}

func (r *Adapter) DeletePhotoPreviews(ctx context.Context, photoID uuid.UUID) error {
	queries := r.getQueries(ctx)

	err := queries.DeletePhotoPreviews(ctx, photoID)
	if err != nil {
		return printError(err)
	}

	return nil
}
