-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_tags_category_type ON tags_category(type);
CREATE UNIQUE INDEX idx_tags_photo_id_name ON tags(photo_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_tags_photo_id_name;
DROP INDEX idx_tags_category_type;
-- +goose StatementEnd
