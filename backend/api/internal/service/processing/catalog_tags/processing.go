package catalog_tags

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/tags"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const (
	PhotoCatalogTag      = "photo_catalog"
	PhotoCatalogTagColor = "#ff0000"
)

type Storage interface {
	service.Transactor
	GetPhotoUploadData(ctx context.Context, photoID uuid.UUID) (model.PhotoUploadData, error)
}

type TagPhoto interface {
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error)
	GetOrCreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error)
	DeletePhotoTagsByCategories(ctx context.Context, photoID uuid.UUID, categoryIDs []uuid.UUID) error
}

type Processing struct {
	logger     log.Logger
	tagService TagPhoto
	storage    Storage
}

func NewService(logger log.Logger, tagService TagPhoto, storage Storage) *Processing {
	return &Processing{
		logger:     logger,
		tagService: tagService,
		storage:    storage,
	}
}

func (s *Processing) createPhotoCatalogTag(ctx context.Context, photo model.Photo, uploadData *model.PhotoUploadData) error {
	tagsMap := make(map[string]struct{})
	for _, path := range uploadData.Paths {
		dirs := getDirectories(path)
		for _, dir := range dirs {
			tagsMap[dir] = struct{}{}
		}
	}

	photoCatalog, err := s.tagService.GetOrCreateCategory(ctx, PhotoCatalogTag, PhotoCatalogTagColor)
	if err != nil {
		return serviceerr.MakeErr(err, "getOrCreateTagCategory")
	}

	for tag, _ := range tagsMap {
		if utf8.RuneCountInString(tag) < tags.TagNameMin {
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

func (s *Processing) Compensate(ctx context.Context, photoID uuid.UUID) error {
	photoCatalog, err := s.tagService.GetOrCreateCategory(ctx, PhotoCatalogTag, PhotoCatalogTagColor)
	if err != nil {
		return serviceerr.MakeErr(err, "getOrCreateTagCategory")
	}

	err = s.tagService.DeletePhotoTagsByCategories(ctx, photoID, []uuid.UUID{photoCatalog.ID})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, serviceerr.ErrNotFound):
		return nil
	default:
		return serviceerr.MakeErr(err, "s.tagService.DeletePhotoTagsByCategories")
	}
}

func (s *Processing) Init(_ context.Context) error {
	return nil
}

func (s *Processing) NeedLoadPhotoBody() bool {
	return false
}

// Processing создание и сохранение автоматических тегов (по метаданным или по путям и тд)
func (s *Processing) Processing(ctx context.Context, photo model.Photo, _ []byte) (bool, error) {
	uploadData, err := s.storage.GetPhotoUploadData(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.GetPhotoUploadData")
	}

	// Теги По каталогу По каталогу КАТАЛОГ
	if err = s.createPhotoCatalogTag(ctx, photo, &uploadData); err != nil {
		return false, serviceerr.MakeErr(err, "s.createPhotoCatalogTag")
	}

	return true, nil
}
