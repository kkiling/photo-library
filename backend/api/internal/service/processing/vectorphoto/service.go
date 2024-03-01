package vectorphoto

import (
	"context"
	"errors"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"go.uber.org/multierr"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"gonum.org/v1/gonum/floats"
)

type Storage interface {
	service.Transactor
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
func (s *Service) Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error) {
	vector, err := s.photoML.GetImageVector(ctx, photoBody)
	if err != nil {
		if errors.Is(err, photoml.ErrInternalServerError) {
			return false, multierr.Append(err, serviceerr.ErrPhotoIsNotValid)
		}
		return false, serviceerr.MakeErr(err, "photoML.GetImageVector")
	}

	norm := floats.Norm(vector, 2)
	if err := s.storage.SaveOrUpdatePhotoVector(ctx, model.PhotoVector{
		PhotoID: photo.ID,
		Vector:  vector,
		Norm:    norm,
	}); err != nil {
		return false, serviceerr.MakeErr(err, "storage.SaveOrUpdatePhotoVector")
	}

	return true, nil
}
