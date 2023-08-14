-- +goose Up
-- +goose StatementBegin
-- Для таблицы exif_data
ALTER TABLE exif_data
    DROP CONSTRAINT exif_data_photo_id_fkey;

ALTER TABLE exif_data
    ADD CONSTRAINT fk_exif_data_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id)
            ON DELETE CASCADE;

-- Для таблицы meta_data
ALTER TABLE meta_data
    DROP CONSTRAINT meta_data_photo_id_fkey;

ALTER TABLE meta_data
    ADD CONSTRAINT fk_meta_data_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id)
            ON DELETE CASCADE;

-- Для таблицы upload_photo_data
ALTER TABLE upload_photo_data
    DROP CONSTRAINT upload_photo_data_photo_id_fkey;

ALTER TABLE upload_photo_data
    ADD CONSTRAINT fk_upload_photo_data_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id)
            ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Для таблицы exif_data
ALTER TABLE exif_data
    DROP CONSTRAINT fk_exif_data_photo_id;

ALTER TABLE exif_data
    ADD CONSTRAINT exif_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);

-- Для таблицы meta_data
ALTER TABLE meta_data
    DROP CONSTRAINT fk_meta_data_photo_id;

ALTER TABLE meta_data
    ADD CONSTRAINT meta_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);

-- Для таблицы upload_photo_data
ALTER TABLE upload_photo_data
    DROP CONSTRAINT fk_upload_photo_data_photo_id;

ALTER TABLE upload_photo_data
    ADD CONSTRAINT upload_photo_data_photo_id_fkey
        FOREIGN KEY (photo_id) REFERENCES photos(id);
-- +goose StatementEnd
