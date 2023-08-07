-- +goose Up
CREATE TABLE file_hash (
   file_path TEXT PRIMARY KEY NOT NULL,
   update_at DATETIME NOT NULL,
   hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE file_hash;
