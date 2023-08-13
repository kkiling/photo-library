package entity

import (
	"github.com/google/uuid"
	"time"
)

type MetaData struct {
	PhotoID      uuid.UUID
	ModelInfo    *string
	SizeBytes    int
	WidthPixel   int
	HeightPixel  int
	DateTime     *time.Time
	UpdateAt     time.Time
	GeoLatitude  *float64
	GeoLongitude *float64
}
