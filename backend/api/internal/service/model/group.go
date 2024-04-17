package model

import (
	"time"

	"github.com/google/uuid"
)

type PhotoGroupsFilter struct {
}

type PhotoGroupsParams struct {
	Paginator  Pagination
	SortOrder  SortOrder
	SortDirect SortDirect
	Filter     PhotoGroupsFilter
}

type PhotoGroup struct {
	ID          uuid.UUID
	MainPhotoID uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PhotoGroupWithPhotoIDs struct {
	PhotoGroup
	PhotoIDs []uuid.UUID
}
