package vectorphoto

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"gonum.org/v1/gonum/floats"
)

type Storage interface {
	service.Transactor
	ExistPhotoVector(ctx context.Context, photoID uuid.UUID) (bool, error)
	SaveOrUpdatePhotoVector(ctx context.Context, photoVector model.PhotoVector) error
}

type PhotoML interface {
	GetImageVector(ctx context.Context, imgBytes []byte) ([]float64, error)
}

type Service struct {
	logger  log.Logger
	storage Storage
	photoML PhotoML
}

func NewService(logger log.Logger, storage Storage, photoML PhotoML) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
		photoML: photoML,
	}
}

// Processing рассчитывает и сохраняет вектора фотографий, для расчета похожих фото
func (s *Service) Processing(ctx context.Context, photo model.Photo, photoBody []byte) error {
	if exist, err := s.storage.ExistPhotoVector(ctx, photo.ID); err != nil {
		return fmt.Errorf("storage.ExistPhotoVector: %e", err)
	} else if exist {
		return nil
	}

	vector, err := s.photoML.GetImageVector(ctx, photoBody)
	if err != nil {
		if errors.Is(err, photoml.ErrInternalServerError) {
			return errors.Join(err, serviceerr.ErrPhotoIsNotValid)
		}
		return fmt.Errorf("photoML.GetImageVector: %w", err)
	}

	norm := floats.Norm(vector, 2)
	if err := s.storage.SaveOrUpdatePhotoVector(ctx, model.PhotoVector{
		PhotoID: photo.ID,
		Vector:  vector,
		Norm:    norm,
	}); err != nil {
		return fmt.Errorf("storage.SaveOrUpdatePhotoVector: %e", err)
	}

	return nil
}
