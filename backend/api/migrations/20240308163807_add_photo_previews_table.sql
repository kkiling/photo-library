-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_previews (
    id UUID PRIMARY KEY NOT NULL,
    photo_id UUID NOT NULL REFERENCES photos(id),
    file_name VARCHAR(2048) NOT NULL,
    width_pixel INTEGER NOT NULL,
    height_pixel INTEGER NOT NULL,
    size_pixel INTEGER NOT NULL
);

CREATE INDEX idx_photo_previews_photo_id ON photo_previews(photo_id);

ALTER TYPE photo_processing_status ADD VALUE 'PHOTO_PREVIEW';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_previews;
-- +goose StatementEnd

