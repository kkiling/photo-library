-- +goose Up
-- +goose StatementBegin
CREATE TYPE photo_processing_status AS ENUM (
    'EXIF_DATA',
    'META_DATA',
    'CATALOG_TAGS',
    'META_TAGS',
    'PHOTO_VECTOR'
);

CREATE TABLE photo_processing_statuses (
    photo_id UUID NOT NULL,
    processed_at TIMESTAMP NOT NULL,
    status photo_processing_status NOT NULL
);

ALTER TABLE photo_processing_statuses
    ADD CONSTRAINT fk_photo_processing_statuses_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_processing_statuses;
DROP TYPE photo_processing_status;
-- +goose StatementEnd
