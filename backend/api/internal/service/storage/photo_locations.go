package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/geo"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) SaveGeoAddress(ctx context.Context, location model.Location) error {
	queries := r.getQueries(ctx)

	params := photo_library.SavePhotoLocationParams{
		PhotoID:          location.PhotoID,
		CreatedAt:        location.CreatedAt,
		GeoLatitude:      location.Latitude,
		GeoLongitude:     location.Longitude,
		FormattedAddress: location.Address.FormattedAddress,
		Street:           location.Address.Street,
		HouseNumber:      location.Address.HouseNumber,
		Suburb:           location.Address.Suburb,
		Postcode:         location.Address.Postcode,
		State:            location.Address.State,
		StateCode:        location.Address.StateCode,
		StateDistrict:    location.Address.StateDistrict,
		County:           location.Address.County,
		Country:          location.Address.Country,
		CountryCode:      location.Address.CountryCode,
		City:             location.Address.City,
	}

	err := queries.SavePhotoLocation(ctx, params)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetGeoAddress(ctx context.Context, photoID uuid.UUID) (model.Location, error) {
	queries := r.getQueries(ctx)
	res, err := queries.GetGeoAddress(ctx, photoID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Location{}, serviceerr.ErrNotFound
		}
		return model.Location{}, printError(err)
	}

	return model.Location{
		PhotoID:   res.PhotoID,
		CreatedAt: res.CreatedAt,
		Latitude:  res.GeoLatitude,
		Longitude: res.GeoLongitude,
		Address: geo.Address{
			FormattedAddress: res.FormattedAddress,
			Street:           res.Street,
			HouseNumber:      res.HouseNumber,
			Suburb:           res.Suburb,
			Postcode:         res.Postcode,
			State:            res.State,
			StateCode:        res.StateCode,
			StateDistrict:    res.StateDistrict,
			County:           res.County,
			Country:          res.Country,
			CountryCode:      res.CountryCode,
			City:             res.City,
		},
	}, nil
}
