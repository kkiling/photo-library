package photo_groups_service

import (
	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

// GetPhotoGroupsRequest запрос на получение групп фото
type GetPhotoGroupsRequest struct {
	Paginator model.Pagination
}

// PhotoPreview Превью фотографии
type PhotoPreview struct {
	// Путь до фотографии, с нужным get параметром size
	Src string
	// Ширина фото
	Width int
	// Высота фото
	Height int
	// Размер фото (max Width Height)
	Size int
}

// PhotoWithPreviews оригинал фотографии с превью
type PhotoWithPreviews struct {
	// ID фотографии
	ID uuid.UUID
	// Оригинал фотографии
	Original PhotoPreview
	// Превью фотографии (разные размеры)
	Previews []PhotoPreview
}

// PhotoContent контент фото и ее расширение
type PhotoContent struct {
	PhotoBody []byte
	Extension model.PhotoExtension
}

// PhotoGroup Группа фотографий
type PhotoGroup struct {
	// ID группы
	ID uuid.UUID
	// Главная фотография группы
	MainPhoto PhotoWithPreviews
	// Количество фоток объединенных в эту группу
	PhotosCount int
	// Все фото объединенные в группу
	Photos []PhotoWithPreviews
}

// PaginatedPhotoGroups ответ на получение групп фото
type PaginatedPhotoGroups struct {
	Items      []PhotoGroup
	TotalItems int
}
