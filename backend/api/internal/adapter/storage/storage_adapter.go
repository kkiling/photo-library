package storage

import (
	"context"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/mapping"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/pgrepo"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type Adapter struct {
	photoRepo *pgrepo.PhotoRepository
}

func NewStorageAdapter(photoRepo *pgrepo.PhotoRepository) *Adapter {
	return &Adapter{
		photoRepo: photoRepo,
	}
}

func (r *Adapter) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	return r.photoRepo.RunTransaction(ctx, txFunc)
}

func (r *Adapter) GetPhotoByHash(ctx context.Context, hash string) (*model.Photo, error) {
	res, err := r.photoRepo.GetPhotoByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return mapping.PhotoEntityToModel(res), nil
}

func (r *Adapter) SavePhoto(ctx context.Context, photo model.Photo) error {
	in := mapping.PhotoModelToEntity(&photo)
	return r.photoRepo.SavePhoto(ctx, *in)
}

func (r *Adapter) MakeNotValidPhoto(ctx context.Context, photoID uuid.UUID, error string) error {
	return r.photoRepo.MakeNotValidPhoto(ctx, photoID, error)
}

func (r *Adapter) GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error) {
	res, err := r.photoRepo.GetPhotoById(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapping.PhotoEntityToModel(res), nil
}

func (r *Adapter) GetPhotosCount(ctx context.Context, filter *model.PhotoFilter) (int64, error) {
	return r.photoRepo.GetPhotosCount(ctx, mapping.PhotoFilter(filter))
}

func (r *Adapter) AddPhotosProcessingStatus(ctx context.Context, photoID uuid.UUID, status model.PhotoProcessingStatus, success bool) error {
	return r.photoRepo.AddPhotosProcessingStatus(ctx, photoID, string(status), success)
}

func (r *Adapter) GetUnprocessedPhotoIDs(ctx context.Context, lastProcessingStatus model.PhotoProcessingStatus, limit int64) ([]uuid.UUID, error) {
	return r.photoRepo.GetUnprocessedPhotoIDs(ctx, string(lastProcessingStatus), limit)
}

func (r *Adapter) GetPhotoProcessingStatuses(ctx context.Context, photoID uuid.UUID) ([]model.PhotoProcessingStatus, error) {
	statuses, err := r.photoRepo.GetPhotoProcessingStatuses(ctx, photoID)
	if err != nil {
		return nil, err
	}

	var result = make([]model.PhotoProcessingStatus, 0, len(statuses))
	for _, s := range statuses {
		result = append(result, model.PhotoProcessingStatus(s))
	}
	return result, nil
}

func (r *Adapter) GetPaginatedPhotos(ctx context.Context, params model.PhotoSelectParams, filter *model.PhotoFilter) ([]model.Photo, error) {
	photos, err := r.photoRepo.GetPaginatedPhotos(ctx, mapping.PhotoSelectParams(params), mapping.PhotoFilter(filter))
	if err != nil {
		return nil, err
	}
	result := make([]model.Photo, 0, len(photos))
	for _, p := range photos {
		result = append(result, *mapping.PhotoEntityToModel(&p))
	}
	return result, nil
}

func (r *Adapter) SaveUploadPhotoData(ctx context.Context, data model.PhotoUploadData) error {
	in := mapping.PhotoUploadDataModelToEntity(&data)
	return r.photoRepo.SavePhotoUploadData(ctx, *in)
}

func (r *Adapter) GetUploadPhotoData(ctx context.Context, photoID uuid.UUID) (*model.PhotoUploadData, error) {
	res, err := r.photoRepo.GetPhotoUploadData(ctx, photoID)
	if err != nil {
		return nil, err
	}
	return mapping.PhotoUploadDataEntityToModel(res), nil
}

func (r *Adapter) SaveOrUpdateExif(ctx context.Context, data *model.ExifPhotoData) error {
	in := mapping.ExifPhotoDataModelToEntity(data)
	return r.photoRepo.SaveOrUpdateExif(ctx, in)
}

func (r *Adapter) GetExif(ctx context.Context, photoID uuid.UUID) (*model.ExifPhotoData, error) {
	res, err := r.photoRepo.GetExif(ctx, photoID)
	if err != nil {
		return nil, err
	}
	return mapping.ExifPhotoDataEntityToModel(res), nil
}

func (r *Adapter) SaveOrUpdateMetaData(ctx context.Context, data model.PhotoMetadata) error {
	in := mapping.PhotoMetadataModelToEntity(&data)
	return r.photoRepo.SaveOrUpdatePhotoMetadata(ctx, in)
}

func (r *Adapter) GetMetaData(ctx context.Context, photoID uuid.UUID) (*model.PhotoMetadata, error) {
	res, err := r.photoRepo.GetPhotoMetadata(ctx, photoID)
	if err != nil {
		return nil, err
	}
	return mapping.PhotoMetadataEntityToModel(res), nil
}

func (r *Adapter) GetTagCategory(ctx context.Context, categoryID uuid.UUID) (*model.TagCategory, error) {
	res, err := r.photoRepo.GetTagCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	return mapping.TagCategoryEntityToModel(res), nil
}

func (r *Adapter) GetTagCategoryByType(ctx context.Context, typeCategory string) (*model.TagCategory, error) {
	res, err := r.photoRepo.GetTagCategoryByType(ctx, typeCategory)
	if err != nil {
		return nil, err
	}
	return mapping.TagCategoryEntityToModel(res), nil
}

func (r *Adapter) SaveTagCategory(ctx context.Context, category model.TagCategory) error {
	in := mapping.TagCategoryModelToEntity(&category)
	return r.photoRepo.SaveTagCategory(ctx, *in)
}

func (r *Adapter) GetTagByName(ctx context.Context, photoID uuid.UUID, name string) (*model.Tag, error) {
	res, err := r.photoRepo.GetTagByName(ctx, photoID, name)
	if err != nil {
		return nil, err
	}
	return mapping.TagEntityToModel(res), nil
}

func (r *Adapter) SaveTag(ctx context.Context, tag model.Tag) error {
	in := mapping.TagModelToEntity(&tag)
	return r.photoRepo.SaveTag(ctx, *in)
}

func (r *Adapter) SaveOrUpdatePhotoVector(ctx context.Context, photoVector model.PhotoVector) error {
	in := mapping.PhotoVectorModelToEntity(&photoVector)
	return r.photoRepo.SaveOrUpdatePhotoVector(ctx, *in)
}

func (r *Adapter) GetPaginatedPhotosVector(ctx context.Context, offset int64, limit int) ([]model.PhotoVector, error) {
	vectors, err := r.photoRepo.GetPaginatedPhotoVectors(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := make([]model.PhotoVector, 0, len(vectors))
	for _, p := range vectors {
		result = append(result, *mapping.PhotoVectorEntityToModel(&p))
	}
	return result, nil
}

func (r *Adapter) SaveSimilarPhotoCoefficient(ctx context.Context, sim model.CoeffSimilarPhoto) error {
	in := mapping.CoeffSimilarPhotoModelToEntity(&sim)
	return r.photoRepo.SaveCoeffSimilarPhoto(ctx, *in)
}
