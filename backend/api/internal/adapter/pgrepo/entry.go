package pgrepo

import (
	"github.com/google/uuid"
	"time"
)

type Photo struct {
	ID        uuid.UUID `db:"id"`
	FilePath  string    `db:"file_path"`
	Hash      string    `db:"hash"`
	UpdateAt  time.Time `db:"update_at"`
	UploadAt  time.Time `db:"upload_at"`
	Extension string    `db:"extension"`
}

type UploadPhotoData struct {
	ID       uuid.UUID `db:"id"`
	PhotoID  uuid.UUID `db:"photo_id"`
	Paths    []string  `db:"paths"`
	UploadAt time.Time `db:"upload_at"`
	ClientId string    `db:"client_id"`
}
