package entity

import "github.com/google/uuid"

type ExifPhotoData struct {
	PhotoID uuid.UUID
	Data    []byte
}
