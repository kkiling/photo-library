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
	GetTags(ctx context.Context, photoID uuid.UUID) ([]model.Tag, error)
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

// CreateCategory создание категории тегов
func (s *Service) CreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error) {
	if err := validateCreateCategory(typeCategory, color); err != nil {
		return model.TagCategory{}, serviceerr.InvalidInputErr(err, "validateCreateCategory")
	}

	if findCategory, err := s.storage.GetTagCategoryByType(ctx, typeCategory); err != nil {
		return model.TagCategory{}, serviceerr.MakeErr(err, "storage.GetTagCategoryByName")
	} else if findCategory != nil {
		return model.TagCategory{}, serviceerr.ConflictError("category already exist")
	}

	newCategory := model.TagCategory{
		ID:    uuid.New(),
		Type:  typeCategory,
		Color: color,
	}

	if err := s.storage.SaveTagCategory(ctx, newCategory); err != nil {
		return model.TagCategory{}, serviceerr.MakeErr(err, "storage.SaveTagCategory")
	}

	return newCategory, nil
}

func (s *Service) GetCategory(ctx context.Context, typeCategory string) (*model.TagCategory, error) {
	// TODO: кешировать
	if findCategory, err := s.storage.GetTagCategoryByType(ctx, typeCategory); err != nil {
		return nil, serviceerr.MakeErr(err, "storage.GetTagCategoryByName")
	} else if findCategory != nil {
		return findCategory, nil
	}

	return nil, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*model.TagCategory, error) {
	// TODO: кешировать
	if findCategory, err := s.storage.GetTagCategory(ctx, categoryID); err != nil {
		return nil, serviceerr.MakeErr(err, "storage.GetTagCategory")
	} else if findCategory != nil {
		return findCategory, nil
	}

	return nil, nil
}

func (s *Service) GetTags(ctx context.Context, photoID uuid.UUID) ([]model.Tag, error) {
	tags, err := s.storage.GetTags(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "storage.GetTagCategoryByName")
	}
	return tags, nil
}

func validateAddPhotoTag(name string) error {
	validate := validator.New()

	// Валидация имени
	if err := validate.Var(name, fmt.Sprintf("min=%d,max=%d", TagNameMin, TagNameMax)); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return serviceerr.MakeErr(validationErrors, "invalid name")
		}

		return serviceerr.MakeErr(validationErrors, "invalid name")
	}

	return nil
}

// AddPhotoTag добавляет тег для фотографии
func (s *Service) AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error) {
	if err := validateAddPhotoTag(name); err != nil {
		return model.Tag{}, serviceerr.InvalidInputError("validateAddPhotoTag", err)
	}

	if findCategory, err := s.storage.GetTagCategory(ctx, categoryID); err != nil {
		return model.Tag{}, serviceerr.MakeErr(err, "storage.GetTypeCategory")
	} else if findCategory == nil {
		return model.Tag{}, serviceerr.NotFoundError("category not exist")
	}

	if findTag, err := s.storage.GetTagByName(ctx, photoID, name); err != nil {
		return model.Tag{}, serviceerr.MakeErr(err, "storage.GetTagByName")
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
		return model.Tag{}, serviceerr.MakeErr(err, "storage.SaveTag")
	}

	return tag, nil
}
