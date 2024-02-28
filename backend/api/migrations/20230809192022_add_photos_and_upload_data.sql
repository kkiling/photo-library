-- +goose Up
-- +goose StatementBegin
CREATE TABLE photos (
    id UUID PRIMARY KEY NOT NULL,
    file_name VARCHAR(2048) NOT NULL,
    hash VARCHAR(512) NOT NULL,
    update_at TIMESTAMP NOT NULL,
    extension VARCHAR(8) NOT NULL
);

CREATE TABLE photo_upload_data (
   id UUID PRIMARY KEY NOT NULL,
   upload_at TIMESTAMP NOT NULL,
   photo_id UUID NOT NULL REFERENCES photos(id),
   paths TEXT[] NOT NULL CHECK (cardinality(paths) <= 2048),
   client_id VARCHAR(256) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_upload_data;
DROP TABLE photos;
-- +goose StatementEnd