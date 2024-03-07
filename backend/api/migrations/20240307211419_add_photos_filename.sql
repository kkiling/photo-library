-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_photos_file_name ON photos (file_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_photos_file_name;
-- +goose StatementEnd
