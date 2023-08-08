package syncphotos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"path/filepath"
	"strings"
	"time"
)

type Transactor interface {
	WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type Storage interface {
	Transactor
	GetPhotoByHash(ctx context.Context, hash string) (*model.Photo, error)
	SavePhoto(ctx context.Context, photo model.Photo) error
	SaveUploadPhotoData(ctx context.Context, data model.UploadPhotoData) error
}

type FileStore interface {
	SaveFileBody(ctx context.Context, body []byte) (url string, err error)
	DeleteFile(ctx context.Context, url string) error
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

func (s *Service) UploadPhoto(ctx context.Context, form model.SyncPhotoRequest) (model.SyncPhotoResponse, error) {

	// Проверяем загружено ли фото
	photo, err := s.storage.GetPhotoByHash(ctx, form.Hash)
	if err != nil {
		return model.SyncPhotoResponse{}, fmt.Errorf("storage.GetPhotoByHash: %w", err)
	}

	if photo != nil {
		return model.SyncPhotoResponse{
			HasBeenUploadedBefore: true,
			Hash:                  photo.Hash,
			UploadAt:              photo.UploadAt,
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
	url, err := s.fileStorage.SaveFileBody(ctx, form.Body)
	if err != nil {
		// TODO:  ошибка
		return model.SyncPhotoResponse{}, fmt.Errorf("fileStorage.SaveFileBody: %w", err)
	}

	newPhoto := model.Photo{
		Id:        uuid.New(),
		Url:       url,
		Hash:      form.Hash,
		UpdateAt:  form.UpdateAt,
		UploadAt:  time.Now(),
		Extension: *ex,
	}

	uploadPhotoData := model.UploadPhotoData{
		Id:       uuid.New(),
		PhotoId:  newPhoto.Id,
		Paths:    form.Paths,
		UploadAt: newPhoto.UploadAt,
		ClientId: form.ClientId,
	}

	// Одной транзакцией сохранить
	err = s.storage.WithTransaction(ctx, func(tx context.Context) error {
		if saveErr := s.storage.SavePhoto(tx, newPhoto); saveErr != nil {
			return saveErr
		}
		if saveErr := s.storage.SaveUploadPhotoData(tx, uploadPhotoData); saveErr != nil {
			return saveErr
		}
		return nil
	})

	if err != nil {
		if delErr := s.fileStorage.DeleteFile(ctx, photo.Url); delErr != nil {
			s.logger.Errorf("fail fileStorage.DeleteFile %s: %w", photo.Url, delErr)
		}
		// TODO:  ошибка
		return model.SyncPhotoResponse{}, fmt.Errorf("storage.WithTransaction: %w", err)
	}

	// Если произошла ошибка при сохранении фото в базе
	// То, пытаемся удалить его из файлового хранилища
	// Если программа упадет до удаления из хранилища, то файл зависнет и будет бесхозным
	// Нужна чистилка бесхозных файлов

	/*// Извлекаем расширение файла
	ext := filepath.Ext(request.Paths[0])
	// Создаем UUID
	uuid := uuid.New()

	// Формируем новое имя файла
	newFilename := fmt.Sprintf("/Users/kkiling/Desktop/syncphotos/%s%s", uuid.String(), ext)

	// Создаем новый файл с новым именем
	newFile, err := os.Create(newFilename)
	defer newFile.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new file: %v", err)
	}

	// Записываем данные в новый файл
	if _, err := newFile.Write(request.Body); err != nil {
		return nil, fmt.Errorf("Failed to write to new file: %v", err)
	}
	*/
	return model.SyncPhotoResponse{
		HasBeenUploadedBefore: false,
		Hash:                  photo.Hash,
		UploadAt:              photo.UploadAt,
	}, nil
}
