-- +goose Up
-- +goose StatementBegin
CREATE TABLE rocket_locks(
    key TEXT NOT NULL UNIQUE CHECK (LENGTH(key) <= 128),
    locked_until TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rocket_locks;
-- +goose StatementEnd
