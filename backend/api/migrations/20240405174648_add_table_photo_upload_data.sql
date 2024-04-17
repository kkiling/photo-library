-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_upload_data (
    photo_id UUID PRIMARY KEY REFERENCES photos(id),
    -- Дата загрузки фотографии
    upload_at TIMESTAMPTZ NOT NULL,
    -- Оригинальные пути файлов фотографий
    paths text[] NOT NULL CHECK (CARDINALITY(paths) <= 2048),
    -- Идентификатор клиента загрузившего фотографию
    client_id TEXT NOT NULL CHECK (LENGTH(client_id) <= 256)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_upload_data;
-- +goose StatementEnd
