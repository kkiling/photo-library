-- +goose Up
-- +goose StatementBegin
CREATE TYPE photo_processing_status AS ENUM (
    'NEW_PHOTO',
    'SAVE_EXIF_DATA',
    'SAVE_META_DATA',
    'CREATE_TAGS_BY_META',
    'SAVE_PHOTO_VECTOR'
);

ALTER TABLE photos ADD COLUMN processing_status photo_processing_status DEFAULT 'NEW_PHOTO';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photos DROP COLUMN processing_status;
DROP TYPE photo_processing_status;
-- +goose StatementEnd
