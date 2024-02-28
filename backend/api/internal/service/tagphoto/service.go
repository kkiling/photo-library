package tagphoto

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

const (
	TagCategoryTypeMin = 3
	TagCategoryTypeMax = 128
	TagNameMin         = 3
	TagNameMax         = 128
)

type Storage interface {
	service.Transactor
	GetTagCategory(ctx context.Context, categoryID uuid.UUID) (*model.TagCategory, error)
	GetTagCategoryByType(ctx context.Context, typeCategory string) (*model.TagCategory, error)
	SaveTagCategory(ctx context.Context, category model.TagCategory) error
	GetTagByName(ctx context.Context, photoID uuid.UUID, name string) (*model.Tag, error)
	SaveTag(ctx context.Context, tag model.Tag) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func validateCreateCategory(typeCategory, color string) error {
	validate := validator.New()

	// Валидация имени
	if err := validate.Var(typeCategory, fmt.Sprintf("min=%d,max=%d", TagCategoryTypeMin, TagCategoryTypeMax)); err != nil {
		// TODO: ошибка

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return fmt.Errorf("invalid typeCategory: %w", validationErrors)
		}

		return fmt.Errorf("invalid typeCategory: %w", err)
	}

	// Валидация цвета в формате HEX (например, "#FFFFFF")
	// Вы можете настроить этот шаблон, если у вас есть другие требования к формату.
	if err := validate.Var(color, "hexcolor"); err != nil {
		// TODO: ошибка
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return fmt.Errorf("invalid color: %w", validationErrors)
		}

		return fmt.Errorf("invalid color: %w", err)
	}

	return nil
}

// CreateCategory создание категории тегов
func (s *Service) CreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error) {
	if err := validateCreateCategory(typeCategory, color); err != nil {
		return model.TagCategory{}, err
	}

	if findCategory, err := s.storage.GetTagCategoryByType(ctx, typeCategory); err != nil {
		return model.TagCategory{}, fmt.Errorf("storage.GetTagCategoryByName: %w", err)
	} else if findCategory != nil {
		return model.TagCategory{}, fmt.Errorf("category already exist")
	}

	newCategory := model.TagCategory{
		ID:    uuid.New(),
		Type:  typeCategory,
		Color: color,
	}

	if err := s.storage.SaveTagCategory(ctx, newCategory); err != nil {
		return model.TagCategory{}, fmt.Errorf("storage.SaveTagCategory: %w", err)
	}

	return newCategory, nil
}

func (s *Service) GetCategory(ctx context.Context, typeCategory string) (*model.TagCategory, error) {
	if findCategory, err := s.storage.GetTagCategoryByType(ctx, typeCategory); err != nil {
		return nil, fmt.Errorf("storage.GetTagCategoryByName: %w", err)
	} else if findCategory != nil {
		// TODO: ошибка
		return findCategory, nil
	}

	return nil, nil
}

func validateAddPhotoTag(name string) error {
	validate := validator.New()

	// Валидация имени
	if err := validate.Var(name, fmt.Sprintf("min=%d,max=%d", TagNameMin, TagNameMax)); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return fmt.Errorf("invalid name: %w", validationErrors)
		}

		return fmt.Errorf("invalid name: %w", err)
	}

	return nil
}

// AddPhotoTag добавляет тег для фотографии
func (s *Service) AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error) {
	if err := validateAddPhotoTag(name); err != nil {
		return model.Tag{}, err
	}

	if findCategory, err := s.storage.GetTagCategory(ctx, categoryID); err != nil {
		return model.Tag{}, fmt.Errorf("storage.GetTypeCategory: %w", err)
	} else if findCategory == nil {
		// TODO: ошибка
		return model.Tag{}, fmt.Errorf("category not exist")
	}

	if findTag, err := s.storage.GetTagByName(ctx, photoID, name); err != nil {
		return model.Tag{}, fmt.Errorf("storage.GetTagByName: %w", err)
	} else if findTag != nil {
		return model.Tag{}, serviceerr.ErrTagAlreadyExist
	}

	tag := model.Tag{
		ID:         uuid.New(),
		CategoryID: categoryID,
		PhotoID:    photoID,
		Name:       name,
	}

	if err := s.storage.SaveTag(ctx, tag); err != nil {
		// TODO: ошибка
		return model.Tag{}, fmt.Errorf("storage.SaveTag: %w", err)
	}

	return tag, nil
}
