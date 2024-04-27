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

// PhotoGroupDTO Группа фотографий
type PhotoGroupDTO struct {
	// GroupID группы
	GroupID uuid.UUID
	// Главная фотография группы
	MainPhoto PhotoWithPreviewsDTO
	// Количество фоток объединенных в эту группу
	PhotosCount int
	// Все фото объединенные в группу
	Photos []PhotoWithPreviewsDTO
}

// PaginatedPhotoGroupsDTO ответ на получение групп фото
type PaginatedPhotoGroupsDTO struct {
	Items      []PhotoGroupDTO
	TotalItems int
}
