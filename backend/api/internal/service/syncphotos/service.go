package syncphotos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"path/filepath"
	"strings"
	"time"
)

type Storage interface {
	service.Transactor
	GetPhotoByHash(ctx context.Context, hash string) (*model.Photo, error)
	SavePhoto(ctx context.Context, photo model.Photo) error
	SaveUploadPhotoData(ctx context.Context, data model.PhotoUploadData) error
}

type FileStore interface {
	SaveFileBody(ctx context.Context, fileName string, body []byte) error
	DeleteFile(ctx context.Context, fileName string) error
}

type Service struct {
	logger      log.Logger
	storage     Storage
	fileStorage FileStore
}

func NewService(logger log.Logger, storage Storage, fileStorage FileStore) *Service {
	return &Service{
		logger:      logger,
		storage:     storage,
		fileStorage: fileStorage,
	}
}

func (s *Service) getPhotoExtension(path string) *model.PhotoExtension {
	// Извлекаем расширение файла из пути.
	ext := strings.ToUpper(filepath.Ext(path))
	// Удаляем точку из начала расширения.
	ext = strings.TrimPrefix(ext, ".")

	// Определяем, какому PhotoExtension соответствует извлеченное расширение.
	switch ext {
	case string(model.PhotoExtensionJpg), string(model.PhotoExtensionJpeg):
		photoExt := model.PhotoExtensionJpeg
		return &photoExt
	case string(model.PhotoExtensionPng):
		photoExt := model.PhotoExtensionPng
		return &photoExt
	case string(model.PhotoExtensionBmb):
		photoExt := model.PhotoExtensionBmb
		return &photoExt
	default:
		// Если расширение не соответствует известным типам, возвращаем nil.
		return nil
	}
}

func (s *Service) UploadPhoto(ctx context.Context, form *model.SyncPhotoRequest) (*model.SyncPhotoResponse, error) {
	// Проверяем загружено ли фото
	findPhoto, err := s.storage.GetPhotoByHash(ctx, form.Hash)
	if err != nil {
		return nil, serviceerr.RuntimeError(err, s.storage.GetPhotoByHash)
	}

	if findPhoto != nil {
		return &model.SyncPhotoResponse{
			HasBeenUploadedBefore: true,
			Hash:                  findPhoto.Hash,
		}, nil
	}

	if len(form.Paths) == 0 {
		return nil, serviceerr.InvalidInputError("paths must not be empty")
	}

	// Проверка, поддерживается ли расширение
	ex := s.getPhotoExtension(form.Paths[0])
	if ex == nil {
		return nil, serviceerr.InvalidInputError("photo extension not found")
	}

	photoID := uuid.New()
	uploadAt := time.Now()

	// Сохранить файл и получить url
	fileName := fmt.Sprintf("%s.%s", photoID, strings.ToLower(string(*ex)))
	if err := s.fileStorage.SaveFileBody(ctx, fileName, form.Body); err != nil {
		return nil, serviceerr.RuntimeError(err, s.fileStorage.SaveFileBody)
	}

	newPhoto := model.Photo{
		ID:               photoID,
		FileName:         fileName,
		Hash:             form.Hash,
		UpdateAt:         form.UpdateAt,
		Extension:        *ex,
		ProcessingStatus: model.NewPhoto,
	}

	uploadPhotoData := model.PhotoUploadData{
		PhotoID:  photoID,
		UploadAt: uploadAt,
		Paths:    form.Paths,
		ClientId: form.ClientId,
	}

	// Если произошла ошибка при сохранении фото в базе
	// То, пытаемся удалить его из файлового хранилища
	// Если программа упадет до удаления из хранилища, то файл зависнет и будет бесхозным
	// Нужна чистилка бесхозных файлов

	// Одной транзакцией сохранить
	err = s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		if saveErr := s.storage.SavePhoto(ctxTx, newPhoto); saveErr != nil {
			return saveErr
		}
		if saveErr := s.storage.SaveUploadPhotoData(ctxTx, uploadPhotoData); saveErr != nil {
			return saveErr
		}
		return nil
	})

	if err != nil {
		if delErr := s.fileStorage.DeleteFile(ctx, newPhoto.FileName); delErr != nil {
			s.logger.Errorf("fail fileStorage.DeleteFile %s: %v", newPhoto.FileName, delErr)
		}
		return nil, serviceerr.RuntimeError(err, s.storage.RunTransaction)
	}

	return &model.SyncPhotoResponse{
		HasBeenUploadedBefore: false,
		Hash:                  newPhoto.Hash,
	}, nil
}
