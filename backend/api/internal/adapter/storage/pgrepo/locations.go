package pgrepo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
)

func (r *PhotoRepository) SaveGeoAddress(ctx context.Context, location *entity.Location) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO locations (photo_id, created_at, geo_latitude, geo_longitude, 
		                       formatted_address, street, house_number, suburb, postcode, state, 
		                       state_code, state_district, county, country, country_code, city)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err := conn.Exec(ctx, query,
		location.PhotoID,
		location.CreatedAt,
		location.Latitude,
		location.Longitude,
		location.FormattedAddress,
		location.Street,
		location.HouseNumber,
		location.Suburb,
		location.Postcode,
		location.State,
		location.StateCode,
		location.StateDistrict,
		location.County,
		location.Country,
		location.CountryCode,
		location.City,
	)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetGeoAddress(ctx context.Context, photoID uuid.UUID) (*entity.Location, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT photo_id, created_at, geo_latitude, geo_longitude, 
			   formatted_address, street, house_number, suburb, postcode, state, 
			   state_code, state_district, county, country, country_code, city
		FROM locations
		WHERE photo_id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, photoID)

	var geo entity.Location
	err := row.Scan(&geo.PhotoID, &geo.CreatedAt, &geo.Latitude, &geo.Longitude,
		&geo.FormattedAddress, &geo.Street, &geo.HouseNumber, &geo.Suburb, &geo.Postcode, &geo.State,
		&geo.StateCode, &geo.StateDistrict, &geo.County, &geo.Country, &geo.CountryCode, &geo.County)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, printError(err)
	}

	return &geo, nil
}
