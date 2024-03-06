package model

import "github.com/google/uuid"

type PhotoVector struct {
	PhotoID uuid.UUID
	Vector  []float64
	Norm    float64
}

type CoeffSimilarPhoto struct {
	PhotoID1    uuid.UUID
	PhotoID2    uuid.UUID
	Coefficient float64
}
