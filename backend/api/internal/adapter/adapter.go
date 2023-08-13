package adapter

import (
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/mapping"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type DbAdapter struct {
	photoRepo *pgrepo.PhotoRepository
}

func NewDbAdapter(photoRepo *pgrepo.PhotoRepository) *DbAdapter {
	return &DbAdapter{
		photoRepo: photoRepo,
	}
}

func (r *DbAdapter) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	return r.photoRepo.RunTransaction(ctx, txFunc)
}

func (r *DbAdapter) GetPhotoByHash(ctx context.Context, hash string) (*model.Photo, error) {
	res, err := r.photoRepo.GetPhotoByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return mapping.PhotoEntityToModel(res), nil
}

func (r *DbAdapter) SavePhoto(ctx context.Context, photo model.Photo) error {
	in := mapping.PhotoModelToEntity(&photo)
	return r.photoRepo.SavePhoto(ctx, *in)
}

func (r *DbAdapter) GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error) {
	res, err := r.photoRepo.GetPhotoById(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapping.PhotoEntityToModel(res), nil
}

func (r *DbAdapter) GetPhotosCount(ctx context.Context) (int64, error) {
	return r.photoRepo.GetPhotosCount(ctx)
}

func (r *DbAdapter) GetPaginatedPhotos(ctx context.Context, offset int64, limit int) ([]model.Photo, error) {
	photos, err := r.photoRepo.GetPaginatedPhotos(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := make([]model.Photo, 0, len(photos))
	for _, p := range photos {
		result = append(result, *mapping.PhotoEntityToModel(&p))
	}

	return result, nil
}

func (r *DbAdapter) SaveUploadPhotoData(ctx context.Context, data model.UploadPhotoData) error {
	in := mapping.UploadPhotoDataModelToEntity(&data)
	return r.photoRepo.SaveUploadPhotoData(ctx, *in)
}

func (r *DbAdapter) GetUploadPhotoData(ctx context.Context, photoID uuid.UUID) (*model.UploadPhotoData, error) {
	res, err := r.photoRepo.GetUploadPhotoData(ctx, photoID)
	if err != nil {
		return nil, err
	}
	return mapping.UploadPhotoDataEntityToModel(res), nil
}

func (r *DbAdapter) SaveOrUpdateExif(ctx context.Context, data *model.ExifData) error {
	in := mapping.ExifModelToExif(data)
	return r.photoRepo.SaveOrUpdateExif(ctx, in)
}

func (r *DbAdapter) GetExif(ctx context.Context, photoID uuid.UUID) (*model.ExifData, error) {
	res, err := r.photoRepo.GetExif(ctx, photoID)
	if err != nil {
		return nil, err
	}
	return mapping.ExifEntityToModel(res), nil
}

func (r *DbAdapter) SaveOrUpdateMetaData(ctx context.Context, data model.MetaData) error {
	in := mapping.MetaDataModelToEntity(&data)
	return r.photoRepo.SaveOrUpdateMetaData(ctx, in)
}

func (r *DbAdapter) GetMetaData(ctx context.Context, photoID uuid.UUID) (*model.MetaData, error) {
	res, err := r.photoRepo.GetMetaData(ctx, photoID)
	if err != nil {
		return nil, err
	}
	return mapping.MetaDataDataEntityToModel(res), nil
}
