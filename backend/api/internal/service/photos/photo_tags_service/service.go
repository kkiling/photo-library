package photo_tags_service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Storage interface {
	service.Transactor
}

type TagPhoto interface {
	GetTags(ctx context.Context, photoID uuid.UUID) ([]model.Tag, error)
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (model.TagCategory, error)
	GetCategories(ctx context.Context) ([]model.TagCategory, error)
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error)
	DeleteTag(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	logger     log.Logger
	storage    Storage
	tagService TagPhoto
}

func NewService(logger log.Logger, tagService TagPhoto, storage Storage) *Service {
	return &Service{
		logger:     logger,
		storage:    storage,
		tagService: tagService,
	}
}

func (s *Service) GetPhotoTags(ctx context.Context, photoID uuid.UUID) ([]model.TagWithCategoryDTO, error) {

	tags, err := s.tagService.GetTags(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.tagService.GetTags")
	}

	tagsWithCategories := make([]model.TagWithCategoryDTO, 0, len(tags))

	// Проще так, потому что категории кешируются
	for _, tag := range tags {
		category, err := s.tagService.GetCategoryByID(ctx, tag.CategoryID)
		if err != nil {
			return nil, serviceerr.MakeErr(err, "s.tagService.GetCategoryByID")
		}
		tagsWithCategories = append(tagsWithCategories, model.TagWithCategoryDTO{
			ID:         tag.ID,
			Name:       tag.Name,
			IDCategory: category.ID,
			Type:       category.Type,
			Color:      category.Color,
		})
	}

	return tagsWithCategories, nil
}

func (s *Service) GetTagCategories(ctx context.Context) ([]model.TagCategory, error) {
	categories, err := s.tagService.GetCategories(ctx)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.tagService.GetCategories")
	}
	return categories, nil
}

func (s *Service) AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, tagName string) error {
	_, err := s.tagService.AddPhotoTag(ctx, photoID, categoryID, tagName)
	if err != nil {
		return serviceerr.MakeErr(err, "s.tagService.AddPhotoTag")
	}
	return nil
}

func (s *Service) DeletePhotoTag(ctx context.Context, tagID uuid.UUID) error {
	err := s.tagService.DeleteTag(ctx, tagID)
	if err != nil {
		return serviceerr.MakeErr(err, "s.tagService.DeleteTag")
	}
	return nil
}
