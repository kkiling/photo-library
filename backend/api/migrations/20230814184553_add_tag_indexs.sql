-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_tag_category_type ON tag_category(type);
CREATE UNIQUE INDEX idx_tag_photo_id_name ON tag(photo_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_tag_photo_id_name;
DROP INDEX idx_tag_category_type;
-- +goose StatementEnd
