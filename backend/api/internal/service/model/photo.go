package model

import (
	"time"

	"github.com/google/uuid"
)

type PhotoFilter struct {
	ProcessingStatusIn []PhotoProcessingStatus
}

type PhotoSortOrder string

const (
	PhotoSortOrderNone PhotoSortOrder = "NONE"
)

type PhotoSelectParams struct {
	Offset     int64
	Limit      int
	SortOrder  PhotoSortOrder
	SortDirect SortDirect
}

type PhotoProcessingStatus string

const (
	NewPhoto         PhotoProcessingStatus = "NEW_PHOTO"
	ExifDataSaved    PhotoProcessingStatus = "EXIF_DATA_SAVED"
	MetaDataSaved    PhotoProcessingStatus = "META_DATA_SAVED"
	SystemTagsSaved  PhotoProcessingStatus = "SYSTEM_TAGS_SAVED"
	PhotoVectorSaved PhotoProcessingStatus = "PHOTO_VECTOR_SAVED" // Конечная в данный момент
)

type PhotoExtension string

const (
	PhotoExtensionJpg  PhotoExtension = "JPG"
	PhotoExtensionJpeg PhotoExtension = "JPEG"
	PhotoExtensionPng  PhotoExtension = "PNG"
	PhotoExtensionBmb  PhotoExtension = "BMP"
)

type Photo struct {
	ID               uuid.UUID
	FileName         string
	Hash             string
	UpdateAt         time.Time
	Extension        PhotoExtension
	ProcessingStatus PhotoProcessingStatus
}

type PhotoUploadData struct {
	PhotoID  uuid.UUID
	UploadAt time.Time
	Paths    []string
	ClientId string
}
