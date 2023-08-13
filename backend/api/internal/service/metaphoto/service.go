package metaphoto

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"time"
)

const timeLayout = "2006:01:02 15:04:05"
const invalidTime = "0000:00:00 00:00:00"

var errInvalidDateFormat = fmt.Errorf("invalid date format")
var ErrExifNotFound = fmt.Errorf("exif not found")

type Database interface {
	service.Transactor
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	SaveOrUpdateMetaData(ctx context.Context, data model.MetaData) error
	GetExif(ctx context.Context, photoID uuid.UUID) (*model.ExifData, error)
	GetUploadPhotoData(ctx context.Context, photoID uuid.UUID) (*model.UploadPhotoData, error)
}

type Service struct {
	database Database
}

func NewService(storage Database) *Service {
	return &Service{
		database: storage,
	}
}

// SavePhotoMetaData рассчитывает meta данные фотографии и сохраняет в базу
func (s *Service) SavePhotoMetaData(ctx context.Context, photo model.Photo, photoBody []byte) error {
	exif, err := s.database.GetExif(ctx, photo.ID)
	if err != nil {
		return fmt.Errorf("database.GetExif: %w", err)
	}

	if exif == nil {
		return ErrExifNotFound
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
		uploadData, err := s.database.GetUploadPhotoData(ctx, photo.ID)
		if err != nil {
			return fmt.Errorf("database.GetUploadPhotoData: %w", err)
		}

		if exif == nil {
			return fmt.Errorf("not found upload data for photo: %s", photo.ID)
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

	var geo *model.Geo
	if exif.GPSLongitude != nil && exif.GPSLatitude != nil {
		geo, err = convertToGeo(exif.GPSLatitude, exif.GPSLongitude)
		if err != nil {
			return fmt.Errorf("convertToGeo: %w", err)
		}
	}

	var modelInfo *string
	if exif.Model != nil || exif.Make != nil {
		sss := ""
		if exif.Model != nil && exif.Make != nil {
			sss = fmt.Sprintf("%s %s", *exif.Model, *exif.Make)
		} else if exif.Model != nil {
			sss = *exif.Model
		} else {
			sss = *exif.Make
		}
		modelInfo = &sss
	}

	widthPixel, heightPixel, err := getImageDetails(photoBody)
	if err != nil {
		return fmt.Errorf("getImageDetails: %w", exif)
	}

	meta := model.MetaData{
		PhotoID:     photo.ID,
		ModelInfo:   modelInfo,
		SizeBytes:   len(photoBody),
		WidthPixel:  widthPixel,
		HeightPixel: heightPixel,
		DateTime:    dateTime,
		UpdateAt:    photo.UpdateAt,
		Geo:         geo,
	}

	err = s.database.SaveOrUpdateMetaData(ctx, meta)
	if err != nil {
		// TODO:  ошибка
		return fmt.Errorf("database.SaveOrUpdateMeta: %w", err)
	}

	return nil
}
