package pgrepo

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
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

func (r *PhotoRepository) GetPaginatedPhotoVectors(ctx context.Context, offset int64, limit int64) ([]entity.PhotoVector, error) {
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

func (r *PhotoRepository) FindSimilarPhotoCoefficients(ctx context.Context, photoID uuid.UUID) ([]entity.CoeffSimilarPhoto, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT photo_id1, photo_id2, coefficient
		FROM coeffs_similar_photos
		WHERE photo_id1 = $1 OR photo_id2 = $2;
	`

	rows, err := conn.Query(ctx, query, photoID, photoID)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result []entity.CoeffSimilarPhoto
	for rows.Next() {
		var coefficient entity.CoeffSimilarPhoto

		errScan := rows.Scan(&coefficient.PhotoID1, &coefficient.PhotoID2, &coefficient.Coefficient)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, errScan
		}

		if errScan != nil {
			return nil, printError(err)
		}

		result = append(result, coefficient)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}

func (r *PhotoRepository) GetPhotosVectorCount(ctx context.Context) (int64, error) {
	conn := r.getConn(ctx)

	var counter int64

	builder := sq.
		Select("count(1)").
		From("photo_vectors").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("builder.ToSql: %w", err)
	}

	err = conn.QueryRow(ctx, query, args...).Scan(&counter)
	if err != nil {
		return 0, printError(err)
	}

	return counter, nil
}

func (r *PhotoRepository) GetPhotoVector(ctx context.Context, photoID uuid.UUID) (*entity.PhotoVector, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT photo_id, vector, norm
		FROM photo_vectors
		WHERE photo_id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, photoID)

	var vector entity.PhotoVector
	err := row.Scan(&vector.PhotoID, &vector.Vector, &vector.Norm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &vector, nil
}
