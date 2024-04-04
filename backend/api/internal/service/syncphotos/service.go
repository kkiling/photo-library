package syncphotos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/utils"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
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

func (s *Service) UploadPhoto(ctx context.Context, form *SyncPhotoRequest) (*SyncPhotoResponse, error) {
	// Проверяем загружено ли фото
	findPhoto, err := s.storage.GetPhotoByHash(ctx, form.Hash)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetPhotoByHash")
	}

	if findPhoto != nil {
		return &SyncPhotoResponse{
			HasBeenUploadedBefore: true,
			Hash:                  findPhoto.Hash,
		}, nil
	}

	if len(form.Paths) == 0 {
		return nil, serviceerr.InvalidInputf("paths must not be empty")
	}

	// Проверка, поддерживается ли расширение
	ex := utils.GetPhotoExtension(form.Paths[0])
	if ex == nil {
		return nil, serviceerr.InvalidInputf("photo extension not found")
	}

	photoID := uuid.New()
	uploadAt := time.Now()

	// Сохранить файл и получить url
	fileName := fmt.Sprintf("%s.%s", photoID, strings.ToLower(string(*ex)))
	if err := s.fileStorage.SaveFileBody(ctx, fileName, form.Body); err != nil {
		return nil, serviceerr.MakeErr(err, "s.fileStorage.SaveFileBody")
	}

	newPhoto := model.Photo{
		ID:        photoID,
		FileName:  fileName,
		Hash:      form.Hash,
		UpdateAt:  form.UpdateAt,
		Extension: *ex,
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
		return nil, serviceerr.MakeErr(err, "s.storage.RunTransaction")
	}

	return &SyncPhotoResponse{
		HasBeenUploadedBefore: false,
		Hash:                  newPhoto.Hash,
	}, nil
}
