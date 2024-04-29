package codes

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/service/utils"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

// Storage интерфейс хранения данных
type Storage interface {
	service.Transactor
	SaveConfirmCode(ctx context.Context, confirmCode model.ConfirmCode) error
	GetActiveConfirmCode(ctx context.Context, code string, confirmType model.ConfirmCodeType) (model.ConfirmCode, error)
	UpdateConfirmCode(ctx context.Context, personID uuid.UUID, confirmCodeType model.ConfirmCodeType, update model.UpdateConfirmCode) error
}

// Service сервис для работы с кодами повержения
type Service struct {
	logger  log.Logger
	storage Storage
}

// NewService создание сервиса для работы с паролями
func NewService(logger log.Logger, storage Storage) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
	}
}

// GetActiveConfirmCode получение активного кода подтверждения
func (s *Service) GetActiveConfirmCode(ctx context.Context, code string, confirmType model.ConfirmCodeType) (model.ConfirmCode, error) {
	return s.storage.GetActiveConfirmCode(ctx, code, confirmType)
}

// SendConfirmCode отправка кода подтверждения
func (s *Service) SendConfirmCode(ctx context.Context, personID uuid.UUID, confirmType model.ConfirmCodeType) error {
	code, err := utils.GenerateCode()
	if err != nil {
		return serviceerr.MakeErr(err, "utils.GenerateCode")
	}
	confirmCode := model.ConfirmCode{
		Base:     model.NewBase(),
		Code:     code,
		PersonID: personID,
		Type:     confirmType,
		Active:   true,
	}

	// Сначала отправляем, потом сохраняем
	// Не страшно если два рада отправим, страшно если кода не будет в базе
	// TODO: Реализация отправки кодов
	fmt.Println("*** *** *** *** *** *** ***")
	fmt.Printf("Confirm code: %s, for personID: %s", confirmCode.Code, personID.String())
	fmt.Println("*** *** *** *** *** *** ***")

	if err := s.storage.SaveConfirmCode(ctx, confirmCode); err != nil {
		return serviceerr.MakeErr(err, "s.storage.SaveConfirmCode")
	}

	return nil
}

// DeactivateCode деактивация кода подтверждения
func (s *Service) DeactivateCode(ctx context.Context, personID uuid.UUID, confirmCodeType model.ConfirmCodeType) error {
	update := model.UpdateConfirmCode{
		BaseUpdate: model.NewBaseUpdate(),
		Active:     model.NewUpdateField(false),
	}
	if err := s.storage.UpdateConfirmCode(ctx, personID, confirmCodeType, update); err != nil {
		return serviceerr.MakeErr(err, "s.storage.UpdateConfirmCode")
	}
	return nil
}
