-- +goose Up
-- +goose StatementBegin
CREATE TABLE photos (
    id UUID PRIMARY KEY NOT NULL,
    file_path VARCHAR(2048) NOT NULL,
    hash VARCHAR(512) NOT NULL,
    update_at TIMESTAMP NOT NULL,
    upload_at TIMESTAMP NOT NULL,
    extension VARCHAR(8) NOT NULL
);

CREATE TABLE upload_photo_data (
   id UUID PRIMARY KEY NOT NULL,
   photo_id UUID NOT NULL REFERENCES photos(id),
   paths TEXT[] NOT NULL CHECK (cardinality(paths) <= 2048),
   upload_at TIMESTAMP NOT NULL,
   client_id VARCHAR(256) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE upload_photo_data;
DROP TABLE photos;
-- +goose StatementEnd