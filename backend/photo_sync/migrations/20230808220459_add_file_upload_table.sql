-- +goose Up
CREATE TABLE file_upload (
   hash TEXT PRIMARY KEY NOT NULL,
   upload_at DATETIME NOT NULL,
   success BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE file_upload;
