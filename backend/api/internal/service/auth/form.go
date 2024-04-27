package auth

// InitAdminForm инициализация администратора
type InitAdminForm struct {
	Name       string  `validate:"required,min=3,max=128"`
	Surname    string  `validate:"required,min=3,max=128"`
	Patronymic *string `validate:"omitempty,min=3,max=128"`
	Email      string  `validate:"required,email"`
}

// ActivateAuthForm активация аккаунта
type ActivateAuthForm struct {
	CodeConfirm string `validate:"required,max=128"`
	Password    string `validate:"required,min=8,containsany=0123456789,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type LoginForm struct {
	Email    string `validate:"required,max=128"`
	Password string `validate:"required,max=128"`
}
