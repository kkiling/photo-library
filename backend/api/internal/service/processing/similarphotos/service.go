package similarphotos

import (
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"sync"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"gonum.org/v1/gonum/floats"
)

type Config struct {
	MinSimilarCoefficient float64 `yaml:"min_similar_coefficient"`
	Limit                 int64   `yaml:"limit"`
}

type Storage interface {
	service.Transactor
	GetPhotosVectorCount(ctx context.Context) (int64, error)
	GetPaginatedPhotosVector(ctx context.Context, offset int64, limit int64) ([]model.PhotoVector, error)
	GetPhotoVector(ctx context.Context, photoID uuid.UUID) (*model.PhotoVector, error)
	SaveSimilarPhotoCoefficient(ctx context.Context, sim model.CoeffSimilarPhoto) error
}

type Service struct {
	logger       log.Logger
	cfg          Config
	storage      Storage
	mu           sync.Mutex
	photoVectors map[uuid.UUID]model.PhotoVector
}

func NewService(logger log.Logger, cfg Config, storage Storage) *Service {
	return &Service{
		logger:       logger,
		cfg:          cfg,
		storage:      storage,
		mu:           sync.Mutex{},
		photoVectors: make(map[uuid.UUID]model.PhotoVector),
	}
}

func similarity(photoVector1, photoVector2 *model.PhotoVector) float64 {
	dotProduct := floats.Dot(photoVector1.Vector, photoVector2.Vector)
	return dotProduct / (photoVector1.Norm * photoVector2.Norm)
}

func (s *Service) Init(ctx context.Context) error {
	count, err := s.storage.GetPhotosVectorCount(ctx)
	if err != nil {
		return serviceerr.MakeErr(err, "storage.GetPhotosVectorCount")
	}

	for offset := int64(0); offset <= count; offset += s.cfg.Limit {
		vectors, err := s.storage.GetPaginatedPhotosVector(ctx, offset, s.cfg.Limit)
		if err != nil {
			return serviceerr.MakeErr(err, "storage.GetPaginatedPhotoVectors")
		}
		for _, vector := range vectors {
			s.photoVectors[vector.PhotoID] = vector
		}
	}

	return nil
}

func (s *Service) NeedLoadPhotoBody() bool {
	return false
}

// Processing рассчитывает коэффициент похожих фотографий
func (s *Service) Processing(ctx context.Context, photo model.Photo, _ []byte) (bool, error) {
	// Расчет векторов должен быть не конкурентным, а один за другим
	s.mu.Lock()
	defer s.mu.Unlock()

	currentPhotoVector, ok := s.photoVectors[photo.ID]
	if !ok {
		pv, err := s.storage.GetPhotoVector(ctx, photo.ID)
		if err != nil {
			return false, serviceerr.MakeErr(err, "storage.GetPhotoVector")
		}
		if pv == nil {
			return false, nil
		}
		s.photoVectors[photo.ID] = *pv
		currentPhotoVector = *pv
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

		err := s.storage.SaveSimilarPhotoCoefficient(ctx, model.CoeffSimilarPhoto{
			PhotoID1:    id1,
			PhotoID2:    id2,
			Coefficient: coefficient,
		})

		if err != nil {
			return false, serviceerr.MakeErr(err, "storage.SaveSimilarPhotoCoefficient")
		}
	}

	return true, nil
}
