-- +goose Up
-- +goose StatementBegin
CREATE TABLE rocket_locks(
    key text NOT NULL UNIQUE,
    locked_until timestamptz NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rocket_locks;
-- +goose StatementEnd
