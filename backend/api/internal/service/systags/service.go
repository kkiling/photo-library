package systags

import (
	"context"
	"errors"
	"fmt"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/tagphoto"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"unicode/utf8"
)

var ErrMetaNotFound = fmt.Errorf("meta not found")
var ErrUploadDataNotFound = fmt.Errorf("upload data not found")

const (
	YearTag              = "year"
	YearTagColor         = "#14ff00"
	CameraModelTag       = "camera_model"
	CameraModelTagColor  = "#ffed00"
	PhotoCatalogTag      = "photo_catalog"
	PhotoCatalogTagColor = "#ff0000"
	LocationTag          = "location"
	LocationTagColor     = "#00ecff"
)

type Database interface {
	service.Transactor
	GetMetaData(ctx context.Context, photoID uuid.UUID) (*model.MetaData, error)
	GetUploadPhotoData(ctx context.Context, photoID uuid.UUID) (*model.UploadPhotoData, error)
}

type TagPhoto interface {
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error)
	GetCategory(ctx context.Context, typeCategory string) (*model.TagCategory, error)
	CreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error)
}

type Service struct {
	logger     log.Logger
	tagService TagPhoto
	database   Database
	geocoder   geo.Geocoder
}

func NewService(logger log.Logger, tagService TagPhoto, database Database) *Service {
	return &Service{
		logger:     logger,
		tagService: tagService,
		database:   database,
		geocoder:   openstreetmap.Geocoder(),
	}
}

func (s *Service) getOrCreateTagCategory(ctx context.Context, tagCategory, color string) (model.TagCategory, error) {
	// TODO: кешировать

	findCategory, err := s.tagService.GetCategory(ctx, tagCategory)
	if err != nil {
		return model.TagCategory{}, fmt.Errorf("tagService.GetCategory: %w", err)
	}

	if findCategory != nil {
		return *findCategory, nil
	}

	category, err := s.tagService.CreateCategory(ctx, tagCategory, color)
	if err != nil {
		return model.TagCategory{}, fmt.Errorf("tagService.CreateCategory: %w", err)
	}

	return category, nil
}

func (s *Service) CreateTagByMeta(ctx context.Context, photo model.Photo) error {

	// По каталогу КАТАЛОГ
	data, err := s.database.GetUploadPhotoData(ctx, photo.ID)
	if err != nil {
		return fmt.Errorf("database.GetUploadPhotoData: %w", err)
	}
	if data == nil {
		return ErrUploadDataNotFound
	}

	tags := make(map[string]struct{})
	for _, path := range data.Paths {
		dirs := getDirectories(path)
		for _, dir := range dirs {
			tags[dir] = struct{}{}
		}
	}

	photoCatalog, err := s.getOrCreateTagCategory(ctx, PhotoCatalogTag, PhotoCatalogTagColor)
	if err != nil {
		return fmt.Errorf("getOrCreateTagCategory: %w", err)
	}

	for tag, _ := range tags {
		if utf8.RuneCountInString(tag) < tagphoto.TagNameMin {
			continue
		}
		_, err = s.tagService.AddPhotoTag(ctx, photo.ID, photoCatalog.ID, tag)
		if err != nil {
			if errors.Is(err, tagphoto.ErrTagAlreadyExist) {
				continue
			}
			return fmt.Errorf("tagService.AddPhotoTag: %w", err)
		}
	}

	// Мета информация
	meta, err := s.database.GetMetaData(ctx, photo.ID)
	if err != nil {
		return fmt.Errorf("databases.GetMetaData: %w", err)
	}

	if meta == nil {
		return ErrMetaNotFound
	}

	// По дате формируем тег ГОД
	if meta.DateTime != nil {
		yearCategory, err := s.getOrCreateTagCategory(ctx, YearTag, YearTagColor)
		if err != nil {
			return fmt.Errorf("getOrCreateTagCategory: %w", err)
		}

		name := fmt.Sprintf("%d", meta.DateTime.Year())
		_, err = s.tagService.AddPhotoTag(ctx, photo.ID, yearCategory.ID, name)
		if err != nil {
			if errors.Is(err, tagphoto.ErrTagAlreadyExist) {
			} else {
				return fmt.Errorf("tagService.AddPhotoTag: %w", err)
			}
		}
	}

	// По модели МОДЕЛЬ
	if meta.ModelInfo != nil && (*meta.ModelInfo != "") {
		cameraCategory, err := s.getOrCreateTagCategory(ctx, CameraModelTag, CameraModelTagColor)
		if err != nil {
			return fmt.Errorf("getOrCreateTagCategory: %w", err)
		}

		_, err = s.tagService.AddPhotoTag(ctx, photo.ID, cameraCategory.ID, *meta.ModelInfo)
		if err != nil {
			if errors.Is(err, tagphoto.ErrTagAlreadyExist) {
			} else {
				return fmt.Errorf("tagService.AddPhotoTag: %w", err)
			}
		}
	}

	if meta.Geo != nil {
		locationCategory, err := s.getOrCreateTagCategory(ctx, LocationTag, LocationTagColor)
		if err != nil {
			return fmt.Errorf("getOrCreateTagCategory: %w", err)
		}

		location, err := s.geocoder.ReverseGeocode(meta.Geo.Latitude, meta.Geo.Longitude)
		if err != nil {
			return fmt.Errorf("geocoder.ReverseGeocode: %w", err)
		}

		locationTags := make([]string, 0)
		if location.State != "" {
			locationTags = append(locationTags, location.State)
		}
		if location.StateDistrict != "" {
			locationTags = append(locationTags, location.StateDistrict)
		}
		if location.County != "" {
			locationTags = append(locationTags, location.County)
		}
		if location.Country != "" {
			locationTags = append(locationTags, location.Country)
		}
		if location.City != "" {
			locationTags = append(locationTags, location.City)
		}

		for _, loc := range locationTags {
			_, err = s.tagService.AddPhotoTag(ctx, photo.ID, locationCategory.ID, loc)
			if err != nil {
				if errors.Is(err, tagphoto.ErrTagAlreadyExist) {
				} else {
					return fmt.Errorf("tagService.AddPhotoTag: %w", err)
				}
			}
		}
	}

	return nil
}
