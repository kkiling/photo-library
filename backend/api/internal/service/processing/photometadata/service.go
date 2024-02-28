package photometadata

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const timeLayout = "2006:01:02 15:04:05"
const invalidTime = "0000:00:00 00:00:00"

type Storage interface {
	service.Transactor
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	SaveOrUpdateMetaData(ctx context.Context, data model.PhotoMetadata) error
	GetExif(ctx context.Context, photoID uuid.UUID) (*model.ExifPhotoData, error)
	GetUploadPhotoData(ctx context.Context, photoID uuid.UUID) (*model.PhotoUploadData, error)
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

func (s *Service) getDateTime(ctx context.Context, exif *model.ExifPhotoData, photoID uuid.UUID) (*time.Time, error) {
	if exif == nil {
		return nil, nil
	}

	var dateTime *time.Time = nil

	if exif.DateTime != nil {
		if d, err := parseDate(*exif.DateTime); err == nil {
			dateTime = &d
		}
	} else if exif.DateTimeOriginal != nil {
		if d, err := parseDate(*exif.DateTimeOriginal); err == nil {
			dateTime = &d
		}
	} else {
		uploadData, err := s.storage.GetUploadPhotoData(ctx, photoID)
		if err != nil {
			return nil, fmt.Errorf("storage.GetPhotoUploadData: %w", err)
		}

		if uploadData == nil {
			return nil, fmt.Errorf("not found upload data for photo: %s", photoID)
		}

		for _, path := range uploadData.Paths {
			toTime, err := fileNameToTime(path)
			if err != nil {
				continue
			}
			dateTime = &toTime
			break
		}
	}

	return dateTime, nil
}

func (s *Service) getGeo(exif *model.ExifPhotoData) (*model.Geo, error) {
	if exif != nil && exif.GPSLongitude != nil && exif.GPSLatitude != nil {
		geo, err := convertToGeo(exif.GPSLatitude, exif.GPSLongitude)
		if err != nil {
			return nil, fmt.Errorf("convertToGeo: %w", err)
		}
		return geo, nil
	}

	return nil, nil
}

func (s *Service) getModelInfo(exif *model.ExifPhotoData) *string {
	if exif != nil && (exif.Model != nil || exif.Make != nil) {
		modelInfo := ""
		if exif.Model != nil && exif.Make != nil {
			modelInfo = fmt.Sprintf("%s %s", *exif.Model, *exif.Make)
		} else if exif.Model != nil {
			modelInfo = *exif.Model
		} else {
			modelInfo = *exif.Make
		}
		return &modelInfo
	}

	return nil
}

// Processing рассчитывает meta данные фотографии и сохраняет в базу
func (s *Service) Processing(ctx context.Context, photo model.Photo, photoBody []byte) error {
	exif, err := s.storage.GetExif(ctx, photo.ID)
	if err != nil {
		return fmt.Errorf("storage.GetExif: %w", err)
	}

	dateTime, err := s.getDateTime(ctx, exif, photo.ID)
	if err != nil {
		return fmt.Errorf("getDateTime: %w", err)
	}
	geo, err := s.getGeo(exif)
	if err != nil {
		return fmt.Errorf("getGeo: %w", err)
	}

	modelInfo := s.getModelInfo(exif)

	widthPixel, heightPixel, err := getImageDetails(photoBody)
	if err != nil {
		return fmt.Errorf("getImageDetails: %w", err)
	}

	meta := model.PhotoMetadata{
		PhotoID:     photo.ID,
		ModelInfo:   modelInfo,
		SizeBytes:   len(photoBody),
		WidthPixel:  widthPixel,
		HeightPixel: heightPixel,
		DateTime:    dateTime,
		UpdateAt:    photo.UpdateAt,
		Geo:         geo,
	}

	err = s.storage.SaveOrUpdateMetaData(ctx, meta)
	if err != nil {
		// TODO:  ошибка
		return fmt.Errorf("storage.SaveOrUpdateMeta: %w", err)
	}

	return nil
}
