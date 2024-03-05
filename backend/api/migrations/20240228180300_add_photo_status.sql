-- +goose Up
-- +goose StatementBegin
CREATE TYPE photo_status AS ENUM (
    'NEW_PHOTO',
    'NOT_VALID'
);

ALTER TABLE photos ADD COLUMN status photo_status DEFAULT 'NEW_PHOTO';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photos DROP COLUMN status;
DROP TYPE photo_status;
-- +goose StatementEnd
