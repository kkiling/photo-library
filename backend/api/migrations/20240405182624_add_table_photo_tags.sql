-- +goose Up
-- +goose StatementBegin
CREATE TABLE tag_categories (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL UNIQUE CHECK (LENGTH(type) <= 64) ,
    color TEXT NOT NULL CHECK (LENGTH(color) <= 7)
);

CREATE TABLE photo_tags (
    id UUID PRIMARY KEY,
    category_id UUID NOT NULL REFERENCES tag_categories(id),
    photo_id UUID NOT NULL REFERENCES photos(id),
    name TEXT  NOT NULL CHECK (LENGTH(name) <= 128)
);

CREATE UNIQUE INDEX idx_tags_photo_id_name ON photo_tags(photo_id, name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tag_categories;
DROP TABLE photo_tags;
-- +goose StatementEnd
