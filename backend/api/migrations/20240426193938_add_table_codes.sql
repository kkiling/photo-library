-- +goose Up
-- +goose StatementBegin
CREATE TYPE code_type AS ENUM (
    'ACTIVATE_AUT'
);

CREATE TABLE codes (
    code TEXT PRIMARY KEY NOT NULL,
    person_id UUID REFERENCES people(id),
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
