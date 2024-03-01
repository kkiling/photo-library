package entity

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

type PhotoStatus string

const (
	NewPhotoStatus PhotoStatus = "NEW_PHOTO"
	NotValidStatus PhotoStatus = "NOT_VALID"
)

type PhotoSelectParams struct {
	Offset     int64
	Limit      int
	SortOrder  PhotoSortOrder
	SortDirect SortDirect
}

type Photo struct {
	ID        uuid.UUID
	FileName  string
	Hash      string
	UpdateAt  time.Time
	Extension string
}

type PhotoUploadData struct {
	PhotoID  uuid.UUID
	UploadAt time.Time
	Paths    []string
	ClientId string
}
