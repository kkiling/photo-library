package exifphoto

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type printer struct {
	dataTypes map[string]metaDataType
}

func (p *printer) Walk(name exif.FieldName, tag *tiff.Tag) error {
	switch tag.Format() {
	case tiff.IntVal:
		var t dataType = dataTypeInt
		if tag.Count > 1 {
			t = dataTypeIntArray
		}
		p.dataTypes[string(name)] = metaDataType{
			Name: string(name),
			Type: t,
		}
	case tiff.FloatVal:
		return fmt.Errorf("unknown type")
	case tiff.RatVal:
		var t dataType = dataTypeFloat
		if tag.Count > 1 {
			t = dataTypeFloatArray
		}
		p.dataTypes[string(name)] = metaDataType{
			Name: string(name),
			Type: t,
		}
	case tiff.StringVal, tiff.UndefVal:
		var t dataType = dataTypeString
		p.dataTypes[string(name)] = metaDataType{
			Name: string(name),
			Type: t,
		}
	case tiff.OtherVal:
		return fmt.Errorf("unknown type")
	}
	return nil
}

func (s *Service) getMetaDataType(ctx context.Context, photoId uuid.UUID) (map[string]metaDataType, error) {
	photo, err := s.database.GetPhotoById(ctx, photoId)
	if err != nil {
		return nil, fmt.Errorf("database.GetPhotoById: %w", err)
	}
	if photo == nil {
		return nil, fmt.Errorf("photo id=%s not found", photoId)
	}

	photoBody, err := s.fileStorage.GetFileBody(ctx, photo.FilePath)
	if err != nil {
		return nil, fmt.Errorf("fileStorage.GetFileBody: %w", err)
	}

	reader := bytes.NewReader(photoBody)
	x, err := exif.Decode(reader)

	if err != nil {
		if err.Error() == "EOF" {
			return nil, ExifEOFErr
		}
		if exif.IsCriticalError(err) {
			return nil, ExifCriticalErr
		}
	}

	var p printer = printer{
		dataTypes: make(map[string]metaDataType),
	}
	if err := x.Walk(&p); err != nil {
		return nil, fmt.Errorf("exif.Walk: %w", err)
	}

	return p.dataTypes, nil
}

// PrintExifData функция печатает возможные свойства exif и их типы, не используется в самом приложении
func (s *Service) PrintExifData(ctx context.Context) error {
	countPhotos, err := s.database.GetPhotosCount(ctx)
	if err != nil {
		return err
	}

	const limit = 1000
	var offset int64

	dataTypes := make(map[string]metaDataType)
	overrides := make(map[string]struct{})
	totalCount := 0
	failCount := 0
	for offset = 0; offset < countPhotos; offset += limit {
		photos, err := s.database.GetPaginatedPhotos(ctx, offset, limit)
		if err != nil {
			return err
		}
		for _, photo := range photos {
			totalCount++
			res, err := s.getMetaDataType(ctx, photo.ID)

			if err != nil {
				if errors.Is(err, ExifCriticalErr) || errors.Is(err, ExifEOFErr) {
					// fmt.Printf("Error photo: %s - %v\n", photo.FilePath, err)
					failCount++
				} else {
					return err
				}
			}

			for k, newVal := range res {
				oldVal, ok := dataTypes[k]
				if ok {
					if oldVal.Type != newVal.Type {
						_, ok2 := overrides[newVal.Name]
						if !ok2 {
							dataTypes[k] = metaDataType{
								Name: newVal.Name,
								Type: dataTypeString,
							}

							overrides[newVal.Name] = struct{}{}
							fmt.Println(newVal.Name + ": OVERRIDE TO STRING")
						}
					}
				} else {
					dataTypes[k] = newVal
				}
			}
		}
	}

	fmt.Printf("Total photos: %d\n", totalCount)
	fmt.Printf("Fail photos: %d\n", failCount)

	for _, newVal := range dataTypes {
		fmt.Printf("%s\t%s\n", newVal.Name, typeNames[newVal.Type])
	}

	return nil
}
