-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_locations (
    photo_id UUID PRIMARY KEY REFERENCES photos(id),
    created_at          TIMESTAMPTZ NOT NULL,
    geo_latitude        DOUBLE PRECISION NOT NULL,
    geo_longitude       DOUBLE PRECISION NOT NULL,
    formatted_address   TEXT NOT NULL CHECK (LENGTH(formatted_address) <= 1024),
    street              TEXT NOT NULL CHECK (LENGTH(street) <= 1024),
    house_number        TEXT NOT NULL CHECK (LENGTH(house_number) <= 1024),
    suburb              TEXT NOT NULL CHECK (LENGTH(suburb) <= 1024),
    postcode            TEXT NOT NULL CHECK (LENGTH(postcode) <= 1024),
    state               TEXT NOT NULL CHECK (LENGTH(state) <= 1024),
    state_code          TEXT NOT NULL CHECK (LENGTH(state_code) <= 1024),
    state_district      TEXT NOT NULL CHECK (LENGTH(state_district) <= 1024),
    county              TEXT NOT NULL CHECK (LENGTH(county) <= 1024) ,
    country             TEXT NOT NULL CHECK (LENGTH(country) <= 1024),
    country_code        TEXT NOT NULL CHECK (LENGTH(country_code) <= 1024),
    city                TEXT NOT NULL CHECK (LENGTH(city) <= 1024)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_locations;
-- +goose StatementEnd
