package photos

import (
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type PhotoPreview struct {
	Src    string
	Width  int
	Height int
	Size   int
}

type Photo struct {
	ID       uuid.UUID
	Src      string
	Width    int
	Height   int
	Size     int
	Metadata *model.PhotoMetadata
	Tags     []Tag
	Preview  []PhotoPreview
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
	PhotoBody          []byte
	Extension          model.PhotoExtension
	WidthPixelCurrent  int
	HeightPixelCurrent int
}

type Tag struct {
	ID    uuid.UUID
	Name  string
	Type  string
	Color string
}

type PhotoGroupData struct {
	ID        uuid.UUID
	MainPhoto Photo
	Photos    []Photo
}
