package model

import (
	"github.com/google/uuid"
	"time"
)

type ApiTokenType string

const (
	ApiTokenSyncPhoto ApiTokenType = "sync_photo"
)

type ApiToken struct {
	Base
	ID        uuid.UUID
	PersonID  uuid.UUID
	Caption   string
	Token     string
	Type      ApiTokenType
	ExpiredAt *time.Time
}

type ApiTokenDTO struct {
	ID        uuid.UUID
	Caption   string
	Type      ApiTokenType
	ExpiredAt *time.Time
}
