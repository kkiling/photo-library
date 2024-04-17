-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_vectors (
    photo_id UUID PRIMARY KEY REFERENCES photos(id),
    vector DOUBLE PRECISION[] NOT NULL,
    norm DOUBLE PRECISION NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_vectors;
-- +goose StatementEnd
