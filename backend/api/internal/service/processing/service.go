package processing

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const (
	defaultTimeout = 12 * time.Second
)

type Config struct {
	MaxGoroutines int  `yaml:"max_goroutines"`
	Limit         int  `yaml:"limit"`
	Debug         bool `yaml:"debug"`
}

type RocketLockService interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (model.RocketLockID, error)
	UnLock(ctx context.Context, lockID model.RocketLockID) error
}

type Storage interface {
	GetPhotoById(ctx context.Context, id uuid.UUID) (model.Photo, error)
	MakeNotValidPhoto(ctx context.Context, photoID uuid.UUID, error string) error
	GetPhotoProcessingTypes(ctx context.Context, photoID uuid.UUID) ([]model.PhotoProcessing, error)
	GetUnprocessedPhotoIDs(ctx context.Context, processingTypes []model.ProcessingType, limit int) ([]uuid.UUID, error)
	AddPhotoProcessing(ctx context.Context, processing model.PhotoProcessing) error
}

type FileStore interface {
	GetFileBody(ctx context.Context, fileKey string) ([]byte, error)
}

type PhotoProcessor interface {
	Init(ctx context.Context) error
	Compensate(ctx context.Context, photoID uuid.UUID) error
	Processing(ctx context.Context, photo model.Photo, photoBody []byte) (bool, error)
	NeedLoadPhotoBody() bool
}

type Service struct {
	logger          log.Logger
	cfg             Config
	storage         Storage
	lock            RocketLockService
	fileStorage     FileStore
	photoProcessors map[model.ProcessingType]PhotoProcessor
}

func NewService(logger log.Logger,
	cfg Config,
	storage Storage,
	fileStorage FileStore,
	lock RocketLockService,
	photoProcessors map[model.ProcessingType]PhotoProcessor,
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
	for _, nextStatus := range model.ProcessingTypes {
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
	if !s.cfg.Debug {
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
	}

	photo, err := s.storage.GetPhotoById(ctx, photoID)
	if err != nil {
		return serviceerr.NotFoundf("photo not found")
	}

	actualPhotoProcessingStatuses, err := s.storage.GetPhotoProcessingTypes(ctx, photoID)
	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.GetPhotoProcessingStatuses")
	}

	var photoBody []byte

	for _, nextType := range model.ProcessingTypes {
		if lo.ContainsBy(actualPhotoProcessingStatuses, func(item model.PhotoProcessing) bool {
			return item.ProcessingType == nextType
		}) {
			continue
		}

		processor, ok := s.photoProcessors[nextType]
		if !ok {
			return serviceerr.NotFoundf("not found processing service for photo status: %s", string(nextType))
		}

		if compensateErr := processor.Compensate(ctx, photoID); compensateErr != nil {
			compensateErr = fmt.Errorf("status %s: %w", nextType, compensateErr)
			return serviceerr.MakeErr(compensateErr, "processor.Compensate")
		}

		// Ленивая загрузка фото, если нужно
		if len(photoBody) == 0 && processor.NeedLoadPhotoBody() {
			photoBody, err = s.fileStorage.GetFileBody(ctx, photo.FileKey)
			if err != nil {
				return serviceerr.MakeErr(err, "s.fileStorage.GetFileBody")
			}
		}

		success, processingErr := processor.Processing(ctx, photo, photoBody)
		if processingErr != nil {
			processingErr = fmt.Errorf("status %s: %w", nextType, processingErr)
			return serviceerr.MakeErr(processingErr, "processor.Processing")
		}

		data := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now(),
			ProcessingType: nextType,
			Success:        success,
		}

		addProcessingErr := s.storage.AddPhotoProcessing(ctx, data)
		if addProcessingErr != nil {
			addProcessingErr = fmt.Errorf("status %s: %w", nextType, addProcessingErr)
			return serviceerr.MakeErr(err, "s.storage.AddPhotoProcessing")
		}
	}

	return nil
}

func (s *Service) ProcessingPhotos(ctx context.Context) (model.PhotoProcessingResult, error) {
	var stats = model.PhotoProcessingResult{
		EOF: true,
	}
	photoIDs, err := s.storage.GetUnprocessedPhotoIDs(ctx, model.ProcessingTypes, s.cfg.Limit)
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
						// Установленна блокировка, просто стоит подождать
						stats.LockProcessedPhotos++
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
