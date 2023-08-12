package model

import (
	"github.com/google/uuid"
	"time"
)

type Geo struct {
	Latitude  float64
	Longitude float64
}

type MetaData struct {
	PhotoID     uuid.UUID
	SizeBytes   int
	WidthPixel  int
	HeightPixel int
	// Дата и время снимка (берем из exif если нет то пробуем из имени файла
	PhotoAt *time.Time
	// Дата последнего обновления файла
	// Берем из файла
	UpdateAt time.Time
	Geo      *Geo
}
