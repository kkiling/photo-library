package photo_tags_handler

import (
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos/photo_tags_service"
	desc "github.com/kkiling/photo-library/backend/api/pkg/common/gen/proto/v1"
)

func mapTags(tags []photo_tags_service.TagWithCategory) []*desc.PhotoTag {
	return lo.Map(tags, func(tag photo_tags_service.TagWithCategory, index int) *desc.PhotoTag {
		return &desc.PhotoTag{
			Id:    tag.ID.String(),
			Name:  tag.Name,
			Type:  tag.Type,
			Color: tag.Color,
		}
	})
}

func mapGetTagCategories(response []model.TagCategory) *desc.GetTagCategoriesResponse {
	return &desc.GetTagCategoriesResponse{
		Items: lo.Map(response, func(tag model.TagCategory, _ int) *desc.TagCategory {
			return &desc.TagCategory{
				Id:    tag.ID.String(),
				Type:  tag.Type,
				Color: tag.Color,
			}
		}),
	}
}
