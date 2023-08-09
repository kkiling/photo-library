package model

import (
	"github.com/google/uuid"
	"time"
)

type PhotoExtension string

const (
	PhotoExtensionJpg  PhotoExtension = "JPG"
	PhotoExtensionJpeg PhotoExtension = "JPEG"
	PhotoExtensionPng  PhotoExtension = "PNG"
	PhotoExtensionBmb  PhotoExtension = "BMP"
)

type Photo struct {
	ID        uuid.UUID
	FilePath  string
	Hash      string
	UpdateAt  time.Time
	UploadAt  time.Time
	Extension PhotoExtension
}

type UploadPhotoData struct {
	ID       uuid.UUID
	PhotoID  uuid.UUID
	Paths    []string
	UploadAt time.Time
	ClientId string
}
