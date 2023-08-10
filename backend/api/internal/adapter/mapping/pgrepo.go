package mapping

import (
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
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
		ID:       in.ID,
		PhotoID:  in.PhotoID,
		Paths:    in.Paths,
		UploadAt: in.UploadAt,
		ClientId: in.ClientId,
	}
}
