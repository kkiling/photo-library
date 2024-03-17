package photogroup

import (
	"context"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/samber/lo"
	"sync"
)

type Config struct {
	MinSimilarCoefficient float64 `yaml:"min_similar_coefficient"`
}

type Storage interface {
	service.Transactor
	FindSimilarPhotoCoefficients(ctx context.Context, photoID uuid.UUID) ([]model.CoeffSimilarPhoto, error)
	FindGroupIDByPhotoID(ctx context.Context, photoID uuid.UUID) (*uuid.UUID, error)
	SaveGroup(ctx context.Context, group model.PhotoGroup) error
	AddPhotoIDsToGroup(ctx context.Context, groupID uuid.UUID, photoIDs []uuid.UUID) error
}

type Service struct {
	logger  log.Logger
	storage Storage
	cfg     Config
	mu      sync.Mutex
}

func NewService(logger log.Logger, cfg Config, storage Storage) *Service {
	return &Service{
		logger:  logger,
		cfg:     cfg,
		storage: storage,
	}
}

func (s *Service) Init(ctx context.Context) error {
	return nil
}

func (s *Service) NeedLoadPhotoBody() bool {
	return false
}

func (s *Service) getGroupPhotoIDs(ctx context.Context, photoID uuid.UUID) ([]uuid.UUID, error) {
	coefficients, err := s.storage.FindSimilarPhotoCoefficients(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.FindGroupByPhotoID")
	}

	groupPhotoIDs := make([]uuid.UUID, 0)
	groupPhotoIDs = append(groupPhotoIDs, photoID)

	// *** *** *** *** ***
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
		groupPhotoIDs = append(groupPhotoIDs, addPhotoID)
	}
	// *** *** *** *** ***

	return groupPhotoIDs, nil
}

func (s *Service) getGroupAllPhotoIDs(ctx context.Context, photoID uuid.UUID) ([]uuid.UUID, error) {
	alreadyGetGroup := make(map[uuid.UUID]struct{})
	groupIDs, err := s.getGroupPhotoIDs(ctx, photoID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.getGroupPhotoIDs")
	}

	if len(groupIDs) == 1 {
		return groupIDs, nil
	}

	alreadyGetGroup[photoID] = struct{}{}
	haveNewID := true
	for haveNewID {
		haveNewID = false
		for _, id := range groupIDs {
			if _, ok := alreadyGetGroup[id]; ok {
				continue
			}

			groupIDsNew, err := s.getGroupPhotoIDs(ctx, id)
			if err != nil {
				return nil, serviceerr.MakeErr(err, "s.getGroupPhotoIDs")
			}

			alreadyGetGroup[id] = struct{}{}

			for _, idNew := range groupIDsNew {
				if !lo.Contains(groupIDs, idNew) {
					groupIDs = append(groupIDs, idNew)
					haveNewID = true
				}
			}
		}
	}

	return groupIDs, nil
}

// Processing создание и сохранение автоматических тегов (по мета данным или по путям и тд)
func (s *Service) Processing(ctx context.Context, photo model.Photo, _ []byte) (bool, error) {
	// Расчет векторов должен быть не конкурентным, а один за другим
	s.mu.Lock()
	defer s.mu.Unlock()

	groupID, err := s.storage.FindGroupIDByPhotoID(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.FindGroupIDByPhotoID")
	}

	if groupID != nil {
		return true, nil
	}

	groupAllPhotoIDs, err := s.getGroupAllPhotoIDs(ctx, photo.ID)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.getGroupAllPhotoIDs")
	}

	group := model.PhotoGroup{
		ID:          uuid.New(),
		MainPhotoID: photo.ID,
	}

	err = s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		saveErr := s.storage.SaveGroup(ctxTx, group)
		if err != nil {
			return serviceerr.MakeErr(saveErr, "s.storage.CreateGroup")
		}
		if saveErr = s.storage.AddPhotoIDsToGroup(ctxTx, group.ID, groupAllPhotoIDs); saveErr != nil {
			return serviceerr.MakeErr(saveErr, "s.storage.CreateGroup")
		}
		return nil
	})

	if err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.RunTransaction")
	}

	return true, nil
}
