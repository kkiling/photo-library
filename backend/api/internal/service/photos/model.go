package photos

import (
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type Photo struct {
	ID  uuid.UUID
	Url string
}

type PhotoGroup struct {
	ID         uuid.UUID
	MainPhoto  Photo
	PhotoCount int
}

type GetPhotoGroupsRequest struct {
	Paginator model.Pagination
}

type GetPhotoGroupsResponse struct {
	Items      []PhotoGroup
	TotalItems int
}

type PhotoContent struct {
	PhotoBody []byte
	Extension model.PhotoExtension
}
