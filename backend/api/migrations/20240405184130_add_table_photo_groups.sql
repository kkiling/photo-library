-- +goose Up
-- +goose StatementBegin
CREATE TABLE photo_groups (
    id UUID PRIMARY KEY,
    main_photo_id UUID NOT NULL REFERENCES photos(id) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE photo_groups_photos (
    photo_id UUID  REFERENCES photos(id),
    group_id UUID REFERENCES photo_groups(id),
    PRIMARY KEY(photo_id, group_id)
);

CREATE INDEX idx_photos_groups_references_group_id ON photo_groups_photos(group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE photo_groups;
DROP TABLE photo_groups_photos;
-- +goose StatementEnd
