package entity

import "github.com/google/uuid"

type PhotoGroupSortOrder string

const (
	PhotoGroupSortOrderNone PhotoSortOrder = "NONE"
)

type PhotoGroupSelectParams struct {
	Offset     uint64
	Limit      uint64
	SortOrder  PhotoSortOrder
	SortDirect SortDirect
}

type PhotoGroup struct {
	ID          uuid.UUID
	MainPhotoID uuid.UUID
	PhotoIDs    []uuid.UUID
}
