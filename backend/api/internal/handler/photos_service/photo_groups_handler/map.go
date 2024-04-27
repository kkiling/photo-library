package photo_groups_handler

import (
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/form"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func mapPhotoPreview(item *model.PhotoPreviewDTO) *desc.PhotoPreview {
	return &desc.PhotoPreview{
		Src:    item.Src,
		Width:  int32(item.Width),
		Height: int32(item.Height),
		Size:   int32(item.Size),
	}
}

func mapPhotoPreviews(items []model.PhotoPreviewDTO) []*desc.PhotoPreview {
	return lo.Map(items, func(item model.PhotoPreviewDTO, index int) *desc.PhotoPreview {
		return mapPhotoPreview(&item)
	})
}

func mapPhotoWithPreviews(item *model.PhotoWithPreviewsDTO) *desc.PhotoWithPreviews {
	return &desc.PhotoWithPreviews{
		PhotoId:  item.PhotoID.String(),
		Original: mapPhotoPreview(&item.Original),
		Previews: mapPhotoPreviews(item.Previews),
	}
}

func mapPhotosWithPreviews(items []model.PhotoWithPreviewsDTO) []*desc.PhotoWithPreviews {
	return lo.Map(items, func(item model.PhotoWithPreviewsDTO, index int) *desc.PhotoWithPreviews {
		return mapPhotoWithPreviews(&item)
	})
}

func mapPhotoGroup(response *model.PhotoGroupDTO) *desc.PhotoGroup {
	return &desc.PhotoGroup{
		GroupId:     response.GroupID.String(),
		MainPhoto:   mapPhotoWithPreviews(&response.MainPhoto),
		PhotosCount: int32(response.PhotosCount),
		Photos:      mapPhotosWithPreviews(response.Photos),
	}
}

func mapGetPhotoGroupsResponse(response *model.PaginatedPhotoGroupsDTO) *desc.GetPhotoGroupsResponse {
	return &desc.GetPhotoGroupsResponse{
		Items: lo.Map(response.Items, func(item model.PhotoGroupDTO, index int) *desc.PhotoGroup {
			return mapPhotoGroup(&item)
		}),
		TotalItems: int32(response.TotalItems),
	}
}

func mapGetPhotoGroupsRequest(request *desc.GetPhotoGroupsRequest) form.GetPhotoGroups {
	return form.GetPhotoGroups{
		Paginator: model.Pagination{
			Page:    uint64(request.Page),
			PerPage: uint64(request.PerPage),
		},
	}
}
