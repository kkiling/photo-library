package model

import "github.com/google/uuid"

type PhotoGroup struct {
	ID          uuid.UUID
	MainPhotoID uuid.UUID
	PhotoIDs    []uuid.UUID
}
