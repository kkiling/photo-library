package photopreview

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Config struct {
	Sizes []int `yaml:"sizes"`
}

type Storage interface {
	service.Transactor
	GetPhotoPreviews(ctx context.Context, photoID uuid.UUID) ([]model.PhotoPreview, error)
	CreatePhotoPreview(ctx context.Context, preview model.PhotoPreview) error
	GetExif(ctx context.Context, photoID uuid.UUID) (*model.ExifPhotoData, error)
}

type FileStore interface {
	SaveFileBody(ctx context.Context, fileName string, body []byte) error
	DeleteFile(ctx context.Context, fileName string) error
}

type Service struct {
	logger      log.Logger
	cfg         Config
	storage     Storage
	fileStorage FileStore
}

func NewService(logger log.Logger, cfg Config, storage Storage, fileStorage FileStore) *Service {
	return &Service{
		logger:      logger,
		cfg:         cfg,
		storage:     storage,
		fileStorage: fileStorage,
	}
}

func (s *Service) Init(ctx context.Context) error {
	return nil
}

func (s *Service) NeedLoadPhotoBody() bool {
	return true
}

// Processing генерирует preview версии фотографий
func (s *Service) Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error) {
	/*var orientation = 1
	exif, err := s.storage.GetExif(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.GetExif")
	}

	if exif != nil && exif.Orientation != nil {
		orientation = *exif.Orientation
	}*/

	// TODO: Временно, удалить
	//previews, err := s.storage.GetPhotoPreviews(ctx, photo.ID)
	//if err != nil {
	//	return false, err
	//}
	//if len(previews) > 0 {
	//	return true, nil
	//}

	var originalImage ImageCV
	if err := originalImage.Load(photoBody); err != nil {
		return false, fmt.Errorf("image.Decode: %w, (%w)", err, serviceerr.ErrPhotoIsNotValid)
	}
	defer originalImage.Close()

	// TODO: СОХРАНЯТЬ превью в одной транзакции!

	widthPixel, heightPixel := originalImage.Size()
	sizePixel := widthPixel
	if heightPixel > widthPixel {
		sizePixel = heightPixel
	}

	photoPreview := model.PhotoPreview{
		ID:          photo.ID,
		PhotoID:     photo.ID,
		FileName:    photo.FileName,
		WidthPixel:  widthPixel,
		HeightPixel: heightPixel,
		SizePixel:   sizePixel,
	}

	/*if orientation == 6 || orientation == 8 {
		photoPreview.WidthPixel, photoPreview.HeightPixel = photoPreview.HeightPixel, photoPreview.WidthPixel
	}*/

	if err := s.storage.CreatePhotoPreview(ctx, photoPreview); err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.CreatePhotoPreview")
	}

	for _, maxSize := range s.cfg.Sizes {
		if maxSize >= sizePixel {
			break
		}

		imgPreview, err := createImagePreview(originalImage, widthPixel, heightPixel, maxSize, photo.Extension)
		if err != nil {
			return false, serviceerr.MakeErr(err, "failed to decode image")
		}

		previewID := uuid.New()
		// Сохранить файл и получить url
		fileName := fmt.Sprintf("preview_%d_%s.%s", maxSize, photo.ID, strings.ToLower(string(photo.Extension)))
		if err := s.fileStorage.SaveFileBody(ctx, fileName, imgPreview.photoBody); err != nil {
			return false, serviceerr.MakeErr(err, "s.fileStorage.SaveFileBody")
		}

		preview := model.PhotoPreview{
			ID:          previewID,
			PhotoID:     photo.ID,
			FileName:    fileName,
			WidthPixel:  imgPreview.width,
			HeightPixel: imgPreview.height,
			SizePixel:   maxSize,
		}
		if err := s.storage.CreatePhotoPreview(ctx, preview); err != nil {
			return false, serviceerr.MakeErr(err, "s.storage.CreatePhotoPreview")
		}
	}

	return true, nil
}
