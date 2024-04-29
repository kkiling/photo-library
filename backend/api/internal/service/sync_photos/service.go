package sync_photos

import (
	"context"
	"errors"
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
	GetPhotoByHash(ctx context.Context, hash string) (model.Photo, error)
	SavePhoto(ctx context.Context, photo model.Photo) error
	SavePhotoUploadData(ctx context.Context, uploadData model.PhotoUploadData) error
}

type FileStore interface {
	SaveFile(_ context.Context, fileName string, body []byte, dirs ...string) (fileKey string, err error)
	DeleteFile(_ context.Context, fileKey string) error
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

func (s *Service) UploadPhoto(ctx context.Context, form *model.SyncPhotoRequest) (model.SyncPhotoResponse, error) {
	apiToken, err := utils.GetApiToken(ctx)
	if err != nil {
		return model.SyncPhotoResponse{}, err
	}
	if len(form.Paths) == 0 {
		return model.SyncPhotoResponse{}, serviceerr.InvalidInputf("paths must not be empty")
	}

	// Проверяем загружено ли фото
	findPhoto, err := s.storage.GetPhotoByHash(ctx, form.Hash)
	switch {
	case errors.Is(err, serviceerr.ErrNotFound):
	case err == nil:
		return model.SyncPhotoResponse{
			HasBeenUploadedBefore: true,
			Hash:                  findPhoto.Hash,
		}, nil
	default:
		return model.SyncPhotoResponse{}, serviceerr.MakeErr(err, "s.storage.GetPhotoByHash")
	}

	// Проверка, поддерживается ли расширение
	ex := utils.GetPhotoExtension(form.Paths[0])
	if ex == nil {
		return model.SyncPhotoResponse{}, serviceerr.InvalidInputf("photo extension not found")
	}

	photoID := uuid.New()

	fileName := fmt.Sprintf("%s.%s", photoID, strings.ToLower(string(*ex)))
	fileKey, err := s.fileStorage.SaveFile(ctx, fileName, form.Body)
	if err != nil {
		return model.SyncPhotoResponse{}, serviceerr.MakeErr(err, "s.fileStorage.SaveFileBody")
	}

	// TODO: добавление в photo space

	newPhoto := model.Photo{
		ID:             photoID,
		FileKey:        fileKey,
		Hash:           form.Hash,
		PhotoUpdatedAt: form.PhotoUpdatedAt,
		Extension:      *ex,
	}

	uploadPhotoData := model.PhotoUploadData{
		PhotoID:    photoID,
		UploadAt:   time.Now(),
		Paths:      form.Paths,
		ClientInfo: form.ClientInfo,
		PersonID:   apiToken.PersonID,
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
		if saveErr := s.storage.SavePhotoUploadData(ctxTx, uploadPhotoData); saveErr != nil {
			return saveErr
		}
		return nil
	})

	if err != nil {
		if delErr := s.fileStorage.DeleteFile(ctx, newPhoto.FileKey); delErr != nil {
			s.logger.Errorf("fail fileStorage.DeleteFile %s: %v", newPhoto.FileKey, delErr)
		}
		return model.SyncPhotoResponse{}, serviceerr.MakeErr(err, "s.storage.RunTransaction")
	}

	return model.SyncPhotoResponse{
		HasBeenUploadedBefore: false,
		Hash:                  newPhoto.Hash,
	}, nil
}
