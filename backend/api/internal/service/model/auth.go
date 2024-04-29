package model

import (
	"time"

	"github.com/google/uuid"
)

// AuthStatus статус авторизации
type AuthStatus string

const (
	// AuthStatusSentInvite был отправлен инвайт на вступление
	AuthStatusSentInvite AuthStatus = "SENT_INVITE"
	// AuthStatusNotActivated не активен
	AuthStatusNotActivated AuthStatus = "NOT_ACTIVATED"
	// AuthStatusActivated активен
	AuthStatusActivated AuthStatus = "ACTIVATED"
	// AuthStatusBlocked заблокирован
	AuthStatusBlocked AuthStatus = "BLOCKED"
)

type AuthRole string

const (
	// AuthRoleAdmin администратор сайта
	AuthRoleAdmin AuthRole = "ADMIN"
	// AuthRoleUser обычный пользователь сайта
	AuthRoleUser AuthRole = "USER"
)

var AuthRoleAll = []AuthRole{AuthRoleAdmin, AuthRoleUser}

type RefreshTokenStatus string

const (
	RefreshTokenStatusActive  RefreshTokenStatus = "ACTIVE"
	RefreshTokenStatusRevoked RefreshTokenStatus = "REVOKED"
	RefreshTokenStatusExpired RefreshTokenStatus = "EXPIRED"
	RefreshTokenStatusLogout  RefreshTokenStatus = "LOGOUT"
)

type Session struct {
	PersonID uuid.UUID `json:"person_id"`
	Role     AuthRole  `json:"role"`
}

type RefreshToken struct {
	Base
	ID       uuid.UUID
	PersonID uuid.UUID
	Status   RefreshTokenStatus
}

type RefreshSession struct {
	RefreshTokenID uuid.UUID `json:"refresh_token_id"`
	PersonID       uuid.UUID `json:"person_id"`
}

// Auth авторизация пользователя
type Auth struct {
	Base
	PersonID     uuid.UUID
	Email        string
	PasswordHash []byte
	Status       AuthStatus
	Role         AuthRole
}

// UpdateAuth Обновление Auth
type UpdateAuth struct {
	BaseUpdate
	PasswordHash UpdateField[[]byte]
	Status       UpdateField[AuthStatus]
	Role         UpdateField[AuthRole]
}

// AuthDataDTO токены авторизации
type AuthDataDTO struct {
	PersonID               uuid.UUID
	Email                  string
	AccessToken            string
	AccessTokenExpiration  time.Time
	RefreshToken           string
	RefreshTokenExpiration time.Time
}
