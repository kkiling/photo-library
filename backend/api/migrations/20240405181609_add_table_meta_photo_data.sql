-- +goose Up
-- +goose StatementBegin
CREATE TABLE meta_photo_data (
    photo_id UUID PRIMARY KEY REFERENCES photos(id),
    model_info TEXT CHECK (LENGTH(model_info) <= 512),
    size_bytes INTEGER NOT NULL,
    width_pixel INTEGER NOT NULL,
    height_pixel INTEGER NOT NULL,
    -- Дата время сьемки
    date_time TIMESTAMPTZ,
    -- Дата последнего обновления файла
    updated_at  TIMESTAMPTZ NOT NULL,
    geo_latitude DOUBLE PRECISION,
    geo_longitude DOUBLE PRECISION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meta_photo_data;
-- +goose StatementEnd
