package metatags

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/geo"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const (
	YearTag             = "year"
	YearTagColor        = "#14ff00"
	CameraModelTag      = "camera_model"
	CameraModelTagColor = "#ffed00"
	LocationTag         = "location"
	LocationTagColor    = "#00ecff"
)

type Storage interface {
	service.Transactor
	GetMetaData(ctx context.Context, photoID uuid.UUID) (*model.PhotoMetadata, error)
	SaveGeoAddress(ctx context.Context, location model.Location) error
	GetGeoAddress(ctx context.Context, photoID uuid.UUID) (*model.Location, error)
}

type TagPhoto interface {
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error)
	GetOrCreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error)
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

func (s *Service) createYearTag(ctx context.Context, photo model.Photo, metaData *model.PhotoMetadata) error {
	if metaData.DateTime == nil {
		return nil
	}
	yearCategory, err := s.tagService.GetOrCreateCategory(ctx, YearTag, YearTagColor)
	if err != nil {
		return serviceerr.MakeErr(err, "s.getOrCreateTagCategory")
	}

	name := fmt.Sprintf("%d", metaData.DateTime.Year())
	_, err = s.tagService.AddPhotoTag(ctx, photo.ID, yearCategory.ID, name)
	if err != nil {
		if errors.Is(err, serviceerr.ErrTagAlreadyExist) {
		} else {
			return serviceerr.MakeErr(err, "tagService.AddPhotoTag")
		}
	}

	return nil
}

func (s *Service) createCameraModelTag(ctx context.Context, photo model.Photo, metaData *model.PhotoMetadata) error {
	if metaData.ModelInfo == nil || (*metaData.ModelInfo == "") {
		return nil
	}
	cameraCategory, err := s.tagService.GetOrCreateCategory(ctx, CameraModelTag, CameraModelTagColor)
	if err != nil {
		return serviceerr.MakeErr(err, "getOrCreateTagCategory")
	}

	_, err = s.tagService.AddPhotoTag(ctx, photo.ID, cameraCategory.ID, *metaData.ModelInfo)
	if err != nil {
		if errors.Is(err, serviceerr.ErrTagAlreadyExist) {
		} else {
			return serviceerr.MakeErr(err, "tagService.AddPhotoTag")
		}
	}

	return nil
}

func (s *Service) reverseGeocode(ctx context.Context, photoID uuid.UUID, latitude, longitude float64) (*geo.Address, error) {
	findLocation, err := s.storage.GetGeoAddress(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetGeoAddress")
	}
	if findLocation != nil {
		return &findLocation.Geo, nil
	}

	location, err := s.geocoder.ReverseGeocode(ctx, latitude, longitude)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "geocoder.ReverseGeocode")
	}

	saveLocation := model.Location{
		PhotoID:   photoID,
		CreatedAt: time.Now(),
		Latitude:  latitude,
		Longitude: longitude,
		Geo:       *location,
	}
	if err := s.storage.SaveGeoAddress(ctx, saveLocation); err != nil {
		return nil, fmt.Errorf("s.storage.SaveGeoAddress: %w", err)
	}

	return location, nil
}

func (s *Service) createLocationTag(ctx context.Context, photo model.Photo, metaData *model.PhotoMetadata) error {
	if metaData.Geo == nil {
		return nil
	}

	locationCategory, err := s.tagService.GetOrCreateCategory(ctx, LocationTag, LocationTagColor)
	if err != nil {
		return serviceerr.MakeErr(err, "getOrCreateTagCategory")
	}

	location, err := s.reverseGeocode(ctx, photo.ID, metaData.Geo.Latitude, metaData.Geo.Longitude)
	if err != nil {
		return serviceerr.MakeErr(err, "geocoder.ReverseGeocode")
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
				return serviceerr.MakeErr(err, "tagService.AddPhotoTag")
			}
		}
	}

	return nil
}

func (s *Service) Init(_ context.Context) error {
	return nil
}

func (s *Service) NeedLoadPhotoBody() bool {
	return false
}

// Processing создание и сохранение автоматических тегов (по мета данным или по путям и тд)
func (s *Service) Processing(ctx context.Context, photo model.Photo, _ []byte) (bool, error) {
	metaData, err := s.storage.GetMetaData(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.GetPhotoMetadata")
	}

	if metaData == nil {
		return false, nil
	}

	// Теги По дате формируем тег ГОД
	if err := s.createYearTag(ctx, photo, metaData); err != nil {
		return false, serviceerr.MakeErr(err, "s.createYearTag")
	}
	// Теги По модели МОДЕЛЬ
	if err := s.createCameraModelTag(ctx, photo, metaData); err != nil {
		return false, serviceerr.MakeErr(err, "s.createCameraModelTag")
	}
	// Теги по геолокации
	if err := s.createLocationTag(ctx, photo, metaData); err != nil {
		return false, serviceerr.MakeErr(err, "s.createLocationTag")
	}

	return true, nil
}
