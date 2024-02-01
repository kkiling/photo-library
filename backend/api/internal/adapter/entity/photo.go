package entity

import (
	"time"

	"github.com/google/uuid"
)

type PhotoFilter struct {
	ProcessingStatusIn []string
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

type Photo struct {
	ID               uuid.UUID
	FileName         string
	Hash             string
	UpdateAt         time.Time
	UploadAt         time.Time
	Extension        string
	ProcessingStatus string
}

type UploadPhotoData struct {
	PhotoID  uuid.UUID
	Paths    []string
	UploadAt time.Time
	ClientId string
}
