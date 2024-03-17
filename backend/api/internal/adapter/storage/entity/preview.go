package entity

import "github.com/google/uuid"

type PhotoPreview struct {
	ID          uuid.UUID
	PhotoID     uuid.UUID
	FileName    string
	WidthPixel  int
	HeightPixel int
	SizePixel   int
}
