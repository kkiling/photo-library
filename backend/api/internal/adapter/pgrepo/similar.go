package pgrepo

import (
	"context"
	"github.com/google/uuid"
)

func (r *PhotoRepository) SaveOrUpdatePhotoVector(ctx context.Context, photoID uuid.UUID, vector []float64, norm float64) error {
	conn := r.getConn(ctx)
	const query = `
		INSERT INTO photo_vector (photo_id, vector, norm)
		VALUES ($1, $2, $3)
		ON CONFLICT (photo_id) 
			DO UPDATE SET 
				vector = EXCLUDED.vector,
				norm = EXCLUDED.norm;
	`
	_, err := conn.Exec(ctx, query, photoID, vector, norm)
	if err != nil {
		return printError(err)
	}
	return nil
}

func (r *PhotoRepository) ExistPhotoVector(ctx context.Context, photoID uuid.UUID) (bool, error) {
	conn := r.getConn(ctx)

	var counter int64

	const query = `
		SELECT count(*)
		FROM photo_vector
		WHERE photo_id = $1
		LIMIT 1
	`

	err := conn.QueryRow(ctx, query, photoID).Scan(&counter)
	if err != nil {
		return false, printError(err)
	}

	return counter > 0, nil
}
