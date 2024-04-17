package model

import "github.com/google/uuid"

type PhotoPreview struct {
	ID          uuid.UUID
	PhotoID     uuid.UUID
	FileKey     string
	WidthPixel  int
	HeightPixel int
	SizePixel   int
	Original    bool
}
