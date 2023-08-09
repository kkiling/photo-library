package adapter

import (
	"context"
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
	//TODO implement me
	panic("implement me")
}

func (r *DbAdapter) SavePhoto(ctx context.Context, photo model.Photo) error {
	//TODO implement me
	panic("implement me")
}

func (r *DbAdapter) SaveUploadPhotoData(ctx context.Context, data model.UploadPhotoData) error {
	//TODO implement me
	panic("implement me")
}
