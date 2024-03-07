-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_groups (
    id UUID PRIMARY KEY NOT NULL,
    main_photo_id UUID NOT NULL REFERENCES photos(id),
    update_at TIMESTAMP NOT NULL
);

CREATE TABLE photo_groups_photos (
    photo_id UUID NOT NULL REFERENCES photos(id),
    group_id UUID NOT NULL REFERENCES photo_groups(id),
    PRIMARY KEY (photo_id)
);

CREATE INDEX idx_photo_groups_photos_group_id ON photo_groups_photos(group_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_groups;
DROP TABLE photo_groups_photos;
-- +goose StatementEnd

