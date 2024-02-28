-- +goose Up
-- +goose StatementBegin
CREATE TABLE coeffs_similar_photos (
    photo_id1 UUID NOT NULL REFERENCES photos(id),
    photo_id2 UUID NOT NULL REFERENCES photos(id),
    coefficient FLOAT NOT NULL,
    PRIMARY KEY (photo_id1, photo_id2)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE coeffs_similar_photos;
-- +goose StatementEnd
