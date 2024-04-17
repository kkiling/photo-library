package similar_photos

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/floats"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/lock"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const (
	defaultTimeout = 6 * time.Second
)

type Config struct {
	MinSimilarCoefficient float64 `yaml:"min_similar_coefficient"`
	Limit                 uint64  `yaml:"limit"`
	Debug                 bool    `yaml:"debug"`
}

type RocketLockService interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (lock.RocketLockID, error)
	UnLock(ctx context.Context, lockID lock.RocketLockID) error
}

type Storage interface {
	service.Transactor
	GetPhotoVectors(ctx context.Context, pagination model.Pagination) ([]model.PhotoVector, error)
	GetPhotoVector(ctx context.Context, photoID uuid.UUID) (model.PhotoVector, error)
	DeleteCoefficientSimilarPhoto(ctx context.Context, photoID uuid.UUID) error
	SaveCoefficientSimilarPhotos(ctx context.Context, coefficient model.CoefficientSimilarPhoto) error
}

type Processing struct {
	logger       log.Logger
	cfg          Config
	storage      Storage
	lock         RocketLockService
	mu           sync.Mutex
	photoVectors map[uuid.UUID]model.PhotoVector
}

func NewService(logger log.Logger, cfg Config, storage Storage, lock RocketLockService) *Processing {
	return &Processing{
		logger:       logger,
		cfg:          cfg,
		storage:      storage,
		lock:         lock,
		mu:           sync.Mutex{},
		photoVectors: make(map[uuid.UUID]model.PhotoVector),
	}
}

func similarity(photoVector1, photoVector2 *model.PhotoVector) float64 {
	dotProduct := floats.Dot(photoVector1.Vector, photoVector2.Vector)
	return dotProduct / (photoVector1.Norm * photoVector2.Norm)
}

func (s *Processing) Compensate(ctx context.Context, photoID uuid.UUID) error {
	err := s.storage.DeleteCoefficientSimilarPhoto(ctx, photoID)
	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.DeleteCoefficientSimilarPhoto")
	}
	return nil
}

func (s *Processing) Init(ctx context.Context) error {
	var page uint64 = 0
	for {
		vectors, err := s.storage.GetPhotoVectors(ctx, model.Pagination{
			Page:    page,
			PerPage: s.cfg.Limit,
		})

		if err != nil {
			return serviceerr.MakeErr(err, "storage.GetPaginatedPhotoVectors")
		}

		if len(vectors) == 0 {
			break
		}

		for _, vector := range vectors {
			s.photoVectors[vector.PhotoID] = vector
		}
		page++
	}

	return nil
}

func (s *Processing) NeedLoadPhotoBody() bool {
	return false
}

// Processing рассчитывает коэффициент похожих фотографий
func (s *Processing) Processing(ctx context.Context, photo model.Photo, _ []byte) (bool, error) {
	// Расчет векторов должен быть не конкурентным, а один за другим
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.cfg.Debug {
		ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
		defer cancel()

		// Так же ставим лок (на случай если идет обработка из нескольких подов
		lockID, err := s.lock.Lock(ctx, "similar_coefficient", defaultTimeout)
		if err != nil {
			return false, err
		}
		defer func() {
			if unlockErr := s.lock.UnLock(ctx, lockID); unlockErr != nil {
				s.logger.Errorf("unlock: %v", unlockErr)
			}
		}()
	}

	currentPhotoVector, ok := s.photoVectors[photo.ID]
	if !ok {
		pv, err := s.storage.GetPhotoVector(ctx, photo.ID)
		switch {
		case err == nil:
			s.photoVectors[photo.ID] = pv
			currentPhotoVector = pv
		case errors.Is(err, serviceerr.ErrNotFound):
			return false, nil
		default:
			return false, serviceerr.MakeErr(err, "storage.GetPhotoVector")

		}
	}

	for _, photoVector := range s.photoVectors {
		if currentPhotoVector.PhotoID == photoVector.PhotoID {
			continue
		}

		coefficient := similarity(&currentPhotoVector, &photoVector)
		if coefficient < s.cfg.MinSimilarCoefficient {
			continue
		}

		id1 := currentPhotoVector.PhotoID
		id2 := photoVector.PhotoID

		// Сравниваем UUID и ставим больший UUID в PhotoID1
		if id1.String() > id2.String() {
			id1, id2 = id2, id1
		}

		err := s.storage.SaveCoefficientSimilarPhotos(ctx, model.CoefficientSimilarPhoto{
			PhotoID1:    id1,
			PhotoID2:    id2,
			Coefficient: coefficient,
		})

		if err != nil {
			return false, serviceerr.MakeErr(err, "storage.SaveCoefficientSimilarPhotos")
		}
	}

	return true, nil
}
