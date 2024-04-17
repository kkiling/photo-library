-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_previews (
    id UUID PRIMARY KEY,
    photo_id UUID NOT NULL REFERENCES photos(id),
    file_key TEXT NOT NULL CHECK (LENGTH(file_key) <= 1024),
    size_pixel INTEGER NOT NULL,
    width_pixel INTEGER NOT NULL,
    height_pixel INTEGER NOT NULL,
    original BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_photo_previews_photo_id ON photo_previews(photo_id);
CREATE UNIQUE INDEX idx_photo_previews_file_key ON photo_previews(file_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_previews;
-- +goose StatementEnd
