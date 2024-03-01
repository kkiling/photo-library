package model

import (
	"time"

	"github.com/google/uuid"
)

type PhotoFilter struct {
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
	ExifDataProcessing    PhotoProcessingStatus = "EXIF_DATA"
	MetaDataProcessing    PhotoProcessingStatus = "META_DATA"
	CatalogTagsProcessing PhotoProcessingStatus = "CATALOG_TAGS"
	MetaTagsProcessing    PhotoProcessingStatus = "META_TAGS"
	PhotoVectorProcessing PhotoProcessingStatus = "PHOTO_VECTOR"
)

const LastProcessingStatus = PhotoVectorProcessing

var PhotoProcessingStatuses = []PhotoProcessingStatus{
	ExifDataProcessing,
	MetaDataProcessing,
	CatalogTagsProcessing,
	MetaTagsProcessing,
	PhotoVectorProcessing,
}

type PhotoExtension string

const (
	PhotoExtensionJpg  PhotoExtension = "JPG"
	PhotoExtensionJpeg PhotoExtension = "JPEG"
	PhotoExtensionPng  PhotoExtension = "PNG"
	PhotoExtensionBmb  PhotoExtension = "BMP"
)

type Photo struct {
	ID        uuid.UUID
	FileName  string
	Hash      string
	UpdateAt  time.Time
	Extension PhotoExtension
}

type PhotoUploadData struct {
	PhotoID  uuid.UUID
	UploadAt time.Time
	Paths    []string
	ClientId string
}
