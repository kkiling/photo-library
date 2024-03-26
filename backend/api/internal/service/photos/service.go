package photos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/utils"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Storage interface {
	service.Transactor
	GetPhotoGroupsCount(ctx context.Context) (uint64, error)
	GetPaginatedPhotoGroups(ctx context.Context, paginator model.Pagination) ([]model.PhotoGroup, error)
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	GetPhotoByFilename(ctx context.Context, fileName string) (*model.Photo, error)
	GetGroupByID(ctx context.Context, id uuid.UUID) (*model.PhotoGroup, error)
	GetMetaData(ctx context.Context, photoID uuid.UUID) (*model.PhotoMetadata, error)
	GetPhotoPreviews(ctx context.Context, photoID uuid.UUID) ([]model.PhotoPreview, error)
	GetPhotoPreviewFileName(ctx context.Context, photoID uuid.UUID, photoSize *int) (string, error)
	GetExif(ctx context.Context, photoID uuid.UUID) (*model.ExifPhotoData, error)
}
type TagPhoto interface {
	GetTags(ctx context.Context, photoID uuid.UUID) ([]model.Tag, error)
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*model.TagCategory, error)
}

type FileStore interface {
	GetFileBody(ctx context.Context, fileName string) ([]byte, error)
}

type Config struct {
	PhotoServerUrl string `yaml:"photo_server_url"`
	// DefaultPreviewSize int    `yaml:"default_preview_size"`
}

type Service struct {
	logger      log.Logger
	storage     Storage
	fileStorage FileStore
	tagService  TagPhoto
	cfg         Config
}

func NewService(logger log.Logger, cfg Config, tagService TagPhoto, fileStorage FileStore, storage Storage) *Service {
	return &Service{
		logger:      logger,
		storage:     storage,
		tagService:  tagService,
		fileStorage: fileStorage,
		cfg:         cfg,
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
		return nil, serviceerr.NotFoundError("not found previews")
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
	if err != nil {
		return PhotoWithPreviews{}, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
	}

	if photo == nil {
		return PhotoWithPreviews{}, serviceerr.NotFoundError("photo %s from group not found", photoID.String())
	}

	// В начале самый маленький, превью всегда есть
	previews, err := s.getPreviews(ctx, photoID)
	if err != nil {
		return PhotoWithPreviews{}, serviceerr.MakeErr(err, "s.getPreview")
	}

	if len(previews) == 0 {
		return PhotoWithPreviews{}, serviceerr.NotFoundError("not found previews for photo %s", photoID.String())
	}

	photoPreviews := make([]PhotoPreview, 0, len(previews))
	for _, preview := range previews {
		photoPreviews = append(photoPreviews, PhotoPreview{
			Src:    fmt.Sprintf("%s/%s?size=%d", s.cfg.PhotoServerUrl, photo.FileName, preview.SizePixel),
			Width:  preview.WidthPixel,
			Height: preview.HeightPixel,
			Size:   preview.SizePixel,
		})
	}

	last := photoPreviews[len(photoPreviews)-1]
	return PhotoWithPreviews{
		ID:       photo.ID,
		Src:      fmt.Sprintf("%s/%s", s.cfg.PhotoServerUrl, photo.FileName),
		Width:    last.Width,
		Height:   last.Height,
		Previews: photoPreviews,
	}, nil
}

// GetPhotoGroups получение списка групп фотографий
func (s *Service) GetPhotoGroups(ctx context.Context, req *GetPhotoGroupsRequest) (*PaginatedPhotoGroups, error) {
	if err := validateGetPhotoGroups(req); err != nil {
		return nil, serviceerr.InvalidInputErr(err, "validateGetPhotoGroups")
	}

	groups, err := s.storage.GetPaginatedPhotoGroups(ctx, req.Paginator)
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

	totalCount, err := s.storage.GetPhotoGroupsCount(ctx)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoGroupsCount")
	}

	return &PaginatedPhotoGroups{
		Items:      items,
		TotalItems: int(totalCount),
	}, nil
}

func (s *Service) GetPhotoContent(ctx context.Context, fileName string, previewSize *int) (*PhotoContent, error) {
	photoID, err := uuid.Parse(utils.FileNameWithoutExtSliceNotation(fileName))
	if err != nil {
		return nil, serviceerr.InvalidInputError("invalid file name uuid")
	}

	previewVileName, err := s.storage.GetPhotoPreviewFileName(ctx, photoID, previewSize)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoByFilename")
	}

	photoBody, err := s.fileStorage.GetFileBody(ctx, previewVileName)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.fileStorage.GetFileBody")
	}

	ext := utils.GetPhotoExtension(previewVileName)
	if ext == nil {
		return nil, serviceerr.InvalidInputError("invalid file extension")
	}

	return &PhotoContent{
		PhotoBody: photoBody,
		Extension: *ext,
	}, nil
}

func (s *Service) GetPhotoGroup(ctx context.Context, groupID uuid.UUID) (*PhotoGroupData, error) {
	group, err := s.storage.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetGroupByID")
	}

	if group == nil {
		return nil, serviceerr.NotFoundError("group not found")
	}

	photoWithPreview, err := s.getPhotoWithPreviews(ctx, group.MainPhotoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
	}

	photosWithPreview := make([]PhotoWithPreviews, 0, len(group.PhotoIDs))
	for _, photoID := range group.PhotoIDs {
		photoPreview, err := s.getPhotoWithPreviews(ctx, photoID)
		if err != nil {
			return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
		}
		photosWithPreview = append(photosWithPreview, photoPreview)
	}

	metaData, err := s.storage.GetMetaData(ctx, group.MainPhotoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetMetaData")
	}

	tags, err := s.tagService.GetTags(ctx, group.MainPhotoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.tagService.GetTags")
	}

	tagsWithCategories := make([]Tag, 0, len(tags))
	for _, tag := range tags {
		category, err := s.tagService.GetCategoryByID(ctx, tag.CategoryID)
		if err != nil {
			return nil, serviceerr.MakeErr(err, "s.tagService.GetCategoryByID")
		}

		tagsWithCategories = append(tagsWithCategories, Tag{
			ID:    tag.ID,
			Name:  tag.Name,
			Type:  category.Type,
			Color: category.Color,
		})
	}

	res := PhotoGroupData{
		PhotoGroup: PhotoGroup{
			ID:          group.MainPhotoID,
			MainPhoto:   photoWithPreview,
			PhotosCount: len(group.PhotoIDs),
		},
		Photos:   photosWithPreview,
		Metadata: metaData,
		Tags:     tagsWithCategories,
	}

	return &res, nil
}
