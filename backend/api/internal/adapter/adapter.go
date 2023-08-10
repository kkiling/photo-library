package adapter

import (
	"context"
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

func (r *DbAdapter) SaveUploadPhotoData(ctx context.Context, data model.UploadPhotoData) error {
	in := mapping.UploadPhotoDataModelToEntity(&data)
	return r.photoRepo.SaveUploadPhotoData(ctx, *in)
}
