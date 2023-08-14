-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_photos_hash ON photos (hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_photos_hash;
-- +goose StatementEnd
