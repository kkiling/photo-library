package systags

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/tagphoto"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
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
	geoPoints  []model.Geo
}

func NewService(logger log.Logger, tagService TagPhoto, database Database) *Service {
	return &Service{
		logger:     logger,
		tagService: tagService,
		database:   database,
		geoPoints:  make([]model.Geo, 0),
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
			return fmt.Errorf("tagService.AddPhotoTag: %w", err)
		}
	}

	// По модели МОДЕЛЬ
	if meta.ModelInfo != nil {
		yearCategory, err := s.getOrCreateTagCategory(ctx, CameraModelTag, CameraModelTagColor)
		if err != nil {
			return fmt.Errorf("getOrCreateTagCategory: %w", err)
		}

		_, err = s.tagService.AddPhotoTag(ctx, photo.ID, yearCategory.ID, *meta.ModelInfo)
		if err != nil {
			return fmt.Errorf("tagService.AddPhotoTag: %w", err)
		}
	}

	if meta.Geo != nil {
		s.geoPoints = append(s.geoPoints, *meta.Geo)
	}

	return nil
}

// TODO: придумать как сделать что бы теги можно было расчитывать по одной фотке, а не пачкой как щас

// CalculateGeoTags вычисление тегов по локации
func (s *Service) CalculateGeoTags(ctx context.Context) error {
	radius := 25.0 // км
	grouped := groupByMean(s.geoPoints, radius)
	for i, group := range grouped {
		fmt.Printf("Group %d: %v\n", i+1, group)
		mean := meanGeo(group)
		fmt.Printf("Mean for Group %d: Latitude: %f, Longitude: %f\n", i+1, mean.Latitude, mean.Longitude)
	}

	return nil
}
