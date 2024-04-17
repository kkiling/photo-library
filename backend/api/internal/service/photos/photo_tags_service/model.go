package photo_tags_service

import "github.com/google/uuid"

// TagWithCategory тег фотографии (Информация о теге с информацией о категории)
type TagWithCategory struct {
	// ID тега
	ID uuid.UUID
	// Имя тега
	Name string
	// ID категории
	IDCategory uuid.UUID
	// Тип категории
	Type string
	// Цвет категории
	Color string
}
