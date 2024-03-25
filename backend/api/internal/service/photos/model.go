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

type PhotoGroup struct {
	ID          uuid.UUID
	Original    PhotoPreview
	Previews    []PhotoPreview
	PhotosCount int
}

type GetPhotoGroupsRequest struct {
	Paginator model.Pagination
}

type PaginatedPhotoGroups struct {
	Items      []PhotoGroup
	TotalItems int
}

type PhotoGroupData struct {
	ID          uuid.UUID
	Original    PhotoPreview
	Previews    []PhotoPreview
	PhotosCount int
	Metadata    *model.PhotoMetadata
	Tags        []Tag
}
