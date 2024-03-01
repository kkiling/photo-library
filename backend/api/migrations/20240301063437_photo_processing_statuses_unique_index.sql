-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_photo_processing_statuses_photo_id_status ON photo_processing_statuses(photo_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_photo_processing_statuses_photo_id_status;
-- +goose StatementEnd



