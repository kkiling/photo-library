package exifphoto

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"strings"
)

type Database interface {
	service.Transactor
	GetPhotosCount(ctx context.Context) (int64, error)
	GetPaginatedPhotos(ctx context.Context, offset int64, limit int) ([]model.Photo, error)
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	DeleteExif(ctx context.Context, photoID uuid.UUID) error
	SaveExif(ctx context.Context, data *model.ExifData) error
}

type FileStore interface {
	GetFileBody(ctx context.Context, filePath string) ([]byte, error)
}

type Service struct {
	database    Database
	fileStorage FileStore
}

func NewService(storage Database, fileStorage FileStore) *Service {
	return &Service{
		database:    storage,
		fileStorage: fileStorage,
	}
}

type write struct {
	data model.ExifData
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

// SavePhotoExifData рассчитывает exif данные фотографии и сохраняет в базу
func (s *Service) SavePhotoExifData(ctx context.Context, photo model.Photo, photoBody []byte) error {
	reader := bytes.NewReader(photoBody)
	x, err := exif.Decode(reader)

	if err != nil {
		if err.Error() == "EOF" {
			return ExifEOFErr
		}
		if exif.IsCriticalError(err) {
			return ExifCriticalErr
		}
	}

	var p = write{
		data: model.ExifData{
			PhotoID: photo.ID,
		},
	}
	if err := x.Walk(&p); err != nil {
		return fmt.Errorf("exif.Walk: %w", err)
	}

	err = s.database.RunTransaction(ctx, func(ctxTx context.Context) error {
		err = s.database.DeleteExif(ctx, photo.ID)
		if err != nil {
			return fmt.Errorf("database.DeleteExif: %w", err)
		}
		err = s.database.SaveExif(ctx, &p.data)
		if err != nil {
			return fmt.Errorf("database.SaveExif: %w", err)
		}
		return nil
	})

	return nil
}
