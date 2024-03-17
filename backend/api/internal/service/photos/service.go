package photos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
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

	// TODO: Костыль с поворотом
	var orientation = 1

	exif, err := s.storage.GetExif(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetExif")
	}

	if exif != nil && exif.Orientation != nil {
		orientation = *exif.Orientation
	}

	if orientation == 6 || orientation == 8 {
		for index, preview := range previews {
			preview.WidthPixel, preview.HeightPixel = preview.HeightPixel, preview.WidthPixel
			previews[index] = preview
		}
	}

	return previews, nil
}

func (s *Service) getPhoto(ctx context.Context, photoID uuid.UUID) (Photo, error) {
	photo, err := s.storage.GetPhotoById(ctx, photoID)
	if err != nil {
		return Photo{}, serviceerr.MakeErr(err, "s.storage.GetPhotoById")
	}

	if photo == nil {
		return Photo{}, serviceerr.NotFoundError("photo %s from group not found", photoID.String())
	}

	metaData, err := s.storage.GetMetaData(ctx, photoID)
	if err != nil {
		return Photo{}, serviceerr.MakeErr(err, "s.storage.GetMetaData")
	}

	tags, err := s.tagService.GetTags(ctx, photoID)
	if err != nil {
		return Photo{}, serviceerr.MakeErr(err, "s.tagService.GetTags")
	}

	tagsWithCategories := make([]Tag, 0, len(tags))
	for _, tag := range tags {
		category, err := s.tagService.GetCategoryByID(ctx, tag.CategoryID)
		if err != nil {
			return Photo{}, serviceerr.MakeErr(err, "s.tagService.GetCategoryByID")
		}

		tagsWithCategories = append(tagsWithCategories, Tag{
			ID:    tag.ID,
			Name:  tag.Name,
			Type:  category.Type,
			Color: category.Color,
		})
	}

	// В начале самый маленький, превью всегда есть
	previews, err := s.getPreviews(ctx, photoID)
	if err != nil {
		return Photo{}, serviceerr.MakeErr(err, "s.getPreview")
	}

	original := previews[len(previews)-1]
	lenPreviews := len(previews) - 1
	photoPreviews := make([]PhotoPreview, 0, lenPreviews)
	for _, preview := range previews[:lenPreviews] {
		photoPreviews = append(photoPreviews, PhotoPreview{
			Src:    fmt.Sprintf("%s/%s?size=%d", s.cfg.PhotoServerUrl, photo.FileName, preview.SizePixel),
			Width:  preview.WidthPixel,
			Height: preview.HeightPixel,
			Size:   preview.SizePixel,
		})
	}

	return Photo{
		ID:       photo.ID,
		Src:      fmt.Sprintf("%s/%s", s.cfg.PhotoServerUrl, photo.FileName),
		Width:    original.WidthPixel,
		Height:   original.HeightPixel,
		Size:     original.SizePixel,
		Metadata: metaData,
		Tags:     tagsWithCategories,
		Preview:  photoPreviews,
	}, nil
}

func validateGetPhotoGroups(req *GetPhotoGroupsRequest) error {
	if err := req.Paginator.Validate(); err != nil {
		return err
	}
	return nil
}

// GetPhotoGroups получение списка групп фотографий
func (s *Service) GetPhotoGroups(ctx context.Context, req *GetPhotoGroupsRequest) (*GetPhotoGroupsResponse, error) {
	if err := validateGetPhotoGroups(req); err != nil {
		return nil, serviceerr.InvalidInputErr(err, "validateGetPhotoGroups")
	}

	groups, err := s.storage.GetPaginatedPhotoGroups(ctx, req.Paginator)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPaginatedPhotoGroups")
	}

	items := make([]PhotoGroup, 0, req.Paginator.PerPage)
	for _, group := range groups {
		mainPhoto, err := s.getPhoto(ctx, group.MainPhotoID)
		if err != nil {
			return nil, serviceerr.MakeErr(err, "s.getPhoto")
		}
		photoGroup := PhotoGroup{
			ID:         group.ID,
			MainPhoto:  mainPhoto,
			PhotoCount: len(group.PhotoIDs),
		}
		items = append(items, photoGroup)
	}

	totalCount, err := s.storage.GetPhotoGroupsCount(ctx)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoGroupsCount")
	}

	return &GetPhotoGroupsResponse{
		Items:      items,
		TotalItems: int(totalCount),
	}, nil
}

func (s *Service) GetPhotoContent(ctx context.Context, fileName string, previewSize *int) (*PhotoContent, error) {
	photo, err := s.storage.GetPhotoByFilename(ctx, fileName)

	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoByFilename")
	}

	if photo == nil {
		return nil, serviceerr.NotFoundError("photo not found")
	}

	previews, err := s.getPreviews(ctx, photo.ID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.getPreview")
	}

	// В начале самый маленький
	currPreview := previews[len(previews)-1]
	if previewSize != nil {
		for _, preview := range previews {
			if preview.SizePixel >= *previewSize {
				currPreview = preview
				break
			}
		}
	}

	photoBody, err := s.fileStorage.GetFileBody(ctx, currPreview.FileName)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.fileStorage.GetFileBody")
	}

	return &PhotoContent{
		PhotoBody: photoBody,
		Extension: photo.Extension,
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

	res := PhotoGroupData{
		ID:        group.ID,
		MainPhoto: Photo{},
		Photos:    make([]Photo, 0, len(group.PhotoIDs)),
	}

	for _, photoID := range group.PhotoIDs {
		photo, err := s.getPhoto(ctx, photoID)
		if err != nil {
			return nil, serviceerr.MakeErr(err, "s.getPhoto")
		}

		res.Photos = append(res.Photos, photo)
		if group.MainPhotoID == photoID {
			res.MainPhoto = photo
		}
	}

	return &res, nil
}
