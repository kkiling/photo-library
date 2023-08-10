package syncphotos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"path/filepath"
	"strings"
	"time"
)

type Database interface {
	service.Transactor
	GetPhotoByHash(ctx context.Context, hash string) (*model.Photo, error)
	SavePhoto(ctx context.Context, photo model.Photo) error
	SaveUploadPhotoData(ctx context.Context, data model.UploadPhotoData) error
}

type FileStore interface {
	SaveFileBody(ctx context.Context, ext string, body []byte) (filePath string, err error)
	DeleteFile(ctx context.Context, filePath string) error
}

type Service struct {
	logger      log.Logger
	database    Database
	fileStorage FileStore
}

func NewService(logger log.Logger, storage Database, fileStorage FileStore) *Service {
	return &Service{
		logger:      logger,
		database:    storage,
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

func (s *Service) UploadPhoto(ctx context.Context, form model.SyncPhotoRequest) (model.SyncPhotoResponse, error) {

	// Проверяем загружено ли фото
	findPhoto, err := s.database.GetPhotoByHash(ctx, form.Hash)
	if err != nil {
		return model.SyncPhotoResponse{}, fmt.Errorf("database.GetPhotoByHash: %w", err)
	}

	if findPhoto != nil {
		return model.SyncPhotoResponse{
			HasBeenUploadedBefore: true,
			Hash:                  findPhoto.Hash,
			UploadAt:              findPhoto.UploadAt,
		}, nil
	}

	if len(form.Paths) == 0 {
		// TODO:  ошибка
		return model.SyncPhotoResponse{}, fmt.Errorf("paths must not be empty")
	}

	// Проверка, поддерживается ли расширение
	ex := s.getPhotoExtension(form.Paths[0])
	if ex == nil {
		// TODO:  ошибка
		return model.SyncPhotoResponse{}, fmt.Errorf("photo extension not found")
	}

	// Сохранить файл и получить url
	filePath, err := s.fileStorage.SaveFileBody(ctx, string(*ex), form.Body)
	if err != nil {
		// TODO:  ошибка
		return model.SyncPhotoResponse{}, fmt.Errorf("fileStorage.SaveFileBody: %w", err)
	}

	newPhoto := model.Photo{
		ID:        uuid.New(),
		FilePath:  filePath,
		Hash:      form.Hash,
		UpdateAt:  form.UpdateAt,
		UploadAt:  time.Now(),
		Extension: *ex,
	}

	uploadPhotoData := model.UploadPhotoData{
		ID:       uuid.New(),
		PhotoID:  newPhoto.ID,
		Paths:    form.Paths,
		UploadAt: newPhoto.UploadAt,
		ClientId: form.ClientId,
	}

	// Если произошла ошибка при сохранении фото в базе
	// То, пытаемся удалить его из файлового хранилища
	// Если программа упадет до удаления из хранилища, то файл зависнет и будет бесхозным
	// Нужна чистилка бесхозных файлов

	// Одной транзакцией сохранить
	err = s.database.RunTransaction(ctx, func(ctxTx context.Context) error {
		if saveErr := s.database.SavePhoto(ctxTx, newPhoto); saveErr != nil {
			return saveErr
		}
		if saveErr := s.database.SaveUploadPhotoData(ctxTx, uploadPhotoData); saveErr != nil {
			return saveErr
		}
		return nil
	})

	if err != nil {
		if delErr := s.fileStorage.DeleteFile(ctx, newPhoto.FilePath); delErr != nil {
			s.logger.Errorf("fail fileStorage.DeleteFile %s: %w", newPhoto.FilePath, delErr)
		}
		// TODO:  ошибка
		return model.SyncPhotoResponse{}, fmt.Errorf("database.WithTransaction: %w", err)
	}

	return model.SyncPhotoResponse{
		HasBeenUploadedBefore: false,
		Hash:                  newPhoto.Hash,
		UploadAt:              newPhoto.UploadAt,
	}, nil
}
