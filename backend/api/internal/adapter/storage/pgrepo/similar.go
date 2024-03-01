package pgrepo

import (
	"context"
	"errors"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"

	"github.com/jackc/pgx/v5"
)

func (r *PhotoRepository) SaveOrUpdatePhotoVector(ctx context.Context, photoVector entity.PhotoVector) error {
	conn := r.getConn(ctx)
	const query = `
		INSERT INTO photo_vectors (photo_id, vector, norm)
		VALUES ($1, $2, $3)
		ON CONFLICT (photo_id) 
			DO UPDATE SET 
				vector = EXCLUDED.vector,
				norm = EXCLUDED.norm;
	`
	_, err := conn.Exec(ctx, query, photoVector.PhotoID, photoVector.Vector, photoVector.Norm)
	if err != nil {
		return printError(err)
	}
	return nil
}

func (r *PhotoRepository) GetPaginatedPhotoVectors(ctx context.Context, offset int64, limit int) ([]entity.PhotoVector, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT photo_id, vector, norm
		FROM photo_vectors
		OFFSET $1
		LIMIT $2
	`

	rows, err := conn.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result []entity.PhotoVector
	for rows.Next() {
		var vector entity.PhotoVector

		errScan := rows.Scan(&vector.PhotoID, &vector.Vector, &vector.Norm)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, errScan
		}

		if errScan != nil {
			return nil, printError(err)
		}

		result = append(result, vector)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}

func (r *PhotoRepository) SaveCoeffSimilarPhoto(ctx context.Context, sim entity.CoeffSimilarPhoto) error {
	conn := r.getConn(ctx)
	const query = `
		INSERT INTO coeffs_similar_photos (photo_id1, photo_id2, coefficient)
		VALUES ($1, $2, $3)
		ON CONFLICT (photo_id1, photo_id2) 
			DO UPDATE SET 
				coefficient = EXCLUDED.coefficient;
	`
	_, err := conn.Exec(ctx, query, sim.PhotoID1, sim.PhotoID2, sim.Coefficient)
	if err != nil {
		return printError(err)
	}
	return nil
}
