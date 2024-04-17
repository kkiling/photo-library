package exif_photo_data

import (
	"bytes"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Storage interface {
	service.Transactor
	GetPhotoById(ctx context.Context, id uuid.UUID) (model.Photo, error)
	SaveExif(ctx context.Context, data model.ExifPhotoData) error
	DeleteExif(ctx context.Context, photoID uuid.UUID) error
}

type Processing struct {
	logger  log.Logger
	storage Storage
}

func NewService(logger log.Logger, storage Storage) *Processing {
	return &Processing{
		logger:  logger,
		storage: storage,
	}
}

func (s *Processing) Compensate(ctx context.Context, photoID uuid.UUID) error {
	err := s.storage.DeleteExif(ctx, photoID)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, serviceerr.ErrNotFound):
		return nil
	default:
		return serviceerr.MakeErr(err, "s.storage.DeleteExif")
	}
}

func (s *Processing) Init(_ context.Context) error {
	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)
	return nil
}

func (s *Processing) NeedLoadPhotoBody() bool {
	return true
}

// Processing рассчитывает exif данные фотографии и сохраняет в базу
func (s *Processing) Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error) {
	reader := bytes.NewReader(photoBody)
	x, err := exif.Decode(reader)

	if err != nil {
		if err.Error() == "EOF" {
			s.logger.Debugf("photo (%s) exif decode error: %v", photo.ID, err)
			return false, nil
		}
		if exif.IsCriticalError(err) {
			s.logger.Debugf("photo (%s) exif decode error: %v", photo.ID, err)
			return false, nil
		}
	}

	var p = write{
		data: make(map[string]interface{}),
	}

	if err = x.Walk(&p); err != nil {
		return false, serviceerr.MakeErr(err, "exif.Walk")
	}

	exifData := model.ExifPhotoData{
		PhotoID: photo.ID,
		Data:    p.data,
	}

	err = s.storage.SaveExif(ctx, exifData)
	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.SaveOrUpdateExif")
	}

	return true, nil
}
