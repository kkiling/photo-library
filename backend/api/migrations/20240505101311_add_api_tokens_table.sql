-- +goose Up
-- +goose StatementBegin
CREATE TYPE api_token_type AS ENUM (
    'SYNC_PHOTO'
);

CREATE TABLE api_tokens (
    id UUID PRIMARY KEY,
    person_id UUID NOT NULL REFERENCES people(id),
    caption TEXT NOT NULL CHECK (LENGTH(caption) <= 128),
    token TEXT NOT NULL UNIQUE CHECK (LENGTH(caption) <= 32),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    expired_at TIMESTAMPTZ,
    type api_token_type NOT NULL
);

CREATE INDEX idx_api_tokens_person_id ON api_tokens(person_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE api_tokens;
DROP TYPE api_token_type;
-- +goose StatementEnd
