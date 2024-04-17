package model

import (
	"time"

	"github.com/google/uuid"
)

type PhotoExtension string

const (
	PhotoExtensionJpeg PhotoExtension = "JPEG"
	PhotoExtensionPng  PhotoExtension = "PNG"
)

type Photo struct {
	ID        uuid.UUID
	FileKey   string
	Hash      string
	UpdateAt  time.Time
	Extension PhotoExtension
}

type PhotoUploadData struct {
	PhotoID  uuid.UUID
	UploadAt time.Time
	Paths    []string
	ClientID string
}
