package tagphoto

import (
	"context"
	"errors"
	"fmt"
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

var ErrCategoryAlreadyExist = fmt.Errorf("category already exist")
var ErrCategoryNotExist = fmt.Errorf("category not exist")
var ErrTagAlreadyExist = fmt.Errorf("tag already exist")

type Database interface {
	service.Transactor
	GetTagCategory(ctx context.Context, categoryID uuid.UUID) (*model.TagCategory, error)
	GetTagCategoryByType(ctx context.Context, typeCategory string) (*model.TagCategory, error)
	SaveTagCategory(ctx context.Context, category model.TagCategory) error
	GetTagByName(ctx context.Context, photoID uuid.UUID, name string) (*model.Tag, error)
	SaveTag(ctx context.Context, tag model.Tag) error
}

type Service struct {
	database Database
}

func NewService(storage Database) *Service {
	return &Service{
		database: storage,
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

	if findCategory, err := s.database.GetTagCategoryByType(ctx, typeCategory); err != nil {
		// TODO: ошибка
		return model.TagCategory{}, fmt.Errorf("database.GetTagCategoryByName: %w", err)
	} else if findCategory != nil {
		// TODO: ошибка
		return model.TagCategory{}, ErrCategoryAlreadyExist
	}

	newCategory := model.TagCategory{
		ID:    uuid.New(),
		Type:  typeCategory,
		Color: color,
	}

	if err := s.database.SaveTagCategory(ctx, newCategory); err != nil {
		// TODO: ошибка
		return model.TagCategory{}, fmt.Errorf("database.SaveTagCategory: %w", err)
	}

	return newCategory, nil
}

func (s *Service) GetCategory(ctx context.Context, typeCategory string) (*model.TagCategory, error) {
	if findCategory, err := s.database.GetTagCategoryByType(ctx, typeCategory); err != nil {
		// TODO: ошибка
		return nil, fmt.Errorf("database.GetTagCategoryByName: %w", err)
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
		// TODO: ошибка

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

	if findCategory, err := s.database.GetTagCategory(ctx, categoryID); err != nil {
		// TODO: ошибка
		return model.Tag{}, fmt.Errorf("database.GetTypeCategory: %w", err)
	} else if findCategory == nil {
		// TODO: ошибка
		return model.Tag{}, ErrCategoryNotExist
	}

	if findTag, err := s.database.GetTagByName(ctx, photoID, name); err != nil {
		// TODO: ошибка
		return model.Tag{}, fmt.Errorf("database.GetTagByName: %w", err)
	} else if findTag != nil {
		// TODO: ошибка
		return model.Tag{}, ErrTagAlreadyExist
	}

	tag := model.Tag{
		ID:         uuid.New(),
		CategoryID: categoryID,
		PhotoID:    photoID,
		Name:       name,
	}

	if err := s.database.SaveTag(ctx, tag); err != nil {
		// TODO: ошибка
		return model.Tag{}, fmt.Errorf("database.SaveTag: %w", err)
	}

	return tag, nil
}
