-- +goose Up
-- +goose StatementBegin
ALTER TYPE photo_processing_status ADD VALUE 'SIMILAR_COEFFICIENT';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
