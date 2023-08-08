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
	Id        uuid.UUID
	Url       string
	Hash      string
	UpdateAt  time.Time
	UploadAt  time.Time
	Extension PhotoExtension
}

type UploadPhotoData struct {
	Id       uuid.UUID
	PhotoId  uuid.UUID
	Paths    []string
	UploadAt time.Time
	ClientId string
}
