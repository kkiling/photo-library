package model

import "github.com/google/uuid"

type TagCategory struct {
	ID    uuid.UUID
	Type  string
	Color string
}

type Tag struct {
	ID         uuid.UUID
	CategoryID uuid.UUID
	PhotoID    uuid.UUID
	Name       string
}

// TagWithCategoryDTO тег фотографии (Информация о теге с информацией о категории)
type TagWithCategoryDTO struct {
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
