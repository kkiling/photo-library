package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/geo"
)

type Location struct {
	PhotoID   uuid.UUID
	CreatedAt time.Time
	Latitude  float64
	Longitude float64
	Geo       geo.Address
}
