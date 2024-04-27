package codes

import (
	"time"

	"github.com/google/uuid"
)

// ConfirmCodeType тип кода подтверждения
type ConfirmCodeType string

const (
	// ConfirmCodeTypeActivateAuth активация авторизации
	ConfirmCodeTypeActivateAuth ConfirmCodeType = "ACTIVATE_AUT"
)

// ConfirmCode код подтверждения
type ConfirmCode struct {
	Code      string
	PersonID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      ConfirmCodeType
	Active    bool
}

// UpdateConfirmCode Обновление Person
type UpdateConfirmCode struct {
	UpdatedAt time.Time
	Type      ConfirmCodeType
	Active    bool
}
