-- +goose Up
-- +goose StatementBegin
CREATE TYPE processing_type AS ENUM (
    'EXIF_DATA',
    'META_DATA',
    'CATALOG_TAGS',
    'META_TAGS',
    'PHOTO_VECTOR',
    'SIMILAR_COEFFICIENT',
    'PHOTO_GROUP',
    'PHOTO_PREVIEW'
);

CREATE TABLE photo_processing (
    photo_id UUID REFERENCES photos(id),
    processed_at TIMESTAMPTZ NOT NULL,
    type processing_type NOT NULL,
    success BOOL NOT NULL,
    PRIMARY KEY (photo_id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_processing;
DROP TYPE processing_type;
-- +goose StatementEnd
