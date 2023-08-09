package pgrepo

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type Repository struct {
	Transactor
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		Transactor: NewTransactor(conn),
	}
}

func (r Repository) GetPhotoByHash(ctx context.Context, hash string) (*model.Photo, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) SavePhoto(ctx context.Context, photo model.Photo) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) SaveUploadPhotoData(ctx context.Context, data model.UploadPhotoData) error {
	//TODO implement me
	panic("implement me")
}
