-- +goose Up
-- +goose StatementBegin
CREATE TYPE refresh_token_status AS ENUM (
    'ACTIVE',
    'REVOKED',
    'EXPIRED',
    'LOGOUT'
);

CREATE TABLE refresh_codes (
    id UUID PRIMARY KEY,
    person_id UUID NOT NULL REFERENCES people(id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    status refresh_token_status NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_codes;
DROP TYPE refresh_token_status;
-- +goose StatementEnd
