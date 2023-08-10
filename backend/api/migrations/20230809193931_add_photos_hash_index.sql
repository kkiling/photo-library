-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_photos_hash ON photos (hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_photos_hash;
-- +goose StatementEnd
