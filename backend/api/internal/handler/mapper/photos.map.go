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

func GetPhotoGroupsResponse(response *photos.GetPhotoGroupsResponse) *desc.GetPhotoGroupsResponse {

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

func MetaData(metadata *model.PhotoMetadata) *desc.Metadata {
	var geo *desc.Geo
	if metadata.Geo != nil {
		geo = &desc.Geo{
			Latitude:  metadata.Geo.Latitude,
			Longitude: metadata.Geo.Longitude,
		}
	}
	return &desc.Metadata{
		ModelInfo:   metadata.ModelInfo,
		SizeBytes:   int32(metadata.SizeBytes),
		WidthPixel:  int32(metadata.WidthPixel),
		HeightPixel: int32(metadata.HeightPixel),
		DataTime:    toTimestampPtr(metadata.DateTime),
		UpdateAt:    toTimestamp(metadata.UpdateAt),
		Geo:         geo,
	}
}

func Tag(tag *photos.Tag) *desc.Tag {
	return &desc.Tag{
		Id:    tag.ID.String(),
		Name:  tag.Name,
		Type:  tag.Type,
		Color: tag.Color,
	}
}

func PhotoWithData(photo *photos.PhotoWithData) *desc.PhotoWithData {
	tags := make([]*desc.Tag, 0, len(photo.Tags))
	for _, tag := range photo.Tags {
		tags = append(tags, Tag(&tag))
	}
	return &desc.PhotoWithData{
		Id:       photo.ID.String(),
		Url:      photo.Url,
		MetaData: MetaData(photo.Metadata),
		Tags:     tags,
	}
}

func GetPhotoGroupResponse(response *photos.PhotoGroupData) *desc.GetPhotoGroupResponse {
	photoArray := make([]*desc.PhotoWithData, 0, len(response.Photos))
	for _, photo := range response.Photos {
		photoArray = append(photoArray, PhotoWithData(&photo))
	}
	return &desc.GetPhotoGroupResponse{
		Id:        response.ID.String(),
		MainPhoto: PhotoWithData(&response.MainPhoto),
		Photos:    photoArray,
	}
}
