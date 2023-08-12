package entity

import (
	"github.com/google/uuid"
	"time"
)

type Photo struct {
	ID        uuid.UUID
	FilePath  string
	Hash      string
	UpdateAt  time.Time
	UploadAt  time.Time
	Extension string
}

type UploadPhotoData struct {
	PhotoID  uuid.UUID
	Paths    []string
	UploadAt time.Time
	ClientId string
}
