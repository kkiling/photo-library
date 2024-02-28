-- +goose Up
-- +goose StatementBegin
-- Для таблицы exif_photo_data
ALTER TABLE exif_photo_data
    DROP CONSTRAINT exif_photo_data_photo_id_fkey;

ALTER TABLE exif_photo_data
    ADD CONSTRAINT fk_exif_photo_data_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id)
            ON DELETE CASCADE;

-- Для таблицы photo_metadata
ALTER TABLE photo_metadata
    DROP CONSTRAINT photo_metadata_photo_id_fkey;

ALTER TABLE photo_metadata
    ADD CONSTRAINT fk_photo_metadata_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id)
            ON DELETE CASCADE;

-- Для таблицы photo_upload_data
ALTER TABLE photo_upload_data
    DROP CONSTRAINT photo_upload_data_photo_id_fkey;

ALTER TABLE photo_upload_data
    ADD CONSTRAINT fk_photo_upload_data_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id)
            ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Для таблицы exif_photo_data
ALTER TABLE exif_photo_data
    DROP CONSTRAINT fk_exif_photo_data_photo_id;

ALTER TABLE exif_photo_data
    ADD CONSTRAINT exif_photo_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);

-- Для таблицы photo_metadata
ALTER TABLE photo_metadata
    DROP CONSTRAINT fk_photo_metadata_photo_id;

ALTER TABLE photo_metadata
    ADD CONSTRAINT photo_metadata_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);

-- Для таблицы photo_upload_data
ALTER TABLE photo_upload_data
    DROP CONSTRAINT fk_photo_upload_data_photo_id;

ALTER TABLE photo_upload_data
    ADD CONSTRAINT photo_upload_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);
-- +goose StatementEnd
