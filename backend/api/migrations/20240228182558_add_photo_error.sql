-- +goose Up
-- +goose StatementBegin
ALTER TABLE photos ADD COLUMN error TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photos DROP COLUMN error;
-- +goose StatementEnd
