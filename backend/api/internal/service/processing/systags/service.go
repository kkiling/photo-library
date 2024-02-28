package systags

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/geo"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/tagphoto"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

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

type Storage interface {
	service.Transactor
	GetMetaData(ctx context.Context, photoID uuid.UUID) (*model.PhotoMetadata, error)
	GetUploadPhotoData(ctx context.Context, photoID uuid.UUID) (*model.PhotoUploadData, error)
}

type TagPhoto interface {
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error)
	GetCategory(ctx context.Context, typeCategory string) (*model.TagCategory, error)
	CreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error)
}

type GeoService interface {
	ReverseGeocode(ctx context.Context, lat, lng float64) (*geo.Address, error)
}

type Service struct {
	logger     log.Logger
	tagService TagPhoto
	storage    Storage
	geocoder   GeoService
}

func NewService(logger log.Logger, tagService TagPhoto, storage Storage, geoService GeoService) *Service {
	return &Service{
		logger:     logger,
		tagService: tagService,
		storage:    storage,
		geocoder:   geoService,
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
		return fmt.Errorf("getOrCreateTagCategory: %w", err)
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
			return fmt.Errorf("tagService.AddPhotoTag: %w", err)
		}
	}

	return nil
}

func (s *Service) createYearTag(ctx context.Context, photo model.Photo, metaData *model.PhotoMetadata) error {
	if metaData.DateTime == nil {
		return nil
	}
	yearCategory, err := s.getOrCreateTagCategory(ctx, YearTag, YearTagColor)
	if err != nil {
		return fmt.Errorf("getOrCreateTagCategory: %w", err)
	}

	name := fmt.Sprintf("%d", metaData.DateTime.Year())
	_, err = s.tagService.AddPhotoTag(ctx, photo.ID, yearCategory.ID, name)
	if err != nil {
		if errors.Is(err, serviceerr.ErrTagAlreadyExist) {
		} else {
			return fmt.Errorf("tagService.AddPhotoTag: %w", err)
		}
	}

	return nil
}

func (s *Service) createCameraModelTag(ctx context.Context, photo model.Photo, metaData *model.PhotoMetadata) error {
	if metaData.ModelInfo == nil || (*metaData.ModelInfo == "") {
		return nil
	}
	cameraCategory, err := s.getOrCreateTagCategory(ctx, CameraModelTag, CameraModelTagColor)
	if err != nil {
		return fmt.Errorf("getOrCreateTagCategory: %w", err)
	}

	_, err = s.tagService.AddPhotoTag(ctx, photo.ID, cameraCategory.ID, *metaData.ModelInfo)
	if err != nil {
		if errors.Is(err, serviceerr.ErrTagAlreadyExist) {
		} else {
			return fmt.Errorf("tagService.AddPhotoTag: %w", err)
		}
	}

	return nil
}

func (s *Service) createLocationTag(ctx context.Context, photo model.Photo, metaData *model.PhotoMetadata) error {
	if metaData.Geo == nil {
		return nil
	}
	locationCategory, err := s.getOrCreateTagCategory(ctx, LocationTag, LocationTagColor)
	if err != nil {
		return fmt.Errorf("getOrCreateTagCategory: %w", err)
	}

	location, err := s.geocoder.ReverseGeocode(ctx, metaData.Geo.Latitude, metaData.Geo.Longitude)
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
			if errors.Is(err, serviceerr.ErrTagAlreadyExist) {
			} else {
				return fmt.Errorf("tagService.AddPhotoTag: %w", err)
			}
		}
	}

	return nil
}

// Processing создание и сохранение автоматических тегов (по мета данным или по путям и тд)
func (s *Service) Processing(ctx context.Context, photo model.Photo, _ []byte) error {
	uploadData, err := s.storage.GetUploadPhotoData(ctx, photo.ID)
	if err != nil {
		return fmt.Errorf("storage.GetPhotoUploadData: %w", err)
	}

	if uploadData == nil {
		return fmt.Errorf("upload data not found")
	}

	metaData, err := s.storage.GetMetaData(ctx, photo.ID)
	if err != nil {
		return fmt.Errorf("storage.GetPhotoMetadata: %w", err)
	}

	if metaData == nil {
		return fmt.Errorf("meta not found")
	}

	// Теги По каталогу По каталогу КАТАЛОГ
	if err := s.createPhotoCatalogTag(ctx, photo, uploadData); err != nil {
		return fmt.Errorf("s.createPhotoCatalogTag: %w", err)
	}
	// Теги По дате формируем тег ГОД
	if err := s.createYearTag(ctx, photo, metaData); err != nil {
		return fmt.Errorf("s.createYearTag: %w", err)
	}
	// Теги По модели МОДЕЛЬ
	if err := s.createCameraModelTag(ctx, photo, metaData); err != nil {
		return fmt.Errorf("s.createCameraModelTag: %w", err)
	}
	// Теги по геолокации
	if err := s.createLocationTag(ctx, photo, metaData); err != nil {
		return fmt.Errorf("s.createLocationTag: %w", err)
	}

	return nil
}
