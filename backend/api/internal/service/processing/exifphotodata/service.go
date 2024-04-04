package exifphotodata

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/rwcarlsen/goexif/exif"
)

type Storage interface {
	service.Transactor
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	SaveExif(ctx context.Context, data *model.ExifPhotoData) error
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

func (s *Service) Init(ctx context.Context) error {
	return nil
}

func (s *Service) NeedLoadPhotoBody() bool {
	return true
}

// Processing рассчитывает exif данные фотографии и сохраняет в базу
func (s *Service) Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error) {
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

	if err := x.Walk(&p); err != nil {
		return false, serviceerr.MakeErr(err, "exif.Walk")
	}

	exifData := model.ExifPhotoData{
		PhotoID: photo.ID,
		Data:    p.data,
	}

	err = s.storage.SaveExif(ctx, &exifData)
	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.SaveOrUpdateExif")
	}

	return true, nil
}
