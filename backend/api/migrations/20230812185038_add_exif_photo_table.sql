-- +goose Up
-- +goose StatementBegin
CREATE TABLE exif_photo_data (
   photo_id UUID PRIMARY KEY,
   data JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exif_photo_data;
-- +goose StatementEnd
