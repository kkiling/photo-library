package photos

import (
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type Photo struct {
	ID  uuid.UUID
	Url string
}

type PhotoWithData struct {
	Photo
	Metadata *model.PhotoMetadata
	Tags     []Tag
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

type Tag struct {
	ID    uuid.UUID
	Name  string
	Type  string
	Color string
}

type PhotoGroupData struct {
	ID        uuid.UUID
	MainPhoto PhotoWithData
	Photos    []PhotoWithData
}
