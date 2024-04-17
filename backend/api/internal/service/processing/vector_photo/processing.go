package vector_photo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/floats"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Storage interface {
	service.Transactor
	SavePhotoVector(ctx context.Context, photoVector model.PhotoVector) error
	DeletePhotoVector(ctx context.Context, photoID uuid.UUID) error
}

type PhotoML interface {
	GetImageVector(ctx context.Context, imgBytes []byte) ([]float64, error)
}

type Processing struct {
	logger  log.Logger
	storage Storage
	photoML PhotoML
}

func NewService(logger log.Logger, storage Storage, photoML PhotoML) *Processing {
	return &Processing{
		logger:  logger,
		storage: storage,
		photoML: photoML,
	}
}

func (s *Processing) Compensate(ctx context.Context, photoID uuid.UUID) error {
	err := s.storage.DeletePhotoVector(ctx, photoID)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, serviceerr.ErrNotFound):
		return nil
	default:
		return serviceerr.MakeErr(err, "s.storage.DeletePhotoVector")
	}
}

func (s *Processing) Init(_ context.Context) error {
	return nil
}

func (s *Processing) NeedLoadPhotoBody() bool {
	return true
}

// Processing рассчитывает и сохраняет вектора фотографий, для расчета похожих фото
func (s *Processing) Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error) {
	vector, err := s.photoML.GetImageVector(ctx, photoBody)
	if err != nil {
		if errors.Is(err, photoml.ErrInternalServerError) {
			return false, fmt.Errorf("s.photoML.GetImageVector: %w (%w)", err, serviceerr.ErrPhotoIsNotValid)
		}
		return false, serviceerr.MakeErr(err, "photoML.GetImageVector")
	}

	norm := floats.Norm(vector, 2)
	if err = s.storage.SavePhotoVector(ctx, model.PhotoVector{
		PhotoID: photo.ID,
		Vector:  vector,
		Norm:    norm,
	}); err != nil {
		return false, serviceerr.MakeErr(err, "storage.SaveOrUpdatePhotoVector")
	}

	return true, nil
}
