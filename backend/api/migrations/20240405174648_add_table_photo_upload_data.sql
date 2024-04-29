-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_upload_data (
    photo_id UUID PRIMARY KEY REFERENCES photos(id),
    -- Дата загрузки фотографии
    upload_at TIMESTAMPTZ NOT NULL,
    -- Оригинальные пути файлов фотографий
    paths text[] NOT NULL CHECK (CARDINALITY(paths) <= 2048),
    -- Информация клиента загрузившего фотографию
    client_info TEXT NOT NULL CHECK (LENGTH(client_info) <= 256),
    -- Информация о загрузившем фотографию
   person_id UUID NOT NULL REFERENCES people(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_upload_data;
-- +goose StatementEnd
