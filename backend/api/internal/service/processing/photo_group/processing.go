package photo_group

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

const (
	defaultTimeout = 6 * time.Second
)

type Config struct {
	MinSimilarCoefficient float64 `yaml:"min_similar_coefficient"`
	Debug                 bool    `yaml:"debug"`
}

type RocketLockService interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (model.RocketLockID, error)
	UnLock(ctx context.Context, lockID model.RocketLockID) error
}

type Storage interface {
	service.Transactor
	FindCoefficientSimilarPhoto(ctx context.Context, photoID uuid.UUID) ([]model.CoefficientSimilarPhoto, error)
	FindGroupIDByPhotoID(ctx context.Context, photoID uuid.UUID) (uuid.UUID, error)
	SaveGroup(ctx context.Context, group model.PhotoGroup) error
	AddPhotoIDsToGroup(ctx context.Context, groupID uuid.UUID, photoIDs []uuid.UUID) error
	GetGroupPhotoIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error)
	DeletePhotoGroup(ctx context.Context, groupID uuid.UUID) error
	DeletePhotoGroupByPhotoID(ctx context.Context, photoID uuid.UUID) error
}

type Processing struct {
	logger  log.Logger
	storage Storage
	lock    RocketLockService
	cfg     Config
	mu      sync.Mutex
}

func NewService(logger log.Logger, cfg Config, storage Storage, lock RocketLockService) *Processing {
	return &Processing{
		logger:  logger,
		cfg:     cfg,
		lock:    lock,
		storage: storage,
	}
}

func (s *Processing) Compensate(ctx context.Context, photoID uuid.UUID) error {
	err := s.storage.DeletePhotoGroupByPhotoID(ctx, photoID)
	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.DeletePhotoGroupByPhotoID")
	}
	return nil
}

func (s *Processing) Init(_ context.Context) error {
	return nil
}

func (s *Processing) NeedLoadPhotoBody() bool {
	return false
}

func (s *Processing) getSimilarPhotosForPhotoID(ctx context.Context, photoID uuid.UUID) ([]uuid.UUID, error) {
	coefficients, err := s.storage.FindCoefficientSimilarPhoto(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.FindCoefficientSimilarPhoto")
	}

	res := make([]uuid.UUID, 0)
	res = append(res, photoID)

	for _, coefficient := range coefficients {
		if coefficient.Coefficient < s.cfg.MinSimilarCoefficient {
			continue
		}
		var addPhotoID uuid.UUID
		if photoID == coefficient.PhotoID1 {
			addPhotoID = coefficient.PhotoID2
		} else {
			addPhotoID = coefficient.PhotoID1
		}
		res = append(res, addPhotoID)
	}

	return res, nil
}

func (s *Processing) getSimilarPhotos(ctx context.Context, photoID uuid.UUID) ([]uuid.UUID, error) {
	alreadyGetGroup := make(map[uuid.UUID]struct{})

	// Ищем фотографии похожие на указанную
	similarPhotos, err := s.getSimilarPhotosForPhotoID(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.getSimilarPhotosForPhotoID")
	}

	// Если фото одно (те само указанное фото), то возвращаем сразу
	if len(similarPhotos) == 1 {
		return similarPhotos, nil
	}

	// Пробегаемся по всем фото similarPhotos
	// Дла каждого фото ищем свои похожие фотографии
	// Всю эту сеть сливаем в общую similarPhotos

	alreadyGetGroup[photoID] = struct{}{}
	haveNewID := true
	for haveNewID {
		haveNewID = false

		for _, id := range similarPhotos {
			if _, ok := alreadyGetGroup[id]; ok {
				continue
			}

			similarPhotosNew, err := s.getSimilarPhotosForPhotoID(ctx, id)
			if err != nil {
				return nil, serviceerr.MakeErr(err, "s.getSimilarPhotosForPhotoID")
			}

			alreadyGetGroup[id] = struct{}{}

			for _, idNew := range similarPhotosNew {
				if !lo.Contains(similarPhotos, idNew) {
					similarPhotos = append(similarPhotos, idNew)
					haveNewID = true
				}
			}
			if haveNewID {
				break
			}
		}

	}

	return similarPhotos, nil
}

func (s *Processing) getGroupsOfSimilarPhotos(ctx context.Context, similarPhotos []uuid.UUID) ([]uuid.UUID, error) {
	groupsOfSimilarPhotos := make([]uuid.UUID, 0, len(similarPhotos))
	for _, photoID := range similarPhotos {
		groupID, err := s.storage.FindGroupIDByPhotoID(ctx, photoID)
		switch {
		case err == nil:
			groupsOfSimilarPhotos = append(groupsOfSimilarPhotos, groupID)
		case errors.Is(err, serviceerr.ErrNotFound):
			continue
		default:
			return nil, serviceerr.MakeErr(err, "s.storage.FindGroupIDByPhotoID")
		}
	}
	return lo.Uniq(groupsOfSimilarPhotos), nil
}

func (s *Processing) mergePhotoGroups(ctx context.Context, photoID uuid.UUID, similarPhotos, groupsOfSimilarPhotos []uuid.UUID) error {
	// Получаем список всех фотографий из всех групп, сливаем в один список.
	photosToMerge := make([]uuid.UUID, 0)
	for _, groupID := range groupsOfSimilarPhotos {
		ids, err := s.storage.GetGroupPhotoIDs(ctx, groupID)
		if err != nil {
			return serviceerr.MakeErr(err, "s.storage.GetGroupPhotoIDs")
		}
		photosToMerge = append(photosToMerge, ids...)
	}

	photosToMerge = append(photosToMerge, similarPhotos...)
	photosToMerge = lo.Uniq(photosToMerge)

	// Удаляем старые группы.
	// Создаем новую группу, с главной фотографией photo.ID.
	// Добавляем все фотографии в новую группу
	group := model.PhotoGroup{
		ID:          uuid.New(),
		MainPhotoID: photoID,
	}

	err := s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		// Удаляем старые группы
		for _, groupID := range groupsOfSimilarPhotos {
			if deleteErr := s.storage.DeletePhotoGroup(ctxTx, groupID); deleteErr != nil {
				return serviceerr.MakeErr(deleteErr, "s.storage.DeletePhotoGroup")
			}
		}

		if saveErr := s.storage.SaveGroup(ctxTx, group); saveErr != nil {
			return serviceerr.MakeErr(saveErr, "s.storage.CreateGroup")
		}

		if saveErr := s.storage.AddPhotoIDsToGroup(ctxTx, group.ID, photosToMerge); saveErr != nil {
			return serviceerr.MakeErr(saveErr, "s.storage.AddPhotoIDsToGroup")
		}
		return nil
	})

	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.RunTransaction")
	}

	return nil
}

func (s *Processing) createNewPhotoGroup(ctx context.Context, photoID uuid.UUID, groupsOfSimilarPhotos []uuid.UUID) error {
	group := model.PhotoGroup{
		ID:          uuid.New(),
		MainPhotoID: photoID,
	}

	err := s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		if saveErr := s.storage.SaveGroup(ctxTx, group); saveErr != nil {
			return serviceerr.MakeErr(saveErr, "s.storage.CreateGroup")
		}

		if saveErr := s.storage.AddPhotoIDsToGroup(ctxTx, group.ID, groupsOfSimilarPhotos); saveErr != nil {
			return serviceerr.MakeErr(saveErr, "s.storage.AddPhotoIDsToGroup")
		}

		return nil
	})

	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.RunTransaction")
	}

	return nil
}

// Processing создание и сохранение групп фотографий на основании коэффициента одинаковости фотографий
func (s *Processing) Processing(ctx context.Context, photo model.Photo, _ []byte) (bool, error) {
	// Расчет векторов должен быть не конкурентным, а один за другим
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.cfg.Debug {
		ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
		defer cancel()

		// Так же ставим лок (на случай если идет обработка из нескольких подов
		lockID, err := s.lock.Lock(ctx, "similar_coefficient", defaultTimeout)
		if err != nil {
			return false, err
		}
		defer func() {
			if unlockErr := s.lock.UnLock(ctx, lockID); unlockErr != nil {
				s.logger.Errorf("unlock: %v", unlockErr)
			}
		}()
	}

	// Если фотография уже принадлежит к какой-то группе, то ничего не делаем
	_, err := s.storage.FindGroupIDByPhotoID(ctx, photo.ID)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, serviceerr.ErrNotFound):
	default:
		return false, serviceerr.MakeErr(err, "s.storage.FindGroupIDByPhotoID")
	}

	similarPhotos, err := s.getSimilarPhotos(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.getSimilarPhotos")
	}

	// Если фотографии groupAllPhotoIDs уже находятся в одной из групп
	// Точнее все groupAllPhotoIDs должны принадлежать оной группе
	// То, просто добавляем новое photo в уже существующую групп

	// Проверяем к каким группам принадлежат все эти фотографии
	groupsOfSimilarPhotos, err := s.getGroupsOfSimilarPhotos(ctx, similarPhotos)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.getGroupsOfSimilarPhotos")
	}

	// Если оказывается так что какие-то из собранных похожих фотографий уже принадлежат к какой-то группе
	if len(groupsOfSimilarPhotos) == 1 {
		// Если все принадлежат уже к какой-то одной группе
		// То, просто добавляем новую фотографию в эту группу
		if saveErr := s.storage.AddPhotoIDsToGroup(ctx, groupsOfSimilarPhotos[0], []uuid.UUID{photo.ID}); saveErr != nil {
			return false, serviceerr.MakeErr(saveErr, "s.storage.AddPhotoIDsToGroup")
		}
		return true, nil
	} else if len(groupsOfSimilarPhotos) > 1 {

		// Если таких групп несколько
		// То, нужно слить все группы в одну
		if err = s.mergePhotoGroups(ctx, photo.ID, similarPhotos, groupsOfSimilarPhotos); err != nil {
			return false, serviceerr.MakeErr(err, "s.mergePhotoGroups")
		}

		return true, nil
	}

	// Если эти фотографии не принадлежат ни к какой группе, то создаем новую группу
	// И добавляем туда все фотографии
	if err = s.createNewPhotoGroup(ctx, photo.ID, similarPhotos); err != nil {
		return false, serviceerr.MakeErr(err, "s.createNewPhotoGroup")
	}

	return true, nil
}
