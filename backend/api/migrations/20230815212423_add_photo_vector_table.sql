-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_vector (
   photo_id UUID PRIMARY KEY,
   vector DOUBLE PRECISION[] NOT NULL,
   norm DOUBLE PRECISION NOT NULL
);

ALTER TABLE photo_vector
    ADD CONSTRAINT fk_photo_vector_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_vector;
-- +goose StatementEnd
