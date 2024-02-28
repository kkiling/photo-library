package processing

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Config struct {
	MaxGoroutines int `yaml:"max_goroutines"`
}

type Storage interface {
	GetPhotosCount(ctx context.Context, filter *model.PhotoFilter) (int64, error)
	GetPaginatedPhotos(ctx context.Context, params model.PhotoSelectParams, filter *model.PhotoFilter) ([]model.Photo, error)
	UpdatePhotosProcessingStatus(ctx context.Context, id uuid.UUID, newProcessingStatus model.PhotoProcessingStatus) error
}

type FileStore interface {
	GetFileBody(ctx context.Context, fileName string) ([]byte, error)
}

type PhotoProcessor interface {
	Processing(ctx context.Context, photo model.Photo, photoBody []byte) error
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

func (s *Service) processingPhoto(ctx context.Context, photo model.Photo) error {
	nextStatus := model.NewPhoto

	switch photo.ProcessingStatus {
	case model.NewPhoto:
		nextStatus = model.ExifDataSaved
	case model.ExifDataSaved:
		nextStatus = model.MetaDataSaved
	case model.MetaDataSaved:
		nextStatus = model.SystemTagsSaved
	case model.SystemTagsSaved:
		nextStatus = model.PhotoVectorSaved
	case model.PhotoVectorSaved:
		// Конечный стутус в данный момент
		return nil
	}

	processor, ok := s.photoProcessors[nextStatus]
	if !ok {
		return serviceerr.NotFoundError("not found processing service for photo status: %s", string(nextStatus))
	}

	photoBody, err := s.fileStorage.GetFileBody(ctx, photo.FileName)
	if err != nil {
		return serviceerr.RuntimeError(err, s.fileStorage.GetFileBody)
	}

	if err := processor.Processing(ctx, photo, photoBody); err != nil {
		return serviceerr.RuntimeError(err, processor.Processing)
	}

	if err := s.storage.UpdatePhotosProcessingStatus(ctx, photo.ID, nextStatus); err != nil {
		return serviceerr.RuntimeError(err, s.storage.UpdatePhotosProcessingStatus)
	}

	return nil
}

func (s *Service) ProcessingPhotos(ctx context.Context,
	status model.PhotoProcessingStatus, limit int) (processedCount int, totalCount int64, err error) {

	filter := model.PhotoFilter{
		ProcessingStatusIn: []model.PhotoProcessingStatus{status},
	}
	photos, err := s.storage.GetPaginatedPhotos(ctx, model.PhotoSelectParams{Limit: limit}, &filter)

	if err != nil {
		return 0, 0, serviceerr.RuntimeError(err, s.storage.GetPaginatedPhotos)
	}

	totalCount, err = s.storage.GetPhotosCount(ctx, &filter)
	if err != nil {
		return 0, 0, serviceerr.RuntimeError(err, s.storage.GetPhotosCount)
	}

	photoChan := make(chan model.Photo)
	errorsChan := make(chan error, s.cfg.MaxGoroutines)
	go func() {
		for _, photo := range photos {
			photoChan <- photo
		}
		close(photoChan)
	}()

	var wg sync.WaitGroup
	for i := 0; i < s.cfg.MaxGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ctx.Err() == nil {
				select {
				case <-ctx.Done():
					// Время выполнения истекло
					errorsChan <- ctx.Err()
					return
				case photo, ok := <-photoChan:
					// Получили результат
					if !ok {
						return
					}
					if processedErr := s.processingPhoto(ctx, photo); processedErr != nil {
						if errors.Is(processedErr, serviceerr.ErrPhotoIsNotValid) {
							// TODO: пометить что фото не валидно
							fmt.Println(processedErr.Error())
						}
						errorsChan <- serviceerr.RuntimeError(processedErr, s.processingPhoto)
						// TODO: остановить все гроутины
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		return 0, 0, <-errorsChan
	}

	return len(photos), totalCount, nil
}
