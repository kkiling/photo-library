-- +goose Up
-- +goose StatementBegin
ALTER TABLE photo_upload_data DROP COLUMN id;
ALTER TABLE photo_upload_data ADD PRIMARY KEY (photo_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE photo_upload_data DROP CONSTRAINT photo_upload_data_pkey;
ALTER TABLE photo_upload_data ADD COLUMN id UUID;
-- +goose StatementEnd
