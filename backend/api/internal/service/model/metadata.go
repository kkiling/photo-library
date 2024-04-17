package model

import (
	"time"

	"github.com/google/uuid"
)

type Geo struct {
	Latitude  float64
	Longitude float64
}

type PhotoMetadata struct {
	PhotoID     uuid.UUID
	ModelInfo   *string
	SizeBytes   int
	WidthPixel  int
	HeightPixel int
	// Дата и время снимка берем из exif, если нет то пробуем из имени файла
	DateTime *time.Time
	// Дата последнего обновления файла (берем из файла)
	UpdatedAt time.Time
	// Geo локация
	Geo *Geo
}
