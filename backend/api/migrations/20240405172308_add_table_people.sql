-- +goose Up
-- +goose StatementBegin
CREATE TABLE people (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    firstname TEXT NOT NULL CHECK (LENGTH(firstname) <= 1024),
    surname TEXT NOT NULL CHECK (LENGTH(surname) <= 1024),
    patronymic TEXT CHECK (LENGTH(patronymic) <= 1024)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE people;
-- +goose StatementEnd
