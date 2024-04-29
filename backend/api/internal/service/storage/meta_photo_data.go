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

func (r *Adapter) GetMetadata(ctx context.Context, photoID uuid.UUID) (model.PhotoMetadata, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetMetadata(ctx, photoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PhotoMetadata{}, serviceerr.ErrNotFound
		}
		return model.PhotoMetadata{}, printError(err)
	}

	var geo *model.Geo
	if res.GeoLongitude != nil && res.GeoLatitude != nil {
		geo = &model.Geo{
			Latitude:  *res.GeoLatitude,
			Longitude: *res.GeoLongitude,
		}
	}

	return model.PhotoMetadata{
		PhotoID:        res.PhotoID,
		ModelInfo:      res.ModelInfo,
		SizeBytes:      res.SizeBytes,
		WidthPixel:     res.WidthPixel,
		HeightPixel:    res.HeightPixel,
		DateTime:       dataTimePtr(res.DateTime),
		PhotoUpdatedAt: res.UpdatedAt,
		Geo:            geo,
	}, nil
}

func (r *Adapter) SaveMetadata(ctx context.Context, metadata model.PhotoMetadata) error {
	queries := r.getQueries(ctx)

	var latitude *float64
	var longitude *float64
	if metadata.Geo != nil {
		latitude = &metadata.Geo.Latitude
		longitude = &metadata.Geo.Longitude
	}

	params := photo_library.SaveMetadataParams{
		PhotoID:      metadata.PhotoID,
		ModelInfo:    metadata.ModelInfo,
		SizeBytes:    metadata.SizeBytes,
		WidthPixel:   metadata.WidthPixel,
		HeightPixel:  metadata.HeightPixel,
		DateTime:     pgTypeTimestamptz(metadata.DateTime),
		UpdatedAt:    metadata.PhotoUpdatedAt,
		GeoLatitude:  latitude,
		GeoLongitude: longitude,
	}

	err := queries.SaveMetadata(ctx, params)

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) DeleteMetadata(ctx context.Context, photoID uuid.UUID) error {
	queries := r.getQueries(ctx)

	_, err := queries.DeleteMetadata(ctx, photoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return serviceerr.ErrNotFound
		}
		return printError(err)
	}

	return nil
}
