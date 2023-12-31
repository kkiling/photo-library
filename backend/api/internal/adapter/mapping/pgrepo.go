package mapping

import (
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"reflect"
)

func PhotoEntityToModel(in *entity.Photo) *model.Photo {
	if in == nil {
		return nil
	}
	return &model.Photo{
		ID:        in.ID,
		FilePath:  in.FilePath,
		Hash:      in.Hash,
		UpdateAt:  in.UpdateAt,
		UploadAt:  in.UploadAt,
		Extension: model.PhotoExtension(in.Extension),
	}
}

func PhotoModelToEntity(in *model.Photo) *entity.Photo {
	if in == nil {
		return nil
	}
	return &entity.Photo{
		ID:        in.ID,
		FilePath:  in.FilePath,
		Hash:      in.Hash,
		UpdateAt:  in.UpdateAt,
		UploadAt:  in.UploadAt,
		Extension: string(in.Extension),
	}
}

func UploadPhotoDataModelToEntity(in *model.UploadPhotoData) *entity.UploadPhotoData {
	if in == nil {
		return nil
	}
	return &entity.UploadPhotoData{
		PhotoID:  in.PhotoID,
		Paths:    in.Paths,
		UploadAt: in.UploadAt,
		ClientId: in.ClientId,
	}
}

func UploadPhotoDataEntityToModel(in *entity.UploadPhotoData) *model.UploadPhotoData {
	if in == nil {
		return nil
	}
	return &model.UploadPhotoData{
		PhotoID:  in.PhotoID,
		Paths:    in.Paths,
		UploadAt: in.UploadAt,
		ClientId: in.ClientId,
	}
}

func mapExifData(in interface{}, outTemplate interface{}) interface{} {
	if in == nil {
		return nil
	}

	inVal := reflect.ValueOf(in).Elem()
	outVal := reflect.ValueOf(outTemplate).Elem()

	inType := inVal.Type()
	outType := outVal.Type()

	if inType.NumField() != outType.NumField() {
		panic("fields count mismatch between structures")
	}

	for i := 0; i < inVal.NumField(); i++ {
		inField := inType.Field(i)
		outField, ok := outType.FieldByName(inField.Name)

		if !ok {
			panic(fmt.Sprintf("field %s is missing in the destination structure", inField.Name))
		}

		if inField.Type != outField.Type {
			panic(fmt.Sprintf("field %s type mismatch between structures", inField.Name))
		}

		outVal.FieldByName(inField.Name).Set(inVal.Field(i))
	}

	return outTemplate
}

func ExifEntityToModel(in *entity.ExifData) *model.ExifData {
	if in == nil {
		return nil
	}
	return mapExifData(in, &model.ExifData{}).(*model.ExifData)
}

func ExifModelToExif(in *model.ExifData) *entity.ExifData {
	return mapExifData(in, &entity.ExifData{}).(*entity.ExifData)
}

func MetaDataModelToEntity(in *model.MetaData) *entity.MetaData {
	if in == nil {
		return nil
	}
	var (
		geoLatitude  *float64
		geoLongitude *float64
	)
	if in.Geo != nil {
		geoLatitude = &in.Geo.Latitude
		geoLongitude = &in.Geo.Longitude
	}
	return &entity.MetaData{
		PhotoID:      in.PhotoID,
		ModelInfo:    in.ModelInfo,
		SizeBytes:    in.SizeBytes,
		WidthPixel:   in.WidthPixel,
		HeightPixel:  in.HeightPixel,
		DateTime:     in.DateTime,
		UpdateAt:     in.UpdateAt,
		GeoLatitude:  geoLatitude,
		GeoLongitude: geoLongitude,
	}
}

func MetaDataDataEntityToModel(in *entity.MetaData) *model.MetaData {
	if in == nil {
		return nil
	}
	var geo *model.Geo
	if in.GeoLatitude != nil && in.GeoLongitude != nil {
		geo = &model.Geo{
			Latitude:  *in.GeoLatitude,
			Longitude: *in.GeoLongitude,
		}
	}

	return &model.MetaData{
		PhotoID:     in.PhotoID,
		ModelInfo:   in.ModelInfo,
		SizeBytes:   in.SizeBytes,
		WidthPixel:  in.WidthPixel,
		HeightPixel: in.HeightPixel,
		DateTime:    in.DateTime,
		UpdateAt:    in.UpdateAt,
		Geo:         geo,
	}
}

func TagEntityToModel(in *entity.Tag) *model.Tag {
	if in == nil {
		return nil
	}
	return &model.Tag{
		ID:         in.ID,
		CategoryID: in.CategoryID,
		PhotoID:    in.PhotoID,
		Name:       in.Name,
	}
}

func TagModelToEntity(in *model.Tag) *entity.Tag {
	if in == nil {
		return nil
	}
	return &entity.Tag{
		ID:         in.ID,
		CategoryID: in.CategoryID,
		PhotoID:    in.PhotoID,
		Name:       in.Name,
	}
}

func TagCategoryEntityToModel(in *entity.TagCategory) *model.TagCategory {
	if in == nil {
		return nil
	}
	return &model.TagCategory{
		ID:    in.ID,
		Type:  in.Type,
		Color: in.Color,
	}
}

func TagCategoryModelToEntity(in *model.TagCategory) *entity.TagCategory {
	if in == nil {
		return nil
	}
	return &entity.TagCategory{
		ID:    in.ID,
		Type:  in.Type,
		Color: in.Color,
	}
}

func PhotoVectorEntityToModel(in *entity.PhotoVector) *model.PhotoVector {
	if in == nil {
		return nil
	}
	return &model.PhotoVector{
		PhotoID: in.PhotoID,
		Vector:  in.Vector,
		Norm:    in.Norm,
	}
}

func PhotoVectorModelToEntity(in *model.PhotoVector) *entity.PhotoVector {
	if in == nil {
		return nil
	}
	return &entity.PhotoVector{
		PhotoID: in.PhotoID,
		Vector:  in.Vector,
		Norm:    in.Norm,
	}
}

func PhotosSimilarCoefficientEntityToModel(in *entity.PhotosSimilarCoefficient) *model.PhotosSimilarCoefficient {
	if in == nil {
		return nil
	}
	return &model.PhotosSimilarCoefficient{
		PhotoID1:    in.PhotoID1,
		PhotoID2:    in.PhotoID2,
		Coefficient: in.Coefficient,
	}
}

func PhotosSimilarCoefficientModelToEntity(in *model.PhotosSimilarCoefficient) *entity.PhotosSimilarCoefficient {
	if in == nil {
		return nil
	}
	return &entity.PhotosSimilarCoefficient{
		PhotoID1:    in.PhotoID1,
		PhotoID2:    in.PhotoID2,
		Coefficient: in.Coefficient,
	}
}
