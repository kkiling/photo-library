package model

import (
	"github.com/google/uuid"
	"time"
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
	PhotoProcessingNew         PhotoProcessingStatus = "NEW_PHOTO"
	PhotoProcessingExifData    PhotoProcessingStatus = "SAVE_EXIF_DATA"
	PhotoProcessingMetaData    PhotoProcessingStatus = "SAVE_META_DATA"
	PhotoProcessingTagsByMeta  PhotoProcessingStatus = "CREATE_TAGS_BY_META"
	PhotoProcessingPhotoVector PhotoProcessingStatus = "SAVE_PHOTO_VECTOR" // Конечная в данный момент
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
	UploadAt         time.Time
	Extension        PhotoExtension
	ProcessingStatus PhotoProcessingStatus
}

type UploadPhotoData struct {
	PhotoID  uuid.UUID
	Paths    []string
	UploadAt time.Time
	ClientId string
}
