-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_coeffs_similar_photos_photo_id2 ON coeffs_similar_photos(photo_id2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_coeffs_similar_photos_photo_id2;
-- +goose StatementEnd

