package tags

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/service/utils"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

const (
	TagCategoryTypeMin = 3
	TagCategoryTypeMax = 128
	TagNameMin         = 3
	TagNameMax         = 128
)

type Storage interface {
	service.Transactor
	GetTagCategory(ctx context.Context, id uuid.UUID) (model.TagCategory, error)
	GetTagCategoryByType(ctx context.Context, typeCategory string) (model.TagCategory, error)
	SaveTagCategory(ctx context.Context, category model.TagCategory) error
	SaveTag(ctx context.Context, tag model.Tag) error
	GetTags(ctx context.Context, photoID uuid.UUID) ([]model.Tag, error)
	DeletePhotoTagsByCategories(ctx context.Context, photoID uuid.UUID, categoryID []uuid.UUID) error
	DeleteTag(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	storage          Storage
	categoriesByType map[string]model.TagCategory
	categoriesByID   map[uuid.UUID]model.TagCategory
	mu               sync.Mutex
}

func NewService(storage Storage) *Service {
	return &Service{
		storage:          storage,
		categoriesByType: make(map[string]model.TagCategory),
		categoriesByID:   make(map[uuid.UUID]model.TagCategory),
	}
}

func validateCreateCategory(typeCategory, color string) error {
	validate := utils.NewValidator()

	// Валидация имени
	if err := validate.Var(typeCategory, fmt.Sprintf("min=%d,max=%d", TagCategoryTypeMin, TagCategoryTypeMax)); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return serviceerr.MakeErr(validationErrors, "invalid typeCategory")
		}

		return serviceerr.MakeErr(err, "invalid typeCategory")
	}

	// Валидация цвета в формате HEX (например, "#FFFFFF")
	// Вы можете настроить этот шаблон, если у вас есть другие требования к формату.
	if err := validate.Var(color, "hexcolor"); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return serviceerr.MakeErr(validationErrors, "invalid color")
		}

		return serviceerr.MakeErr(err, "invalid color")
	}

	return nil
}

func (s *Service) GetOrCreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error) {
	// Категории создаются много из каких мест, может быть гонка
	s.mu.Lock()
	defer s.mu.Unlock()

	if find, ok := s.categoriesByType[typeCategory]; ok {
		return find, nil
	}

	findCategory, err := s.storage.GetTagCategoryByType(ctx, typeCategory)
	if err == nil {
		s.categoriesByType[typeCategory] = findCategory
		return findCategory, nil
	}

	if !errors.Is(err, serviceerr.ErrNotFound) {
		return model.TagCategory{}, serviceerr.MakeErr(err, "s.storage.SaveTagCategory")
	}

	if err = validateCreateCategory(typeCategory, color); err != nil {
		return model.TagCategory{}, serviceerr.InvalidInputErr(err, "validateCreateCategory")
	}

	newCategory := model.TagCategory{
		ID:    uuid.New(),
		Type:  typeCategory,
		Color: color,
	}

	if err = s.storage.SaveTagCategory(ctx, newCategory); err != nil {
		return model.TagCategory{}, serviceerr.MakeErr(err, "storage.SaveTagCategory")
	}

	s.categoriesByType[typeCategory] = newCategory

	return newCategory, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (model.TagCategory, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if find, ok := s.categoriesByID[categoryID]; ok {
		return find, nil
	}
	if findCategory, err := s.storage.GetTagCategory(ctx, categoryID); err != nil {
		if errors.Is(err, serviceerr.ErrNotFound) {
			return model.TagCategory{}, serviceerr.NotFoundf("category not found")
		}
		return model.TagCategory{}, serviceerr.MakeErr(err, "storage.GetTagCategory")
	} else {
		return findCategory, nil
	}
}

func (s *Service) GetCategories(ctx context.Context) ([]model.TagCategory, error) {
	panic("not implemented")
}

func (s *Service) GetTags(ctx context.Context, photoID uuid.UUID) ([]model.Tag, error) {
	tags, err := s.storage.GetTags(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "storage.GetTagCategoryByName")
	}
	return tags, nil
}

func validateAddPhotoTag(name string) error {
	validate := utils.NewValidator()

	// Валидация имени
	if err := validate.Var(name, fmt.Sprintf("min=%d,max=%d", TagNameMin, TagNameMax)); err != nil {
		return serviceerr.MakeErr(err, "invalid name")
	}

	return nil
}

// AddPhotoTag добавляет тег для фотографии
func (s *Service) AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error) {
	if err := validateAddPhotoTag(name); err != nil {
		return model.Tag{}, serviceerr.InvalidInputErr(err, "validateAddPhotoTag")
	}

	tag := model.Tag{
		ID:         uuid.New(),
		CategoryID: categoryID,
		PhotoID:    photoID,
		Name:       name,
	}

	if err := s.storage.SaveTag(ctx, tag); err != nil {
		return model.Tag{}, serviceerr.MakeErr(err, "storage.SaveTag")
	}

	return tag, nil
}

// DeletePhotoTagsByCategories удаление тегов фотографии по категории
func (s *Service) DeletePhotoTagsByCategories(ctx context.Context, photoID uuid.UUID, categoryIDs []uuid.UUID) error {
	if err := s.storage.DeletePhotoTagsByCategories(ctx, photoID, categoryIDs); err != nil {
		return serviceerr.MakeErr(err, "storage.DeletePhotoTagsByCategories")
	}
	return nil
}

func (s *Service) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return s.DeleteTag(ctx, id)
}
