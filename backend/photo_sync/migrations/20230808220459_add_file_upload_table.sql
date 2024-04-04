-- +goose Up
-- +goose StatementBegin
CREATE TABLE file_upload (
   hash TEXT PRIMARY KEY NOT NULL,
   upload_at DATETIME NOT NULL,
   success BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE file_upload;
-- +goose StatementEnd