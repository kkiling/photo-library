-- name: GetPhotoById :one
SELECT id, file_key, hash, photo_updated_at, extension
FROM photos
WHERE id = $1 AND status = 'ACTIVE'
LIMIT 1;

-- name: GetPhotoByFileKey :one
SELECT id, file_key, hash, photo_updated_at, extension
FROM photos
WHERE file_key = $1 AND status = 'ACTIVE'
LIMIT 1;

-- name: GetPhotoByHash :one
SELECT id, file_key, hash, photo_updated_at, extension
FROM photos
WHERE hash = $1 AND status = 'ACTIVE'
LIMIT 1;

-- name: SavePhoto :exec
INSERT INTO photos (id, file_key, hash, photo_updated_at, extension, status)
VALUES ($1, $2, $3, $4, $5, 'ACTIVE');

-- name: MakeNotValidPhoto :exec
UPDATE photos SET status='NOT_VALID', error=$1 WHERE id=$2;

------------------------------------------------------------------------------------------------------------------------

-- name: SavePhotoUploadData :exec
INSERT INTO photo_upload_data (photo_id, paths, upload_at, client_info, person_id)
VALUES ($1, $2, $3, $4, $5);

-- name: GetPhotoUploadData :one
SELECT photo_id, paths, upload_at, client_info, person_id
FROM photo_upload_data
WHERE photo_id=$1
LIMIT 1;

------------------------------------------------------------------------------------------------------------------------

-- name: AddPhotoProcessing :exec
INSERT INTO photo_processing (photo_id, processed_at, type, success)
VALUES ($1, $2, $3, $4);

-- name: GetPhotoProcessing :many
SELECT photo_id, processed_at, type, success FROM photo_processing
WHERE photo_id=$1
ORDER BY processed_at;

-- name: GetUnprocessedPhotos :many
SELECT p.id FROM photos p
LEFT JOIN photo_processing ps ON p.id = ps.photo_id AND ps.type = $1
WHERE ps.photo_id is NULL and p.status = 'ACTIVE'
ORDER BY p.photo_updated_at -- TODO: Индекс на updated_at
LIMIT $2;
------------------------------------------------------------------------------------------------------------------------

-- name: SaveExif :exec
INSERT INTO exif_photo_data (photo_id, data)
VALUES ($1, $2);

-- name: GetExif :one
SELECT photo_id, data FROM exif_photo_data
WHERE photo_id=$1
LIMIT 1;

-- name: DeleteExif :one
DELETE FROM exif_photo_data
WHERE photo_id=$1
RETURNING photo_id;

------------------------------------------------------------------------------------------------------------------------

-- name: SaveMetadata :exec
INSERT INTO meta_photo_data (photo_id, model_info, size_bytes, width_pixel, height_pixel,
                             date_time, updated_at, geo_latitude, geo_longitude)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetMetadata :one
SELECT photo_id, model_info, size_bytes, width_pixel, height_pixel,
       date_time, updated_at, geo_latitude, geo_longitude FROM meta_photo_data
WHERE photo_id=$1
LIMIT 1;

-- name: DeleteMetadata :one
DELETE FROM meta_photo_data
WHERE photo_id=$1
RETURNING photo_id;

------------------------------------------------------------------------------------------------------------------------

-- name: GetTagCategory :one
SELECT id, type, color
FROM tag_categories
WHERE id = $1
LIMIT 1;

-- name: GetTagCategoryByType :one
SELECT id, type, color
FROM tag_categories
WHERE type = $1
LIMIT 1;

-- name: SaveTagCategory :exec
INSERT INTO tag_categories (id, type, color)
VALUES ($1, $2, $3);

-- name: SaveTag :exec
INSERT INTO photo_tags (id, category_id, photo_id, name)
VALUES ($1, $2, $3, $4);

-- name: GetTags :many
SELECT id, category_id, name FROM photo_tags
WHERE photo_id=$1;

-- name: DeleteTag :one
DELETE FROM photo_tags
WHERE id=$1
RETURNING id;

-- name: DeletePhotoTagsByCategories :one
DELETE FROM photo_tags
WHERE photo_id=$1 and category_id = ANY (sqlc.arg(category_ids)::uuid[]) -- IN (sqlc.slice(category_ids)::uuid[])
RETURNING photo_id;

------------------------------------------------------------------------------------------------------------------------

-- name: SavePhotoVector :exec
INSERT INTO photo_vectors (photo_id, vector, norm)
VALUES ($1, $2, $3);

-- name: GetPhotoVector :one
SELECT photo_id, vector, norm FROM photo_vectors
WHERE photo_id=$1
LIMIT 1;

-- name: DeletePhotoVector :one
DELETE FROM photo_vectors
WHERE photo_id=$1
RETURNING photo_id;

-- name: GetPhotoVectors :many
SELECT photo_id, vector, norm
FROM photo_vectors
OFFSET $1
LIMIT $2;

------------------------------------------------------------------------------------------------------------------------

-- name: SaveCoefficientSimilarPhoto :exec
INSERT INTO coefficients_similar_photos (photo_id1, photo_id2, coefficient)
VALUES ($1, $2, $3);

-- name: FindCoefficientSimilarPhoto :many
SELECT photo_id1, photo_id2, coefficient
FROM coefficients_similar_photos
WHERE photo_id1 = $1 OR photo_id2 = $1;

-- name: DeleteCoefficientSimilarPhoto :exec
DELETE FROM coefficients_similar_photos
WHERE photo_id1 = $1 OR photo_id2 = $1;

------------------------------------------------------------------------------------------------------------------------

-- name: FindGroupIDByPhotoID :one
SELECT group_id
FROM photo_groups_photos
WHERE photo_id = $1
LIMIT 1;

-- name: GetGroupByID :many
SELECT id, main_photo_id, updated_at, created_at, p.photo_id as photo_id
FROM photo_groups
LEFT JOIN photo_groups_photos p ON photo_groups.id = p.group_id
WHERE id = $1;

-- name: SaveGroup :exec
INSERT INTO photo_groups (id, main_photo_id, created_at, updated_at)
VALUES ($1, $2, $3, $4);

-- name: AddPhotoIDToGroup :exec
INSERT INTO photo_groups_photos (photo_id, group_id)
VALUES ($1, $2);

-- name: GetPhotoGroupsCount :one
SELECT count(1) as count FROM photo_groups;

-- name: GetPaginatedPhotoGroups :many
SELECT id, main_photo_id, updated_at, created_at, p.photo_id as photo_id FROM photo_groups
LEFT JOIN photo_groups_photos p ON photo_groups.id = p.group_id
OFFSET $1 LIMIT $2;

-- name: GetGroupPhotoIDs :many
SELECT photo_id
FROM photo_groups_photos
WHERE group_id = $1;

-- name: SetPhotoGroupMainPhoto :exec
UPDATE photo_groups
SET main_photo_id = $1, updated_at = $2
WHERE id = $3;

-- name: DeletePhotoGroup :one
DELETE FROM photo_groups where id = $1  RETURNING id;

-- name: DeletePhotoGroupPhotos :exec
DELETE FROM photo_groups_photos where group_id = $1;

-- name: DeletePhotoGroupByMainPhoto :exec
DELETE FROM photo_groups where main_photo_id = $1;

-- name: DeletePhotoGroupPhotosByPhoto :exec
DELETE FROM photo_groups_photos where photo_id = $1;

------------------------------------------------------------------------------------------------------------------------

-- name: SavePhotoPreview :exec
INSERT INTO photo_previews (id, photo_id, file_key, size_pixel, width_pixel, height_pixel, original)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetPhotoPreviews :many
SELECT id, photo_id, file_key, size_pixel, width_pixel, height_pixel, original FROM photo_previews
WHERE photo_id = $1
ORDER BY size_pixel;

-- name: DeletePhotoPreviews :exec
DELETE FROM photo_previews
WHERE photo_id=$1;

------------------------------------------------------------------------------------------------------------------------

-- name: SavePhotoLocation :exec
INSERT INTO photo_locations (photo_id, created_at, geo_latitude, geo_longitude,
                             formatted_address, street, house_number, suburb,
                             postcode, state, state_code, state_district, county,
                             country, country_code, city)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

-- name: GetGeoAddress :one
SELECT photo_id, created_at, geo_latitude, geo_longitude,
       formatted_address, street, house_number, suburb,
       postcode, state, state_code, state_district, county,
       country, country_code, city
FROM photo_locations WHERE photo_id = $1;

------------------------------------------------------------------------------------------------------------------------

-- name: RocketLock :one
INSERT INTO rocket_locks (key, locked_until) VALUES ($1, now() + sqlc.narg('interval')::interval)
ON CONFLICT (key) DO UPDATE SET locked_until = (now() + sqlc.narg('interval')::interval) WHERE rocket_locks.locked_until < now()
RETURNING floor(extract(epoch from locked_until))::bigint;

-- name: RocketLockDelete :exec
DELETE FROM rocket_locks where Key = $1 AND floor(extract(epoch from locked_until)) = sqlc.narg('ts')::bigint;

------------------------------------------------------------------------------------------------------------------------

-- name: SavePerson :exec
INSERT INTO people (id, created_at, updated_at, firstname, surname, patronymic)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetPeopleCount :one
SELECT count(1) as count FROM people;

-- name: GetPerson :one
SELECT id, created_at, updated_at, firstname, surname, patronymic FROM people WHERE id = $1;

------------------------------------------------------------------------------------------------------------------------

-- name: SavePersonAuth :exec
INSERT INTO auth (person_id, created_at, updated_at, email, password_hash, status, role)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: EmailExists :one
SELECT count(1) as count FROM auth WHERE email = $1;

-- name: GetAuth :one
SELECT person_id, created_at, updated_at, email, password_hash, status, role FROM auth WHERE person_id = $1;

-- name: GetAuthByEmail :one
SELECT person_id, created_at, updated_at, email, password_hash, status, role FROM auth WHERE email = $1;

------------------------------------------------------------------------------------------------------------------------

-- name: SaveConfirmCode :exec
INSERT INTO codes (code, person_id, created_at, updated_at, active, type)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetActiveConfirmCode :one
SELECT code, person_id, created_at, updated_at, active, type FROM codes
WHERE code = $1 AND type = $2 and active = true;

------------------------------------------------------------------------------------------------------------------------

-- name: SaveRefreshToken :exec
INSERT INTO refresh_codes (id, person_id, created_at, updated_at, status)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateRefreshTokenStatus :exec
UPDATE refresh_codes
SET status = $1, updated_at = $2
WHERE id = $3;

-- name: GetLastActiveRefreshToken :one
SELECT id, person_id, created_at, updated_at, status FROM refresh_codes
WHERE id=$1 and status='ACTIVE'
ORDER BY created_at DESC
LIMIT 1;

------------------------------------------------------------------------------------------------------------------------

-- name: GetApiTokens :many
SELECT id, person_id, caption, token, created_at, updated_at, expired_at, type FROM api_tokens
WHERE person_id=$1
ORDER BY created_at;

-- name: SaveApiToken :exec
INSERT INTO api_tokens (id, person_id, caption, token, created_at, updated_at, expired_at, type)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: DeleteApiToken :one
DELETE FROM api_tokens
WHERE id=$1 and person_id=$2
RETURNING id;

-- name: GetApiToken :one
SELECT id, person_id, caption, token, created_at, updated_at, expired_at, type FROM api_tokens
WHERE token=$1
LIMIT 1;