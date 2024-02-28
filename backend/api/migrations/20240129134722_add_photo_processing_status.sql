-- +goose Up
-- +goose StatementBegin
CREATE TYPE photo_processing_status AS ENUM (
    'NEW_PHOTO',
    'EXIF_DATA_SAVED',
    'META_DATA_SAVED',
    'SYSTEM_TAGS_SAVED',
    'PHOTO_VECTOR_SAVED'
);

ALTER TABLE photos ADD COLUMN processing_status photo_processing_status DEFAULT 'NEW_PHOTO';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photos DROP COLUMN processing_status;
DROP TYPE photo_processing_status;
-- +goose StatementEnd
