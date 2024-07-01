package model

import (
	"github.com/google/uuid"
)

// ConfirmCodeType тип кода подтверждения
type ConfirmCodeType string

const (
	// ConfirmCodeTypeActivateInvite активация инвайта
	ConfirmCodeTypeActivateInvite ConfirmCodeType = "ACTIVATE_INVITE"
	// ConfirmCodeTypeActivateRegistration активация регистрации
	ConfirmCodeTypeActivateRegistration ConfirmCodeType = "ACTIVATE_REGISTRATION"
)

// ConfirmCode код подтверждения
type ConfirmCode struct {
	Base
	Code     string
	PersonID uuid.UUID
	Type     ConfirmCodeType
	Active   bool
}

// UpdateConfirmCode Обновление Person
type UpdateConfirmCode struct {
	BaseUpdate
	Active UpdateField[bool]
}
