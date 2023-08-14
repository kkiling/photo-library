package similarphotos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"gonum.org/v1/gonum/floats"
	"sync"
)

type Database interface {
	service.Transactor
	GetPhotosCount(ctx context.Context) (int64, error)
	ExistPhotoVector(ctx context.Context, photoID uuid.UUID) (bool, error)
	SaveOrUpdatePhotoVector(ctx context.Context, photoVector model.PhotoVector) error
	GetPaginatedPhotosVector(ctx context.Context, offset int64, limit int) ([]model.PhotoVector, error)
	SaveSimilarPhotoCoefficient(ctx context.Context, sim model.PhotosSimilarCoefficient) error
}

type TagPhoto interface {
	AddPhotoTag(ctx context.Context, photoID, categoryID uuid.UUID, name string) (model.Tag, error)
	GetCategory(ctx context.Context, typeCategory string) (*model.TagCategory, error)
	CreateCategory(ctx context.Context, typeCategory, color string) (model.TagCategory, error)
}

type PhotoML interface {
	GetImageVector(ctx context.Context, imgBytes []byte) ([]float64, error)
}

type Service struct {
	logger     log.Logger
	tagService TagPhoto
	database   Database
	photoML    PhotoML
}

func NewService(logger log.Logger, tagService TagPhoto, database Database, photoML PhotoML) *Service {
	return &Service{
		logger:     logger,
		tagService: tagService,
		database:   database,
		photoML:    photoML,
	}
}

func (s *Service) SavePhotoVector(ctx context.Context, photo model.Photo, photoBody []byte) error {
	if exist, err := s.database.ExistPhotoVector(ctx, photo.ID); err != nil {
		return fmt.Errorf("database.ExistPhotoVector: %e", err)
	} else if exist {
		return nil
	}

	vector, err := s.photoML.GetImageVector(ctx, photoBody)
	if err != nil {
		return fmt.Errorf("photoML.GetImageVector: %e", err)
	}

	norm := floats.Norm(vector, 2)
	if err := s.database.SaveOrUpdatePhotoVector(ctx, model.PhotoVector{
		PhotoID: photo.ID,
		Vector:  vector,
		Norm:    norm,
	}); err != nil {
		return fmt.Errorf("database.SaveOrUpdatePhotoVector: %e", err)
	}

	return nil
}

func similarity(photoVector1, photoVector2 *model.PhotoVector) float64 {
	dotProduct := floats.Dot(photoVector1.Vector, photoVector2.Vector)
	return dotProduct / (photoVector1.Norm * photoVector2.Norm)
}

func (s *Service) SavePhotoSimilarCoefficient(ctx context.Context, photo model.Photo, photoBody []byte) error {
	countPhotos, err := s.database.GetPhotosCount(ctx)
	if err != nil {
		return fmt.Errorf("database.GetPhotosCount: %w", err)
	}
	const limit = 1000
	const maxGoroutines = 20
	var offset int64
	var wg sync.WaitGroup
	photoVectors := make([]model.PhotoVector, 0, countPhotos)
	photoVectorsChan := make(chan model.PhotoVector)

	s.logger.Infof("Start read photos vectors")

	for offset = 0; offset < countPhotos; offset += limit {
		vectors, err := s.database.GetPaginatedPhotosVector(ctx, offset, limit)
		if err != nil {
			return fmt.Errorf("database.GetPaginatedPhotosVector: %w", err)
		}
		photoVectors = append(photoVectors, vectors...)
	}

	s.logger.Infof("Calculate photos Coefficient")

	go func() {
		defer close(photoVectorsChan)
		for _, vector := range photoVectors {
			photoVectorsChan <- vector
		}
	}()

	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for thisPhotoVector := range photoVectorsChan {
				for _, photoVector := range photoVectors {
					if thisPhotoVector.PhotoID == photoVector.PhotoID {
						continue
					}
					coefficient := similarity(&thisPhotoVector, &photoVector)
					if coefficient > 0.8 { // TODO: в конфиг
						err := s.database.SaveSimilarPhotoCoefficient(ctx, model.PhotosSimilarCoefficient{
							PhotoID1:    thisPhotoVector.PhotoID,
							PhotoID2:    photoVector.PhotoID,
							Coefficient: coefficient,
						})
						if err != nil {
							panic(err) // вынест в канал ошибок
						}
					}
				}
			}
		}()
	}
	wg.Wait()

	return nil
}
