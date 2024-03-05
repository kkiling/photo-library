package processing

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"sync"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Config struct {
	MaxGoroutines int   `yaml:"max_goroutines"`
	Limit         int64 `yaml:"limit"`
}

type Storage interface {
	GetPhotoById(ctx context.Context, id uuid.UUID) (*model.Photo, error)
	GetPhotoProcessingStatuses(ctx context.Context, photoID uuid.UUID) ([]model.PhotoProcessingStatus, error)
	GetUnprocessedPhotoIDs(ctx context.Context, lastProcessingStatus model.PhotoProcessingStatus, limit int64) ([]uuid.UUID, error)
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
	fileStorage     FileStore
	photoProcessors map[model.PhotoProcessingStatus]PhotoProcessor
}

func NewService(logger log.Logger,
	cfg Config,
	storage Storage,
	fileStorage FileStore,
	photoProcessors map[model.PhotoProcessingStatus]PhotoProcessor,
) *Service {
	return &Service{
		logger:          logger,
		cfg:             cfg,
		storage:         storage,
		fileStorage:     fileStorage,
		photoProcessors: photoProcessors,
	}
}

func (s *Service) Init(ctx context.Context) error {
	for _, nextStatus := range model.PhotoProcessingStatuses {
		processor, ok := s.photoProcessors[nextStatus]
		if !ok {
			return serviceerr.NotFoundError("not found processing service for photo status: %s", string(nextStatus))
		}

		if err := processor.Init(ctx); err != nil {
			return serviceerr.MakeErr(fmt.Errorf("status %s: %w", nextStatus, err), "processor.Init")
		}
	}

	return nil
}

func (s *Service) processingPhoto(ctx context.Context, photoID uuid.UUID) error {
	// TODO: Тут нужно ставить лок на обработку фотографии
	photo, err := s.storage.GetPhotoById(ctx, photoID)
	if err != nil {
		return serviceerr.NotFoundError("photo not found")
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
			return serviceerr.NotFoundError("not found processing service for photo status: %s", string(nextStatus))
		}

		// Ленивая загрузка фото, если нужно
		if len(photoBody) == 0 && processor.NeedLoadPhotoBody() {
			photoBody, err = s.fileStorage.GetFileBody(ctx, photo.FileName)
			if err != nil {
				return serviceerr.MakeErr(err, "s.fileStorage.GetFileBody")
			}
		}

		success, err := processor.Processing(ctx, *photo, photoBody)
		if err != nil {
			return serviceerr.MakeErr(fmt.Errorf("status %s: %w", nextStatus, err), "processor.Processing")
		}

		if err := s.storage.AddPhotosProcessingStatus(ctx, photo.ID, nextStatus, success); err != nil {
			return serviceerr.MakeErr(fmt.Errorf("status %s: %w", nextStatus, err), "s.storage.UpdatePhotosProcessingStatus")
		}
	}

	return nil
}

func (s *Service) ProcessingPhotos(ctx context.Context) (bool, error) {
	photoIDs, err := s.storage.GetUnprocessedPhotoIDs(ctx, model.LastProcessingStatus, s.cfg.Limit)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.GetUnprocessedPhotoIDs")
	}

	if len(photoIDs) == 0 {
		return false, nil
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
					if processedErr == nil {
						continue
					}

					mu.Lock()
					// Необходимо сделать фото невалидным
					if errors.Is(processedErr, serviceerr.ErrPhotoIsNotValid) {
						if notValidErr := s.storage.MakeNotValidPhoto(ctx, photoID, processedErr.Error()); notValidErr != nil {
							returnErr = notValidErr
							return
						}
					} else {
						returnErr = serviceerr.MakeErr(processedErr, "s.processingPhoto")
						s.logger.Errorf("%v: %v\n", photoID, processedErr)
					}
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	return true, returnErr
}
