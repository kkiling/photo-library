package form

import "github.com/kkiling/photo-library/backend/api/internal/service/model"

type AdminInitInviteForm struct {
	Email string `validate:"required,email"`
}

// SendInviteForm отправление инвайта на присоединение в системе
type SendInviteForm struct {
	Email string         `validate:"required,email"`
	Role  model.AuthRole `validate:"required"`
}

// ActivateInviteForm активация инвайта
type ActivateInviteForm struct {
	FirstName   string  `validate:"required,min=3,max=128"`
	Surname     string  `validate:"required,min=3,max=128"`
	Patronymic  *string `validate:"omitempty,min=3,max=128"`
	CodeConfirm string  `validate:"required,max=128"`
	Password    string  `validate:"required,min=8,containsany=0123456789,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type LoginForm struct {
	Email    string `validate:"required,max=128"`
	Password string `validate:"required,max=128"`
}

type RegisterForm struct {
	FirstName  string  `validate:"required,min=3,max=128"`
	Surname    string  `validate:"required,min=3,max=128"`
	Patronymic *string `validate:"omitempty,min=3,max=128"`
	Email      string  `validate:"required,max=128"`
	Password   string  `validate:"required,max=128"`
}

// ActivateRegisterForm активация регистрации
type ActivateRegisterForm struct {
	CodeConfirm string `validate:"required,max=128"`
}

type EmailAvailableForm struct {
	Email string `validate:"required,max=128"`
}
