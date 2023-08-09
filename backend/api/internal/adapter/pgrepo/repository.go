package pgrepo

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type PhotoRepository struct {
	Transactor
}

func NewPhotoRepository(conn *pgx.Conn) *PhotoRepository {
	return &PhotoRepository{
		Transactor: NewTransactor(conn),
	}
}

func (r *PhotoRepository) GetPhotoByHash(ctx context.Context, hash string) (*Photo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PhotoRepository) SavePhoto(ctx context.Context, photo Photo) error {
	//TODO implement me
	panic("implement me")
}

func (r *PhotoRepository) SaveUploadPhotoData(ctx context.Context, data UploadPhotoData) error {
	//TODO implement me
	panic("implement me")
}
