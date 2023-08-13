-- +goose Up
-- +goose StatementBegin
CREATE TABLE meta_data (
  photo_id UUID PRIMARY KEY,
  model_info TEXT,
  size_bytes INTEGER NOT NULL,
  width_pixel INTEGER NOT NULL,
  height_pixel INTEGER NOT NULL,
  date_time TIMESTAMP,
  update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  geo_latitude DOUBLE PRECISION,
  geo_longitude DOUBLE PRECISION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meta_data;
-- +goose StatementEnd
