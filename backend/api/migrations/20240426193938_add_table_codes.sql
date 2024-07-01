-- +goose Up
-- +goose StatementBegin
CREATE TYPE code_type AS ENUM (
    'ACTIVATE_INVITE',
    'ACTIVATE_REGISTRATION'
);

CREATE TABLE codes (
    code TEXT PRIMARY KEY,
    person_id UUID NOT NULL REFERENCES people(id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    active BOOLEAN NOT NULL,
    type code_type NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE codes;
DROP TYPE code_type;
-- +goose StatementEnd
