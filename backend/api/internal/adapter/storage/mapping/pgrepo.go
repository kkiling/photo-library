package mapping

import (
	"fmt"
	entity2 "github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
	"reflect"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

func PhotoEntityToModel(in *entity2.Photo) *model.Photo {
	if in == nil {
		return nil
	}
	return &model.Photo{
		ID:               in.ID,
		FileName:         in.FileName,
		Hash:             in.Hash,
		UpdateAt:         in.UpdateAt,
		Extension:        model.PhotoExtension(in.Extension),
		ProcessingStatus: model.PhotoProcessingStatus(in.ProcessingStatus),
	}
}

func PhotoModelToEntity(in *model.Photo) *entity2.Photo {
	if in == nil {
		return nil
	}
	return &entity2.Photo{
		ID:               in.ID,
		FileName:         in.FileName,
		Hash:             in.Hash,
		UpdateAt:         in.UpdateAt,
		Extension:        string(in.Extension),
		ProcessingStatus: string(in.ProcessingStatus),
	}
}

func PhotoUploadDataModelToEntity(in *model.PhotoUploadData) *entity2.PhotoUploadData {
	if in == nil {
		return nil
	}
	return &entity2.PhotoUploadData{
		PhotoID:  in.PhotoID,
		Paths:    in.Paths,
		UploadAt: in.UploadAt,
		ClientId: in.ClientId,
	}
}

func PhotoUploadDataEntityToModel(in *entity2.PhotoUploadData) *model.PhotoUploadData {
	if in == nil {
		return nil
	}
	return &model.PhotoUploadData{
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

func ExifPhotoDataEntityToModel(in *entity2.ExifPhotoData) *model.ExifPhotoData {
	if in == nil {
		return nil
	}
	return mapExifData(in, &model.ExifPhotoData{}).(*model.ExifPhotoData)
}

func ExifPhotoDataModelToEntity(in *model.ExifPhotoData) *entity2.ExifPhotoData {
	return mapExifData(in, &entity2.ExifPhotoData{}).(*entity2.ExifPhotoData)
}

func PhotoMetadataModelToEntity(in *model.PhotoMetadata) *entity2.PhotoMetadata {
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
	return &entity2.PhotoMetadata{
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

func PhotoMetadataEntityToModel(in *entity2.PhotoMetadata) *model.PhotoMetadata {
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

	return &model.PhotoMetadata{
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

func TagEntityToModel(in *entity2.Tag) *model.Tag {
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

func TagModelToEntity(in *model.Tag) *entity2.Tag {
	if in == nil {
		return nil
	}
	return &entity2.Tag{
		ID:         in.ID,
		CategoryID: in.CategoryID,
		PhotoID:    in.PhotoID,
		Name:       in.Name,
	}
}

func TagCategoryEntityToModel(in *entity2.TagCategory) *model.TagCategory {
	if in == nil {
		return nil
	}
	return &model.TagCategory{
		ID:    in.ID,
		Type:  in.Type,
		Color: in.Color,
	}
}

func TagCategoryModelToEntity(in *model.TagCategory) *entity2.TagCategory {
	if in == nil {
		return nil
	}
	return &entity2.TagCategory{
		ID:    in.ID,
		Type:  in.Type,
		Color: in.Color,
	}
}

func PhotoVectorEntityToModel(in *entity2.PhotoVector) *model.PhotoVector {
	if in == nil {
		return nil
	}
	return &model.PhotoVector{
		PhotoID: in.PhotoID,
		Vector:  in.Vector,
		Norm:    in.Norm,
	}
}

func PhotoVectorModelToEntity(in *model.PhotoVector) *entity2.PhotoVector {
	if in == nil {
		return nil
	}
	return &entity2.PhotoVector{
		PhotoID: in.PhotoID,
		Vector:  in.Vector,
		Norm:    in.Norm,
	}
}

func CoeffSimilarPhotoEntityToModel(in *entity2.CoeffSimilarPhoto) *model.CoeffSimilarPhoto {
	if in == nil {
		return nil
	}
	return &model.CoeffSimilarPhoto{
		PhotoID1:    in.PhotoID1,
		PhotoID2:    in.PhotoID2,
		Coefficient: in.Coefficient,
	}
}

func CoeffSimilarPhotoModelToEntity(in *model.CoeffSimilarPhoto) *entity2.CoeffSimilarPhoto {
	if in == nil {
		return nil
	}
	return &entity2.CoeffSimilarPhoto{
		PhotoID1:    in.PhotoID1,
		PhotoID2:    in.PhotoID2,
		Coefficient: in.Coefficient,
	}
}

func PhotoFilter(in *model.PhotoFilter) *entity2.PhotoFilter {
	if in == nil {
		return nil
	}

	var filter *entity2.PhotoFilter = nil
	filter = &entity2.PhotoFilter{}
	filter.ProcessingStatusIn = make([]string, 0, len(in.ProcessingStatusIn))
	for _, s := range in.ProcessingStatusIn {
		filter.ProcessingStatusIn = append(filter.ProcessingStatusIn, string(s))
	}

	return filter
}

func PhotoSelectParams(in model.PhotoSelectParams) entity2.PhotoSelectParams {
	return entity2.PhotoSelectParams{
		Offset:     in.Offset,
		Limit:      in.Limit,
		SortOrder:  entity2.PhotoSortOrder(in.SortOrder),
		SortDirect: entity2.SortDirect(in.SortOrder),
	}
}
