package photometadata

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const timeLayout = "2006:01:02 15:04:05"
const invalidTime = "0000:00:00 00:00:00"

type Storage interface {
	service.Transactor
	SavePhotoMetadata(ctx context.Context, data model.PhotoMetadata) error
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
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

func (s *Service) getDateTime(exif *model.ExifPhotoData) (*time.Time, error) {

	if dateTimeStr, ok := exif.GetString(model.DateTimeExifData); ok {
		if d, err := parseDate(dateTimeStr); err == nil {
			return &d, nil
		}
	}

	if dateTimeStr, ok := exif.GetString(model.DateTimeOriginalExifData); ok {
		if d, err := parseDate(dateTimeStr); err == nil {
			return &d, nil
		}
	}

	return nil, nil
}

func (s *Service) getDateTimeFromPaths(ctx context.Context, photoID uuid.UUID) (*time.Time, error) {
	var dateTime *time.Time = nil

	// Попытка получения даты из имени файлов
	uploadData, err := s.storage.GetUploadPhotoData(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "storage.GetPhotoUploadData")
	}

	if uploadData == nil {
		return nil, serviceerr.NotFoundf("not found upload data for photo: %s", photoID)
	}

	for _, path := range uploadData.Paths {
		toTime, err := fileNameToTime(path)
		if err != nil {
			continue
		}
		dateTime = &toTime
		break
	}

	return dateTime, nil
}

func (s *Service) getGeo(exif *model.ExifPhotoData) (*model.Geo, error) {
	latitude, ok := exif.GetFloatArray(model.GPSLatitudeExifData)
	if !ok {
		return nil, nil
	}
	longitude, ok := exif.GetFloatArray(model.GPSLongitudeExifData)
	if !ok {
		return nil, nil
	}

	geo, err := convertToGeo(latitude, longitude)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "convertToGeo")
	}
	return geo, nil
}

func (s *Service) getModelInfo(exif *model.ExifPhotoData) *string {
	modelStr, okModel := exif.GetString(model.ModelExifData)
	makeStr, okMake := exif.GetString(model.MakeExifData)

	if !okModel && !okMake {
		return nil
	}

	if okModel && okMake {
		return lo.ToPtr(fmt.Sprintf("%s %s", modelStr, makeStr))
	}

	if okModel {
		return &modelStr
	}

	return &makeStr
}

func (s *Service) Init(_ context.Context) error {
	return nil
}

func (s *Service) NeedLoadPhotoBody() bool {
	return true
}

// Processing рассчитывает meta данные фотографии и сохраняет в базу
func (s *Service) Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error) {
	exif, err := s.storage.GetExif(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.GetExif")
	}

	var dateTime *time.Time = nil
	var geo *model.Geo = nil
	var modelInfo *string
	if exif != nil {
		dateTime, err = s.getDateTime(exif)
		if err != nil {
			return false, serviceerr.MakeErr(err, "getDateTime")
		}

		geo, err = s.getGeo(exif)
		if err != nil {
			return false, serviceerr.MakeErr(err, "getGeo")
		}

		modelInfo = s.getModelInfo(exif)
	}

	if dateTime == nil {
		dateTime, err = s.getDateTimeFromPaths(ctx, photo.ID)
		if err != nil {
			return false, serviceerr.MakeErr(err, "getDateTime")
		}
	}

	widthPixel, heightPixel, err := getImageDetails(photoBody)
	if err != nil {
		return false, serviceerr.MakeErr(err, "getImageDetails")
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

	err = s.storage.SavePhotoMetadata(ctx, meta)
	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.SaveOrUpdateMeta")
	}

	return true, nil
}
