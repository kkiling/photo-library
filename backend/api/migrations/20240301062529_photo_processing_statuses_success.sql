-- +goose Up
-- +goose StatementBegin
ALTER TABLE photo_processing_statuses ADD COLUMN success BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photo_processing_statuses DROP COLUMN success;
-- +goose StatementEnd
