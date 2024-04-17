package photo_metadata_service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Storage interface {
	service.Transactor
	GetMetadata(ctx context.Context, photoID uuid.UUID) (model.PhotoMetadata, error)
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

func (s *Service) GetPhotoMetaData(ctx context.Context, photoID uuid.UUID) (*model.PhotoMetadata, error) {
	metaData, err := s.storage.GetMetadata(ctx, photoID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return nil, serviceerr.NotFoundf("metadata not found")
	default:
		return nil, serviceerr.MakeErr(err, "s.storage.GetMetaData")
	}
	return &metaData, nil
}
