package codes

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

// Storage интерфейс хранения данных
type Storage interface {
	service.Transactor
	SaveConfirmCode(ctx context.Context, confirmCode ConfirmCode) error
	GetConfirmCode(ctx context.Context, code string) (ConfirmCode, error)
	UpdateConfirmCode(ctx context.Context, personID uuid.UUID, update UpdateConfirmCode) error
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
func (s *Service) GetActiveConfirmCode(ctx context.Context, code string) (ConfirmCode, error) {
	return s.storage.GetConfirmCode(ctx, code)
}

// SendConfirmCode отправка кода подтверждения
func (s *Service) SendConfirmCode(ctx context.Context, personID uuid.UUID, confirmType ConfirmCodeType) error {
	code := ConfirmCode{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Code:      uuid.NewString(), // TODO: сделать более человеческий код
		PersonID:  personID,
		Type:      confirmType,
		Active:    true,
	}

	// Сначала отправляем, потом сохраняем
	// Не страшно если два рада отправим, страшно если кода не будет в базе
	// TODO: Реализация отправки кодов
	fmt.Println("*** *** *** *** *** *** ***")
	fmt.Printf("Confirm code: %s, for personID: %d", code.Code, personID)
	fmt.Println("*** *** *** *** *** *** ***")

	if err := s.storage.SaveConfirmCode(ctx, code); err != nil {
		return serviceerr.MakeErr(err, "s.storage.SaveConfirmCode")
	}

	return nil
}

// DeactivateCode деактивация кода подтверждения
func (s *Service) DeactivateCode(ctx context.Context, personID uuid.UUID, confirmType ConfirmCodeType) error {
	update := UpdateConfirmCode{
		UpdatedAt: time.Now(),
		Type:      confirmType,
		Active:    false,
	}
	if err := s.storage.UpdateConfirmCode(ctx, personID, update); err != nil {
		return serviceerr.MakeErr(err, "s.storage.UpdateConfirmCode")
	}
	return nil
}
