package vectorphoto

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"gonum.org/v1/gonum/floats"
)

type Database interface {
	service.Transactor
	ExistPhotoVector(ctx context.Context, photoID uuid.UUID) (bool, error)
	SaveOrUpdatePhotoVector(ctx context.Context, photoVector model.PhotoVector) error
}

type PhotoML interface {
	GetImageVector(ctx context.Context, imgBytes []byte) ([]float64, error)
}

type Service struct {
	logger   log.Logger
	database Database
	photoML  PhotoML
}

func NewService(logger log.Logger, database Database, photoML PhotoML) *Service {
	return &Service{
		logger:   logger,
		database: database,
		photoML:  photoML,
	}
}

// Processing рассчитывает и сохраняет вектора фотографий, для расчета похожих фото
func (s *Service) Processing(ctx context.Context, photo model.Photo, photoBody []byte) error {
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
