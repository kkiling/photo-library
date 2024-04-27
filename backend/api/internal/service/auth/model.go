package auth

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// AuthStatus статус авторизации
type AuthStatus string

const (
	// AuthStatusNotActivated не активен
	AuthStatusNotActivated AuthStatus = "NOT_ACTIVATED"
	// AuthStatusActivated активен
	AuthStatusActivated AuthStatus = "ACTIVATED"
	// AuthStatusBlocked заблокирован
	AuthStatusBlocked AuthStatus = "BLOCKED"
)

type Session struct {
	PersonID uuid.UUID `json:"person_id"`
}

// Person информация о человеке
type Person struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	Surname    string
	Patronymic *string
}

// FullName полное имя человека
func (p Person) FullName() string {
	return fmt.Sprintf("%s %s", p.Name, p.Surname)
}

// Auth авторизация пользователя
type Auth struct {
	PersonID     uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	PasswordHash []byte
	Status       AuthStatus
}

// PersonFull человек с полной информацией о нем
type PersonFull struct {
	Person
	Auth *Auth
}

// AuthData токены авторизации
type AuthData struct {
	Email                  string
	AccessToken            string
	AccessTokenExpiration  time.Time
	RefreshToken           string
	RefreshTokenExpiration time.Time
}

// UpdateAuth Обновление Auth
type UpdateAuth struct {
	UpdatedAt    time.Time
	PasswordHash []byte
	Status       AuthStatus
}
