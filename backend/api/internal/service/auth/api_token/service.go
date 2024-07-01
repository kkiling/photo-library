package api_token

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/form"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/utils"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/samber/lo"
	"time"
)

type Storage interface {
	service.Transactor
	GetApiTokens(ctx context.Context, personID uuid.UUID) ([]model.ApiToken, error)
	SaveApiToken(ctx context.Context, apiToken model.ApiToken) error
	DeleteApiToken(ctx context.Context, personID, tokenID uuid.UUID) error
	GetApiToken(ctx context.Context, token string) (model.ApiToken, error)
}

type Service struct {
	logger   log.Logger
	storage  Storage
	validate *validator.Validate
}

func NewService(logger log.Logger, storage Storage) *Service {
	return &Service{
		logger:   logger,
		storage:  storage,
		validate: utils.NewValidator(),
	}
}

func (s *Service) GetApiTokens(ctx context.Context) ([]model.ApiTokenDTO, error) {
	session, err := utils.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	tokens, err := s.storage.GetApiTokens(ctx, session.PersonID)
	if err != nil {
		return nil, serviceerr.MakeErr(err, "s.storage.GetApiTokens")
	}

	return lo.Map(tokens, func(item model.ApiToken, _ int) model.ApiTokenDTO {
		return model.ApiTokenDTO{
			ID:        item.ID,
			Caption:   item.Caption,
			Type:      item.Type,
			ExpiredAt: item.ExpiredAt,
		}
	}), nil
}

func (s *Service) CreateApiToken(ctx context.Context, req form.CreateApiToken) (string, error) {
	session, err := utils.GetSession(ctx)
	if err != nil {
		return "", err
	}

	if err = s.validate.Struct(req); err != nil {
		return "", serviceerr.InvalidInputErr(err, "error in create api token")
	}

	apiToken, err := utils.GenerateAPIToken()
	if err != nil {
		return "", serviceerr.MakeErr(err, "utils.GenerateAPIToken")
	}

	var expiresAt *time.Time
	if req.TimeDuration != nil {
		expiresAt = lo.ToPtr(time.Now().Add(*req.TimeDuration))
	}

	token := model.ApiToken{
		Base:      model.NewBase(),
		ID:        uuid.New(),
		PersonID:  session.PersonID,
		Caption:   req.Caption,
		Token:     apiToken,
		Type:      req.Type,
		ExpiredAt: expiresAt,
	}

	err = s.storage.SaveApiToken(ctx, token)
	if err != nil {
		return "", serviceerr.MakeErr(err, " s.storage.SaveApiToken")
	}

	return "", nil
}

func (s *Service) DeleteApiToken(ctx context.Context, tokenID uuid.UUID) error {
	session, err := utils.GetSession(ctx)
	if err != nil {
		return err
	}

	err = s.storage.DeleteApiToken(ctx, session.PersonID, tokenID)
	if err != nil {
		return serviceerr.MakeErr(err, " s.storage.DeleteApiToken")
	}

	return nil
}

func (s *Service) GetApiToken(ctx context.Context, token string) (model.ApiToken, error) {
	apiToken, err := s.storage.GetApiToken(ctx, token)
	if err != nil {
		if errors.Is(err, serviceerr.ErrNotFound) {
			return model.ApiToken{}, serviceerr.NotFoundf("api token not found")
		}
		return model.ApiToken{}, serviceerr.MakeErr(err, " s.storage.DeleteApiToken")
	}
	return apiToken, nil
}
