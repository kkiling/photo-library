-- +goose Up
-- +goose StatementBegin
CREATE TYPE photo_status AS ENUM (
    'ACTIVE',
    'NOT_VALID'
);

CREATE TYPE photo_extension AS ENUM (
    'JPEG',
    'PNG'
);

CREATE TABLE photos (
    id UUID PRIMARY KEY,
    -- Ключ оригинального файла в хранилище
    file_key TEXT NOT NULL CHECK (LENGTH(file_key) <= 1024) UNIQUE,
    -- Хеш фотографии
    hash TEXT NOT NULL CHECK (LENGTH(hash) <= 512) UNIQUE,
    -- Дата последнего изменения оригинального фото
    photo_updated_at TIMESTAMPTZ NOT NULL,
    -- Расширение фотографии
    extension photo_extension NOT NULL,
    -- Статус текущей фотографии
    status photo_status NOT NULL DEFAULT 'ACTIVE',
    -- Текст ошибки (если фотография не валидна)
    error TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photos;
DROP TYPE photo_status;
DROP TYPE photo_extension;
-- +goose StatementEnd
