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

func GetPhotoGroupsResponse(response *photos.PaginatedPhotoGroups) *desc.PaginatedPhotoGroups {
	items := make([]*desc.PhotoGroup, 0, len(response.Items))
	for _, group := range response.Items {
		items = append(items, &desc.PhotoGroup{
			Id:          group.ID.String(),
			MainPhoto:   PhotoWithPreviews(&group.MainPhoto),
			PhotosCount: int32(group.PhotosCount),
		})
	}

	return &desc.PaginatedPhotoGroups{
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

func Tag(tag *photos.Tag) *desc.PhotoTag {
	return &desc.PhotoTag{
		Id:    tag.ID.String(),
		Name:  tag.Name,
		Type:  tag.Type,
		Color: tag.Color,
	}
}

func PhotoWithPreviews(photo *photos.PhotoWithPreviews) *desc.PhotoWithPreviews {
	return &desc.PhotoWithPreviews{
		Id:       photo.ID.String(),
		Src:      photo.Src,
		Width:    int32(photo.Width),
		Height:   int32(photo.Height),
		Previews: PhotoPreviews(photo.Previews),
	}
}

func PhotosWithPreview(photos []photos.PhotoWithPreviews) []*desc.PhotoWithPreviews {
	previews := make([]*desc.PhotoWithPreviews, 0, len(photos))
	for _, pr := range photos {
		previews = append(previews, PhotoWithPreviews(&pr))
	}
	return previews
}

func PhotoPreview(preview *photos.PhotoPreview) *desc.PhotoPreview {
	return &desc.PhotoPreview{
		Src:    preview.Src,
		Width:  int32(preview.Width),
		Height: int32(preview.Height),
		Size:   int32(preview.Size),
	}
}

func PhotoPreviews(preview []photos.PhotoPreview) []*desc.PhotoPreview {
	previews := make([]*desc.PhotoPreview, 0, len(preview))
	for _, pr := range preview {
		previews = append(previews, PhotoPreview(&pr))
	}
	return previews
}

func GetPhotoGroupResponse(response *photos.PhotoGroupData) *desc.PhotoGroupData {
	tags := make([]*desc.PhotoTag, 0, len(response.Tags))
	for _, tag := range response.Tags {
		tags = append(tags, Tag(&tag))
	}

	return &desc.PhotoGroupData{
		Id:          response.ID.String(),
		MainPhoto:   PhotoWithPreviews(&response.MainPhoto),
		PhotosCount: int32(response.PhotosCount),
		Photos:      PhotosWithPreview(response.Photos),
		MetaData:    MetaData(response.Metadata),
		Tags:        tags,
	}
}
