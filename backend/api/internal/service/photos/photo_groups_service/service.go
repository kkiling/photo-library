package photo_groups_service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/form"
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

func validateGetPhotoGroups(req *form.GetPhotoGroups) error {
	if err := req.Paginator.Validate(); err != nil {
		return err
	}
	return nil
}

func (s *Service) getPhotoWithPreviews(ctx context.Context, photoID uuid.UUID) (model.PhotoWithPreviewsDTO, error) {
	// В начале самый маленький, превью всегда есть
	previews, err := s.getPreviews(ctx, photoID)
	if err != nil {
		return model.PhotoWithPreviewsDTO{}, serviceerr.MakeErr(err, "s.getPreview")
	}

	if len(previews) == 0 {
		return model.PhotoWithPreviewsDTO{}, serviceerr.NotFoundf("not found previews for photo %s", photoID.String())
	}

	original := model.PhotoPreviewDTO{}
	photoPreviews := make([]model.PhotoPreviewDTO, 0, len(previews))
	for _, preview := range previews {
		curr := model.PhotoPreviewDTO{
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
	return model.PhotoWithPreviewsDTO{
		PhotoID:  photoID,
		Original: original,
		Previews: photoPreviews,
	}, nil
}

// GetPhotoGroups получение списка групп фотографий
func (s *Service) GetPhotoGroups(ctx context.Context, req form.GetPhotoGroups) (model.PaginatedPhotoGroupsDTO, error) {
	if err := validateGetPhotoGroups(&req); err != nil {
		return model.PaginatedPhotoGroupsDTO{}, serviceerr.InvalidInputErr(err, "validateGetPhotoGroups")
	}

	groups, err := s.storage.GetPaginatedPhotoGroups(ctx, model.PhotoGroupsParams{
		Paginator: req.Paginator,
		// TODO: параметры
		//SortOrder:  "",
		//SortDirect: "",
		//Filter:     model.PhotoGroupsFilter{},
	})
	if err != nil {
		return model.PaginatedPhotoGroupsDTO{}, serviceerr.MakeErr(err, "s.storage.GetPaginatedPhotoGroups")
	}

	items := make([]model.PhotoGroupDTO, 0, req.Paginator.PerPage)
	for _, group := range groups {
		photoWithPreviews, err := s.getPhotoWithPreviews(ctx, group.MainPhotoID)
		if err != nil {
			return model.PaginatedPhotoGroupsDTO{}, serviceerr.MakeErr(err, "s.storage.getPhotoWithPreviews")
		}

		photoGroup := model.PhotoGroupDTO{
			GroupID:     group.ID,
			MainPhoto:   photoWithPreviews,
			PhotosCount: len(group.PhotoIDs),
		}
		items = append(items, photoGroup)
	}

	totalCount, err := s.storage.GetPhotoGroupsCount(ctx, model.PhotoGroupsFilter{
		// TODO: параметры
	})
	if err != nil {
		return model.PaginatedPhotoGroupsDTO{}, serviceerr.MakeErr(err, "s.storage.GetPhotoGroupsCount")
	}

	return model.PaginatedPhotoGroupsDTO{
		Items:      items,
		TotalItems: int(totalCount),
	}, nil
}

func (s *Service) GetPhotoGroup(ctx context.Context, groupID uuid.UUID) (model.PhotoGroupDTO, error) {
	group, err := s.storage.GetGroupByID(ctx, groupID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return model.PhotoGroupDTO{}, serviceerr.NotFoundf("group not found")
	default:
		return model.PhotoGroupDTO{}, serviceerr.MakeErr(err, "s.storage.GetGroupByID")
	}

	mainPhotoPreview, err := s.getPhotoWithPreviews(ctx, group.MainPhotoID)
	if err != nil {
		return model.PhotoGroupDTO{}, serviceerr.MakeErr(err, "s.storage.getPhotoWithPreviews")
	}

	photosPreviews := make([]model.PhotoWithPreviewsDTO, 0, len(group.PhotoIDs))
	for _, photoID := range group.PhotoIDs {
		photoPreview, err := s.getPhotoWithPreviews(ctx, photoID)
		if err != nil {
			return model.PhotoGroupDTO{}, serviceerr.MakeErr(err, "s.storage.getPhotoWithPreviews")
		}
		photosPreviews = append(photosPreviews, photoPreview)
	}

	res := model.PhotoGroupDTO{
		GroupID:     group.MainPhotoID,
		MainPhoto:   mainPhotoPreview,
		PhotosCount: len(group.PhotoIDs),
		Photos:      photosPreviews,
	}

	return res, nil
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

func (s *Service) GetPhotoContent(ctx context.Context, fileKey string) (model.PhotoContentDTO, error) {
	photoBody, err := s.fileStorage.GetFileBody(ctx, fileKey)
	if err != nil {
		return model.PhotoContentDTO{}, serviceerr.MakeErr(err, "s.fileStorage.GetFileBody")
	}
	ext := utils.GetPhotoExtension(fileKey)
	if ext == nil {
		return model.PhotoContentDTO{}, serviceerr.InvalidInputf("invalid file extension")
	}
	return model.PhotoContentDTO{
		PhotoBody: photoBody,
		Extension: *ext,
	}, nil
}
