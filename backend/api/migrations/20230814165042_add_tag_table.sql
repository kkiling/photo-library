-- +goose Up
-- +goose StatementBegin
CREATE TABLE tag_category (
  id UUID PRIMARY KEY,
  type VARCHAR(128) NOT NULL,
  color VARCHAR(7) NOT NULL
);

CREATE TABLE tag (
 id UUID PRIMARY KEY,
 category_id UUID NOT NULL,
 photo_id UUID NOT NULL,
 name VARCHAR(128) NOT NULL
);

ALTER TABLE tag
    ADD CONSTRAINT fk_tag_category_id
        FOREIGN KEY (category_id) REFERENCES tag_category(id) ON DELETE CASCADE;

ALTER TABLE tag
    ADD CONSTRAINT fk_tag_photo_id
        FOREIGN KEY (photo_id) REFERENCES photos(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tag;
DROP TABLE tag_category;
-- +goose StatementEnd
