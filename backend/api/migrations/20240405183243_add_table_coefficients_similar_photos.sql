-- +goose Up
-- +goose StatementBegin
CREATE TABLE coefficients_similar_photos (
    photo_id1 UUID REFERENCES photos(id),
    photo_id2 UUID REFERENCES photos(id),
    coefficient FLOAT NOT NULL,
    PRIMARY KEY (photo_id1, photo_id2)
);

CREATE INDEX idx_coefficients_similar_photos_photo_id2 ON coefficients_similar_photos USING btree (photo_id2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE coefficients_similar_photos;
-- +goose StatementEnd