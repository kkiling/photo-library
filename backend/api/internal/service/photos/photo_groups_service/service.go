package photo_groups_service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/utils"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Storage interface {
	service.Transactor
	GetPhotoGroupsCount(ctx context.Context, _ model.PhotoGroupsFilter) (int64, error)
	GetPaginatedPhotoGroups(ctx context.Context, params model.PhotoGroupsParams) ([]model.PhotoGroupWithPhotoIDs, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (model.PhotoGroupWithPhotoIDs, error)
	GetPhotoById(ctx context.Context, id uuid.UUID) (model.Photo, error)
	GetPhotoByFilename(ctx context.Context, fileName string) (model.Photo, error)
	GetPhotoPreviews(ctx context.Context, photoID uuid.UUID) ([]model.PhotoPreview, error)
	SetPhotoGroupMainPhoto(ctx context.Context, groupID, photoID uuid.UUID) error
}

type FileStore interface {
	GetFileBody(ctx context.Context, fileKey string) ([]byte, error)
	GetFileUrl(ctx context.Context, fileKey string) string
}

type Service struct {
	logger      log.Logger
	storage     Storage
	fileStorage FileStore
}

func NewService(logger log.Logger, fileStorage FileStore, storage Storage) *Service {
	return &Service{
		logger:      logger,
		storage:     storage,
		fileStorage: fileStorage,
	}
}

func (s *Service) getPreviews(ctx context.Context, photoID uuid.UUID) ([]model.PhotoPreview, error) {
	// load preview
	// В начале самый маленький
	previews, err := s.storage.GetPhotoPreviews(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoPreviews")
	}

	if len(previews) == 0 {
		return nil, serviceerr.NotFoundf("not found previews")
	}

	return previews, nil
}

func validateGetPhotoGroups(req *GetPhotoGroupsRequest) error {
	if err := req.Paginator.Validate(); err != nil {
		return err
	}
	return nil
}

func (s *Service) getPhotoWithPreviews(ctx context.Context, photoID uuid.UUID) (PhotoWithPreviews, error) {
	photo, err := s.storage.GetPhotoById(ctx, photoID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return PhotoWithPreviews{}, serviceerr.NotFoundf("photo %s from group not found", photoID.String())
	default:
		return PhotoWithPreviews{}, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
	}

	// В начале самый маленький, превью всегда есть
	previews, err := s.getPreviews(ctx, photoID)
	if err != nil {
		return PhotoWithPreviews{}, serviceerr.MakeErr(err, "s.getPreview")
	}

	if len(previews) == 0 {
		return PhotoWithPreviews{}, serviceerr.NotFoundf("not found previews for photo %s", photoID.String())
	}

	original := PhotoPreview{}
	photoPreviews := make([]PhotoPreview, 0, len(previews))
	for _, preview := range previews {
		curr := PhotoPreview{
			Src:    s.fileStorage.GetFileUrl(ctx, preview.FileKey),
			Width:  preview.WidthPixel,
			Height: preview.HeightPixel,
			Size:   preview.SizePixel,
		}
		if preview.Original {
			original = curr
		}
		photoPreviews = append(photoPreviews, curr)
	}
	return PhotoWithPreviews{
		ID:       photo.ID,
		Original: original,
		Previews: photoPreviews,
	}, nil
}

// GetPhotoGroups получение списка групп фотографий
func (s *Service) GetPhotoGroups(ctx context.Context, req *GetPhotoGroupsRequest) (*PaginatedPhotoGroups, error) {
	if err := validateGetPhotoGroups(req); err != nil {
		return nil, serviceerr.InvalidInputErr(err, "validateGetPhotoGroups")
	}

	groups, err := s.storage.GetPaginatedPhotoGroups(ctx, model.PhotoGroupsParams{
		Paginator: req.Paginator,
		// TODO: параметры
		//SortOrder:  "",
		//SortDirect: "",
		//Filter:     model.PhotoGroupsFilter{},
	})
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPaginatedPhotoGroups")
	}

	items := make([]PhotoGroup, 0, req.Paginator.PerPage)
	for _, group := range groups {
		photoWithPreviews, err := s.getPhotoWithPreviews(ctx, group.MainPhotoID)
		if err != nil {
			return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
		}

		photoGroup := PhotoGroup{
			ID:          group.ID,
			MainPhoto:   photoWithPreviews,
			PhotosCount: len(group.PhotoIDs),
		}
		items = append(items, photoGroup)
	}

	totalCount, err := s.storage.GetPhotoGroupsCount(ctx, model.PhotoGroupsFilter{
		// TODO: параметры
	})
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoGroupsCount")
	}

	return &PaginatedPhotoGroups{
		Items:      items,
		TotalItems: int(totalCount),
	}, nil
}

func (s *Service) GetPhotoGroup(ctx context.Context, groupID uuid.UUID) (*PhotoGroup, error) {
	group, err := s.storage.GetGroupByID(ctx, groupID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return nil, serviceerr.NotFoundf("group not found")
	default:
		return nil, serviceerr.MakeErr(err, "s.storage.GetGroupByID")
	}

	mainPhotoPreview, err := s.getPhotoWithPreviews(ctx, group.MainPhotoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
	}

	photosPreviews := make([]PhotoWithPreviews, 0, len(group.PhotoIDs))
	for _, photoID := range group.PhotoIDs {
		photoPreview, err := s.getPhotoWithPreviews(ctx, photoID)
		if err != nil {
			return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
		}
		photosPreviews = append(photosPreviews, photoPreview)
	}

	res := PhotoGroup{
		ID:          group.MainPhotoID,
		MainPhoto:   mainPhotoPreview,
		PhotosCount: len(group.PhotoIDs),
		Photos:      photosPreviews,
	}

	return &res, nil
}

func (s *Service) SetMainPhotoGroup(ctx context.Context, groupID, photoID uuid.UUID) error {
	group, err := s.storage.GetGroupByID(ctx, groupID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.NotFoundf("group not found")
	default:
		return serviceerr.MakeErr(err, "s.storage.GetGroupByID")
	}

	if !lo.Contains(group.PhotoIDs, photoID) {
		return serviceerr.NotFoundf("photo does not belong to the group")
	}

	if group.MainPhotoID == photoID {
		return serviceerr.Conflictf("photo already set to main photo")
	}

	if err = s.storage.SetPhotoGroupMainPhoto(ctx, groupID, photoID); err != nil {
		return serviceerr.MakeErr(err, "s.storage.SetPhotoGroupMainPhoto")
	}

	return nil
}

func (s *Service) GetPhotoContent(ctx context.Context, fileKey string) (*PhotoContent, error) {
	photoBody, err := s.fileStorage.GetFileBody(ctx, fileKey)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.fileStorage.GetFileBody")
	}
	ext := utils.GetPhotoExtension(fileKey)
	if ext == nil {
		return nil, serviceerr.InvalidInputf("invalid file extension")
	}
	return &PhotoContent{
		PhotoBody: photoBody,
		Extension: *ext,
	}, nil
}
