-- +goose Up
-- +goose StatementBegin
CREATE TABLE file_hash (
   file_path TEXT PRIMARY KEY NOT NULL,
   update_at DATETIME NOT NULL,
   hash TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE file_hash;
-- +goose StatementEnd
