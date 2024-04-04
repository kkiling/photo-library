-- +goose Up
-- +goose StatementBegin
CREATE TABLE locations (
    photo_id UUID       PRIMARY KEY,
    created_at          TIMESTAMP NOT NULL,
    geo_latitude        DOUBLE PRECISION,
    geo_longitude       DOUBLE PRECISION,
    formatted_address   TEXT,
    street              TEXT,
    house_number        TEXT,
    suburb              TEXT,
    postcode            TEXT,
    state               TEXT,
    state_code          TEXT,
    state_district      TEXT,
    county              TEXT,
    country             TEXT,
    country_code        TEXT,
    city                TEXT
);

ALTER TABLE locations
    ADD CONSTRAINT fk_locations_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id) ON DELETE CASCADE;

ALTER TABLE locations
    ADD CONSTRAINT locations_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE locations;
-- +goose StatementEnd
