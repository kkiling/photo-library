-- +goose Up
-- +goose StatementBegin
ALTER TABLE exif_photo_data
    ADD CONSTRAINT exif_photo_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);
ALTER TABLE photo_metadata
    ADD CONSTRAINT photo_metadata_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE exif_photo_data
    DROP CONSTRAINT exif_photo_data_photo_id_fkey;
ALTER TABLE photo_metadata
    DROP CONSTRAINT photo_metadata_photo_id_fkey;
-- +goose StatementEnd
