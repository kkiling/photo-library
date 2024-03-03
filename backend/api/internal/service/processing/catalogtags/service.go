package catalogtags

import (
	"context"
	"errors"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/tagphoto"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const (
	PhotoCatalogTag      = "photo_catalog"
	PhotoCatalogTagColor = "#ff0000"
)

type Storage interface {
	service.Transactor
	GetUploadPhotoData(ctx context.Context, photoID uuid.UUID) (*model.PhotoUploadData, error)
}

type TagPhoto interface {
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error)
	GetCategory(ctx context.Context, typeCategory string) (*model.TagCategory, error)
	CreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error)
}

type Service struct {
	logger     log.Logger
	tagService TagPhoto
	storage    Storage
}

func NewService(logger log.Logger, tagService TagPhoto, storage Storage) *Service {
	return &Service{
		logger:     logger,
		tagService: tagService,
		storage:    storage,
	}
}

func (s *Service) getOrCreateTagCategory(ctx context.Context, tagCategory, color string) (model.TagCategory, error) {
	// TODO: кешировать

	findCategory, err := s.tagService.GetCategory(ctx, tagCategory)
	if err != nil {
		return model.TagCategory{}, serviceerr.MakeErr(err, "tagService.GetCategory")
	}

	if findCategory != nil {
		return *findCategory, nil
	}

	category, err := s.tagService.CreateCategory(ctx, tagCategory, color)
	if err != nil {
		return model.TagCategory{}, serviceerr.MakeErr(err, "tagService.CreateCategory")
	}

	return category, nil
}

func (s *Service) createPhotoCatalogTag(ctx context.Context, photo model.Photo, uploadData *model.PhotoUploadData) error {
	tags := make(map[string]struct{})
	for _, path := range uploadData.Paths {
		dirs := getDirectories(path)
		for _, dir := range dirs {
			tags[dir] = struct{}{}
		}
	}

	photoCatalog, err := s.getOrCreateTagCategory(ctx, PhotoCatalogTag, PhotoCatalogTagColor)
	if err != nil {
		return serviceerr.MakeErr(err, "getOrCreateTagCategory")
	}

	for tag, _ := range tags {
		if utf8.RuneCountInString(tag) < tagphoto.TagNameMin {
			continue
		}
		_, err = s.tagService.AddPhotoTag(ctx, photo.ID, photoCatalog.ID, tag)
		if err != nil {
			if errors.Is(err, serviceerr.ErrTagAlreadyExist) {
				continue
			}
			return serviceerr.MakeErr(err, "tagService.AddPhotoTag")
		}
	}

	return nil
}

func (s *Service) Init(ctx context.Context) error {
	return nil
}

func (s *Service) NeedLoadPhotoBody() bool {
	return false
}

// Processing создание и сохранение автоматических тегов (по мета данным или по путям и тд)
func (s *Service) Processing(ctx context.Context, photo model.Photo, _ []byte) (bool, error) {
	uploadData, err := s.storage.GetUploadPhotoData(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.GetPhotoUploadData")
	}

	if uploadData == nil {
		return false, serviceerr.NotFoundError("upload data not found")
	}

	// Теги По каталогу По каталогу КАТАЛОГ
	if err := s.createPhotoCatalogTag(ctx, photo, uploadData); err != nil {
		return false, serviceerr.MakeErr(err, "s.createPhotoCatalogTag")
	}

	return true, nil
}
