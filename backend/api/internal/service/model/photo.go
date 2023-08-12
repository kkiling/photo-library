package model

import (
	"github.com/google/uuid"
	"time"
)

type PhotoProcessing string

const (
	PhotoProcessingExifData       PhotoProcessing = "GET_EXIF"
	PhotoProcessingMetaData       PhotoProcessing = "CALC_META"
	PhotoProcessingTagsByCatalogs PhotoExtension  = "TAGS_BY_CATALOGS"
	PhotoProcessingTagsByMeta     PhotoExtension  = "TAGS_BY_META"
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
	PhotoID  uuid.UUID
	Paths    []string
	UploadAt time.Time
	ClientId string
}
