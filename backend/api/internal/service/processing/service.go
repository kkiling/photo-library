package processing

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/samber/lo"
)

const (
	defaultTimeout = 12 * time.Second
)

type Config struct {
	MaxGoroutines int    `yaml:"max_goroutines"`
	Limit         uint64 `yaml:"limit"`
}

type RocketLockService interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (*model.RocketLockID, error)
	UnLock(ctx context.Context, lockID *model.RocketLockID) error
}

type Storage interface {
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	GetPhotoProcessingStatuses(ctx context.Context, photoID uuid.UUID) ([]model.PhotoProcessingStatus, error)
	GetUnprocessedPhotoIDs(ctx context.Context, lastProcessingStatus model.PhotoProcessingStatus, perPage uint64) ([]uuid.UUID, error)
	AddPhotosProcessingStatus(ctx context.Context, photoID uuid.UUID, status model.PhotoProcessingStatus, success bool) error
	MakeNotValidPhoto(ctx context.Context, photoID uuid.UUID, error string) error
}

type FileStore interface {
	GetFileBody(ctx context.Context, fileName string) ([]byte, error)
}

type PhotoProcessor interface {
	Init(ctx context.Context) error
	Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error)
	NeedLoadPhotoBody() bool
}

type Service struct {
	logger          log.Logger
	cfg             Config
	storage         Storage
	lock            RocketLockService
	fileStorage     FileStore
	photoProcessors map[model.PhotoProcessingStatus]PhotoProcessor
}

func NewService(logger log.Logger,
	cfg Config,
	storage Storage,
	fileStorage FileStore,
	lock RocketLockService,
	photoProcessors map[model.PhotoProcessingStatus]PhotoProcessor,
) *Service {
	return &Service{
		logger:          logger,
		cfg:             cfg,
		storage:         storage,
		lock:            lock,
		fileStorage:     fileStorage,
		photoProcessors: photoProcessors,
	}
}

func (s *Service) Init(ctx context.Context) error {
	for _, nextStatus := range model.PhotoProcessingStatuses {
		processor, ok := s.photoProcessors[nextStatus]
		if !ok {
			return serviceerr.NotFoundf("not found processing service for photo status: %s", string(nextStatus))
		}

		if err := processor.Init(ctx); err != nil {
			err = fmt.Errorf("status %s: %w", nextStatus, err)
			return serviceerr.MakeErr(err, "processor.Init")
		}
	}

	return nil
}

func (s *Service) processingPhoto(ctx context.Context, photoID uuid.UUID) error {
	// Тут нужно ставить лок на обработку фотографии
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	key := fmt.Sprintf("processing_photo_%s", photoID.String())
	lockID, err := s.lock.Lock(ctx, key, defaultTimeout)
	if err != nil {
		return err
	}
	defer func() {
		err := s.lock.UnLock(ctx, lockID)
		if err != nil {
			s.logger.Errorf("unlock: %v", err)
		}
	}()

	photo, err := s.storage.GetPhotoById(ctx, photoID)
	if err != nil {
		return serviceerr.NotFoundf("photo not found")
	}

	actualPhotoProcessingStatuses, err := s.storage.GetPhotoProcessingStatuses(ctx, photoID)
	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.GetPhotoProcessingStatuses")
	}

	var photoBody []byte

	for _, nextStatus := range model.PhotoProcessingStatuses {
		if lo.Contains(actualPhotoProcessingStatuses, nextStatus) {
			continue
		}

		processor, ok := s.photoProcessors[nextStatus]
		if !ok {
			return serviceerr.NotFoundf("not found processing service for photo status: %s", string(nextStatus))
		}

		// Ленивая загрузка фото, если нужно
		if len(photoBody) == 0 && processor.NeedLoadPhotoBody() {
			photoBody, err = s.fileStorage.GetFileBody(ctx, photo.FileName)
			if err != nil {
				return serviceerr.MakeErr(err, "s.fileStorage.GetFileBody")
			}
		}

		// TODO: С учетом того что группы постоянно меняются, сливаются и тд, пока не закончен процесс обработки фотографий
		// Группы этих фотографий нельзя отдавать пользователям
		// TODO: МОЖЕТ БЫТЬ ТАКАЯ СИТУАЦИЯ ЧТО Processing сохранит данные
		// Но при этом AddPhotosProcessingStatus не обновит состояние
		// Либо в одной транзакции, либо все Processing должны быть идемпотентными
		success, err := processor.Processing(ctx, *photo, photoBody)
		if err != nil {
			err = fmt.Errorf("status %s: %w", nextStatus, err)
			return serviceerr.MakeErr(err, "processor.Processing")
		}

		if err := s.storage.AddPhotosProcessingStatus(ctx, photo.ID, nextStatus, success); err != nil {
			err = fmt.Errorf("status %s: %w", nextStatus, err)
			return serviceerr.MakeErr(err, "s.storage.UpdatePhotosProcessingStatus")
		}
	}

	return nil
}

func (s *Service) ProcessingPhotos(ctx context.Context) (model.ProcessingPhotos, error) {
	var stats = model.ProcessingPhotos{
		EOF: true,
	}
	photoIDs, err := s.storage.GetUnprocessedPhotoIDs(ctx, model.LastProcessingStatus, s.cfg.Limit)
	if err != nil {
		return stats, serviceerr.MakeErr(err, "s.storage.GetUnprocessedPhotoIDs")
	}

	if len(photoIDs) == 0 {
		return stats, nil
	}

	var photoIDsChan = make(chan uuid.UUID)
	go func() {
		defer close(photoIDsChan)
		for _, id := range photoIDs {
			photoIDsChan <- id
		}
	}()

	var returnErr error
	var mu = sync.Mutex{}
	var wg sync.WaitGroup

	for i := 0; i < s.cfg.MaxGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ctx.Err() == nil {
				mu.Lock()
				if returnErr != nil {
					mu.Unlock()
					break
				}
				mu.Unlock()

				select {
				case <-ctx.Done():
					// Время выполнения истекло
					return
				case photoID, ok := <-photoIDsChan:
					if !ok {
						return // Канал закрыт
					}

					processedErr := s.processingPhoto(ctx, photoID)
					mu.Lock()
					if processedErr == nil {
						stats.SuccessProcessedPhotos++
						mu.Unlock()
						continue
					}

					if errors.Is(processedErr, serviceerr.ErrAlreadyLocked) {
						stats.LockProcessedPhotos++
						// Установленна блокировка, просто стоит подождать
						// s.logger.Errorf("ErrAlreadyLocked: %v: %v\n", photoID, processedErr)
					} else if errors.Is(processedErr, serviceerr.ErrPhotoIsNotValid) { // Необходимо сделать фото невалидным
						if notValidErr := s.storage.MakeNotValidPhoto(ctx, photoID, processedErr.Error()); notValidErr != nil {
							returnErr = notValidErr
							stats.ErrorProcessedPhotos++
							mu.Unlock()
							return // Критичная ошибка, выходим
						}
					} else { // Неизвестная критичная ошибка, выходим
						stats.ErrorProcessedPhotos++
						returnErr = serviceerr.MakeErr(processedErr, "s.processingPhoto")
						s.logger.Errorf("%v: %v\n", photoID, processedErr)
					}
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	stats.EOF = false
	return stats, returnErr
}
