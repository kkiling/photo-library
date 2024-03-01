package exifphotodata

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"strings"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type Storage interface {
	service.Transactor
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	SaveOrUpdateExif(ctx context.Context, data *model.ExifPhotoData) error
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

type write struct {
	data model.ExifPhotoData
}

func ratToFloat(tag *tiff.Tag, i int) (float64, error) {
	r1, r2, err := tag.Rat2(i)
	if err != nil {
		return 0, err
	}
	if r2 == 0 {
		return 0, nil
	}
	return float64(r1) / float64(r2), nil
}

func getIntFromTag(tag *tiff.Tag) (int, error) {
	if tag.Format() == tiff.IntVal {
		return tag.Int(0)
	}
	return 0, fmt.Errorf("unexpected tag format")
}

func getFloatFromTag(tag *tiff.Tag) (float64, error) {
	switch tag.Format() {
	case tiff.RatVal:
		return ratToFloat(tag, 0)
	case tiff.FloatVal:
		return tag.Float(0)
	case tiff.IntVal:
		val, err := tag.Int(0)
		if err != nil {
			return 0, err
		}
		return float64(val), nil
	default:
		return 0, fmt.Errorf("unexpected tag format")
	}
}

func getIntArrayFromTag(tag *tiff.Tag) ([]int, error) {
	if tag.Format() != tiff.IntVal {
		return nil, fmt.Errorf("unexpected tag format")
	}
	res := make([]int, tag.Count)
	for i := 0; i < int(tag.Count); i++ {
		val, err := tag.Int(i)
		if err != nil {
			return nil, err
		}
		res[i] = val
	}
	return res, nil
}

func getFloatArrayFromTag(tag *tiff.Tag) ([]float64, error) {
	count := int(tag.Count)
	res := make([]float64, count)

	var err error
	for i := 0; i < count; i++ {
		switch tag.Format() {
		case tiff.RatVal:
			res[i], err = ratToFloat(tag, i)
		case tiff.FloatVal:
			res[i], err = tag.Float(i)
		default:
			return nil, fmt.Errorf("unexpected tag format")
		}
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (p *write) Walk(name exif.FieldName, tag *tiff.Tag) error {
	fieldName := string(name)
	tp := determineDataType(&p.data, fieldName)

	switch tp {
	case dataTypeInt:
		if val, err := getIntFromTag(tag); err == nil {
			return setField(&p.data, fieldName, val)
		}
	case dataTypeString:
		val := strings.TrimSpace(strings.Trim(tag.String(), `"`))
		return setField(&p.data, fieldName, val)
	case dataTypeFloat:
		val, err := getFloatFromTag(tag)
		if err == nil {
			return setField(&p.data, fieldName, val)
		}
	case dataTypeIntArray:
		val, err := getIntArrayFromTag(tag)
		if err == nil {
			return setField(&p.data, fieldName, val)
		}
	case dataTypeFloatArray:
		val, err := getFloatArrayFromTag(tag)
		if err == nil {
			return setField(&p.data, fieldName, val)
		}
	default:
		return fmt.Errorf("unknown type: %v", tag.Format())
	}
	return nil
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
		data: model.ExifPhotoData{
			PhotoID: photo.ID,
		},
	}
	if err := x.Walk(&p); err != nil {
		return false, serviceerr.MakeErr(err, "exif.Walk")
	}

	err = s.storage.SaveOrUpdateExif(ctx, &p.data)

	if err != nil {
		return false, serviceerr.MakeErr(err, "storage.SaveOrUpdateExif")
	}

	return true, nil
}
