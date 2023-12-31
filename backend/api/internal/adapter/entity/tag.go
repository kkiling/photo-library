package entity

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
