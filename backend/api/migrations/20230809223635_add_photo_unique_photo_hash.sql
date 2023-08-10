-- +goose Up
-- +goose StatementBegin
ALTER TABLE photos ADD CONSTRAINT unique_photo_hash UNIQUE (hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photos DROP CONSTRAINT unique_photo_hash;
-- +goose StatementEnd
