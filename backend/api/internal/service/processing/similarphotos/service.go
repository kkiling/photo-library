package similarphotos

import (
	"context"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"gonum.org/v1/gonum/floats"
)

type Storage interface {
	service.Transactor
	GetPhotosCount(ctx context.Context, filter *model.PhotoFilter) (int64, error)
	GetPaginatedPhotosVector(ctx context.Context, offset int64, limit int) ([]model.PhotoVector, error)
	SaveSimilarPhotoCoefficient(ctx context.Context, sim model.CoeffSimilarPhoto) error
}

type Service struct {
	logger  log.Logger
	storage Storage
}

func NewService(logger log.Logger, storage Storage) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
	}
}

func similarity(photoVector1, photoVector2 *model.PhotoVector) float64 {
	dotProduct := floats.Dot(photoVector1.Vector, photoVector2.Vector)
	return dotProduct / (photoVector1.Norm * photoVector2.Norm)
}

func (s *Service) SavePhotoSimilarCoefficient(ctx context.Context) error {
	const limit = 1000
	const maxGoroutines = 20
	const minSimilarCoefficient = 0.8

	countPhotos, err := s.storage.GetPhotosCount(ctx, nil)
	if err != nil {
		return serviceerr.MakeErr(err, "storage.GetPhotosCount")
	}

	var offset int64
	var wg sync.WaitGroup
	photoVectors := make([]model.PhotoVector, 0, countPhotos)
	photoVectorsChan := make(chan model.PhotoVector)
	errorsChan := make(chan error, maxGoroutines)

	for offset = 0; offset < countPhotos; offset += limit {
		vectors, err := s.storage.GetPaginatedPhotosVector(ctx, offset, limit)
		if err != nil {
			return serviceerr.MakeErr(err, "storage.GetPaginatedPhotoVectors")
		}
		photoVectors = append(photoVectors, vectors...)
	}

	go func() {
		defer close(photoVectorsChan)
		for _, vector := range photoVectors {
			photoVectorsChan <- vector
		}
	}()

	bar := pb.New(int(countPhotos)).Start()
	defer bar.Finish()

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
					if coefficient > minSimilarCoefficient {
						id1 := thisPhotoVector.PhotoID
						id2 := photoVector.PhotoID
						// Сравниваем UUID и ставим больший UUID в PhotoID1
						if id1.String() > id2.String() {
							id1, id2 = id2, id1
						}
						err := s.storage.SaveSimilarPhotoCoefficient(ctx, model.CoeffSimilarPhoto{
							PhotoID1:    id1,
							PhotoID2:    id2,
							Coefficient: coefficient,
						})
						if err != nil {
							errorsChan <- serviceerr.MakeErr(err, "storage.SaveCoeffSimilarPhoto")
							return
						}
					}
				}
				bar.Increment()
			}
		}()
	}
	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		return <-errorsChan
	}

	return nil
}
