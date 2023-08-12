package pgrepo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotoRepository struct {
	Transactor
}

func NewPhotoRepository(pool *pgxpool.Pool) *PhotoRepository {
	return &PhotoRepository{
		Transactor: NewTransactor(pool),
	}
}
