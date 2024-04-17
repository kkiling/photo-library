package photo_preview

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

const previewsDir = "previews"

type Config struct {
	Sizes []int `yaml:"sizes"`
}

type Storage interface {
	service.Transactor
	GetExif(ctx context.Context, photoID uuid.UUID) (model.ExifPhotoData, error)
	GetPhotoPreviews(ctx context.Context, photoID uuid.UUID) ([]model.PhotoPreview, error)
	SavePhotoPreview(ctx context.Context, preview model.PhotoPreview) error
	DeletePhotoPreviews(ctx context.Context, photoID uuid.UUID) error
}

type FileStore interface {
	SaveFile(ctx context.Context, fileName string, body []byte, dirs ...string) (fileKey string, err error)
	DeleteFile(ctx context.Context, fileKey string) error
	DeleteFiles(ctx context.Context, fileKey []string) error
	DeleteDirectory(ctx context.Context, dirs ...string) error
}

type Processing struct {
	logger      log.Logger
	cfg         Config
	storage     Storage
	fileStorage FileStore
}

func NewService(logger log.Logger, cfg Config, storage Storage, fileStorage FileStore) *Processing {
	return &Processing{
		logger:      logger,
		cfg:         cfg,
		storage:     storage,
		fileStorage: fileStorage,
	}
}

func (s *Processing) Compensate(ctx context.Context, photoID uuid.UUID) error {
	err := s.storage.DeletePhotoPreviews(ctx, photoID)
	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.DeletePhotoPreviews")
	}

	// Удаляем файлы превью фотографий
	baseCatalog := fmt.Sprintf("preview_%s", photoID.String())
	err = s.fileStorage.DeleteDirectory(ctx, previewsDir, baseCatalog)
	if err != nil {
		return serviceerr.MakeErr(err, "s.fileStorage.DeleteFile")
	}

	return nil
}

func (s *Processing) Init(_ context.Context) error {
	return nil
}

func (s *Processing) NeedLoadPhotoBody() bool {
	return true
}

// Processing генерирует preview версии фотографий
func (s *Processing) Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error) {
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
	photoPreviews := make([]model.PhotoPreview, 0)
	files := make([]string, 0)

	widthPixel, heightPixel := originalImage.Size()
	sizePixel := widthPixel
	if heightPixel > widthPixel {
		sizePixel = heightPixel
	}

	mainPhotoPreview := model.PhotoPreview{
		ID:          photo.ID,
		PhotoID:     photo.ID,
		FileKey:     photo.FileKey,
		WidthPixel:  widthPixel,
		HeightPixel: heightPixel,
		SizePixel:   sizePixel,
		Original:    true,
	}
	photoPreviews = append(photoPreviews, mainPhotoPreview)

	/*if orientation == 6 || orientation == 8 {
		photoPreview.WidthPixel, photoPreview.HeightPixel = photoPreview.HeightPixel, photoPreview.WidthPixel
	}*/

	//if err := s.storage.SavePhotoPreview(ctx, mainPhotoPreview); err != nil {
	//	return false, serviceerr.MakeErr(err, "s.storage.CreatePhotoPreview")
	//}

	for _, size := range s.cfg.Sizes {
		if size >= sizePixel {
			break
		}

		preview, err := createImagePreview(originalImage, widthPixel, heightPixel, size, photo.Extension)
		if err != nil {
			return false, serviceerr.MakeErr(err, "failed to decode image")
		}

		// Сохранить файл и получить url
		fileName := fmt.Sprintf("%s.%s", photo.ID.String(), strings.ToLower(string(photo.Extension)))
		baseCatalog := fmt.Sprintf("preview_%s", photo.ID.String())
		sizeCatalog := fmt.Sprintf("%d", size)
		fileKey, err := s.fileStorage.SaveFile(ctx, fileName, preview.photoBody, previewsDir, baseCatalog, sizeCatalog)
		if err != nil {
			return false, serviceerr.MakeErr(err, "s.fileStorage.SaveFileBody")
		}
		files = append(files, fileKey)

		photoPreviews = append(photoPreviews, model.PhotoPreview{
			ID:          uuid.New(),
			PhotoID:     photo.ID,
			FileKey:     fileKey,
			WidthPixel:  preview.width,
			HeightPixel: preview.height,
			SizePixel:   size,
			Original:    false,
		})
	}

	err := s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		for _, preview := range photoPreviews {
			if err := s.storage.SavePhotoPreview(ctxTx, preview); err != nil {
				return serviceerr.MakeErr(err, "s.storage.CreatePhotoPreview")
			}
		}
		return nil
	})

	if err != nil {
		// Если произошла ошибка создания превью в базе, то удаляем каталог с превью фотографий
		baseCatalog := fmt.Sprintf("preview_%s", photo.ID)
		delErr := s.fileStorage.DeleteDirectory(ctx, previewsDir, baseCatalog)
		if delErr != nil {
			s.logger.Errorf("fail fileStorage.DeleteFiles: %v", delErr)
		}
		return false, serviceerr.MakeErr(err, "s.storage.RunTransaction")
	}

	return true, nil
}
