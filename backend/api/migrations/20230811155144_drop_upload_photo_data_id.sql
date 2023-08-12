-- +goose Up
-- +goose StatementBegin
ALTER TABLE upload_photo_data DROP COLUMN id;
ALTER TABLE upload_photo_data ADD PRIMARY KEY (photo_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE upload_photo_data DROP CONSTRAINT upload_photo_data_pkey;
ALTER TABLE upload_photo_data ADD COLUMN id UUID;
-- +goose StatementEnd
