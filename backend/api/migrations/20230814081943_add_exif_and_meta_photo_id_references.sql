-- +goose Up
-- +goose StatementBegin
ALTER TABLE exif_data
    ADD CONSTRAINT exif_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);
ALTER TABLE meta_data
    ADD CONSTRAINT meta_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE exif_data
    DROP CONSTRAINT exif_data_photo_id_fkey;
ALTER TABLE meta_data
    DROP CONSTRAINT meta_data_photo_id_fkey;
-- +goose StatementEnd
