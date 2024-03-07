package mapper

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func GetPhotoGroupsRequest(request *desc.GetPhotoGroupsRequest) *photos.GetPhotoGroupsRequest {
	return &photos.GetPhotoGroupsRequest{
		Paginator: model.Pagination{
			Page:    uint64(request.Page),
			PerPage: uint64(request.PerPage),
		},
	}
}

func GetPhotoGroupResponse(response *photos.GetPhotoGroupsResponse) *desc.GetPhotoGroupsResponse {

	items := make([]*desc.PhotoGroup, 0, len(response.Items))
	for _, group := range response.Items {
		items = append(items, &desc.PhotoGroup{
			Id: group.ID.String(),
			MainPhoto: &desc.Photo{
				Id:  group.MainPhoto.ID.String(),
				Url: group.MainPhoto.Url,
			},
			PhotosCount: int32(group.PhotoCount),
		})
	}

	return &desc.GetPhotoGroupsResponse{
		Items:      items,
		TotalItems: int32(response.TotalItems),
	}
}
