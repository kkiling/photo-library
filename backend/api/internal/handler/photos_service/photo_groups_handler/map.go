package photo_groups_handler

import (
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos/photo_groups_service"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func mapPhotoPreview(item *photo_groups_service.PhotoPreview) *desc.PhotoPreview {
	return &desc.PhotoPreview{
		Src:    item.Src,
		Width:  int32(item.Width),
		Height: int32(item.Height),
		Size:   int32(item.Size),
	}
}

func mapPhotoPreviews(items []photo_groups_service.PhotoPreview) []*desc.PhotoPreview {
	return lo.Map(items, func(item photo_groups_service.PhotoPreview, index int) *desc.PhotoPreview {
		return mapPhotoPreview(&item)
	})
}

func mapPhotoWithPreviews(item *photo_groups_service.PhotoWithPreviews) *desc.PhotoWithPreviews {
	return &desc.PhotoWithPreviews{
		Id:       item.ID.String(),
		Original: mapPhotoPreview(&item.Original),
		Previews: mapPhotoPreviews(item.Previews),
	}
}

func mapPhotosWithPreviews(items []photo_groups_service.PhotoWithPreviews) []*desc.PhotoWithPreviews {
	return lo.Map(items, func(item photo_groups_service.PhotoWithPreviews, index int) *desc.PhotoWithPreviews {
		return mapPhotoWithPreviews(&item)
	})
}

func mapPhotoGroup(response *photo_groups_service.PhotoGroup) *desc.PhotoGroup {
	return &desc.PhotoGroup{
		Id:          response.ID.String(),
		MainPhoto:   mapPhotoWithPreviews(&response.MainPhoto),
		PhotosCount: int32(response.PhotosCount),
		Photos:      mapPhotosWithPreviews(response.Photos),
	}
}

func mapGetPhotoGroupsResponse(response *photo_groups_service.PaginatedPhotoGroups) *desc.GetPhotoGroupsResponse {
	return &desc.GetPhotoGroupsResponse{
		Items: lo.Map(response.Items, func(item photo_groups_service.PhotoGroup, index int) *desc.PhotoGroup {
			return mapPhotoGroup(&item)
		}),
		TotalItems: int32(response.TotalItems),
	}
}

func mapGetPhotoGroupsRequest(request *desc.GetPhotoGroupsRequest) *photo_groups_service.GetPhotoGroupsRequest {
	return &photo_groups_service.GetPhotoGroupsRequest{
		Paginator: model.Pagination{
			Page:    uint64(request.Page),
			PerPage: uint64(request.PerPage),
		},
	}
}
