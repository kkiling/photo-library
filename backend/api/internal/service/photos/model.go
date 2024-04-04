package photos

import (
	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

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
	// Путь до фотографии
	Src string
	// Ширина фото
	Width int
	// Высота фото
	Height int
	// Превью фотографии (разные размеры)
	Previews []PhotoPreview
}

// PhotoContent контент фото и ее расширение
type PhotoContent struct {
	PhotoBody []byte
	Extension model.PhotoExtension
}

// Tag тег фотографии
type Tag struct {
	ID    uuid.UUID
	Name  string
	Type  string
	Color string
}

// PhotoGroup Группа фотографий
type PhotoGroup struct {
	// ID группы
	ID uuid.UUID
	// Оригинал основной фотографии в группе
	MainPhoto PhotoWithPreviews
	// Количество фоток объединенных в эту группу
	PhotosCount int
}

// PhotoGroupData более детальная информация о конкретной группе
type PhotoGroupData struct {
	PhotoGroup
	// Все фото объединенные в группу
	Photos []PhotoWithPreviews
	// Метаданные mainPhoto
	Metadata *model.PhotoMetadata
	// Теги mainPhoto
	Tags []Tag
}

// GetPhotoGroupsRequest запрос на получение групп фото
type GetPhotoGroupsRequest struct {
	Paginator model.Pagination
}

// PaginatedPhotoGroups ответ на получение групп фото
type PaginatedPhotoGroups struct {
	Items      []PhotoGroup
	TotalItems int
}
