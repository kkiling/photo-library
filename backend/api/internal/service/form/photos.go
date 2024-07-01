package form

import "github.com/kkiling/photo-library/backend/api/internal/service/model"

// GetPhotoGroups запрос на получение групп фото
type GetPhotoGroups struct {
	Paginator model.Pagination
}
