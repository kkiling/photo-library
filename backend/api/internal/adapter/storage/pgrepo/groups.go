package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
	"time"
)

func (r *PhotoRepository) FindGroupIDByPhotoID(ctx context.Context, photoID uuid.UUID) (*uuid.UUID, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT group_id
		FROM photo_groups_photos AS pg
		WHERE photo_id = $1
		LIMIT 1
	`

	row := conn.QueryRow(ctx, query, photoID)

	var groupID uuid.UUID
	err := row.Scan(&groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &groupID, nil
}

func (r *PhotoRepository) FindGroupByPhotoID(ctx context.Context, photoID uuid.UUID) (*entity.PhotoGroup, error) {

	groupID, err := r.FindGroupIDByPhotoID(ctx, photoID)
	if err != nil {
		return nil, fmt.Errorf("r.findGroupIDByPhotoID: %w", err)
	}
	if groupID == nil {
		return nil, nil
	}

	conn := r.getConn(ctx)
	const query = `
		SELECT g.id, g.main_photo_id, gp.photo_id
		FROM photo_groups AS g
		JOIN photo_groups_photos gp ON g.id = gp.group_id
		WHERE g.id = $1
	`

	rows, err := conn.Query(ctx, query, photoID)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result entity.PhotoGroup
	for rows.Next() {
		var (
			id           uuid.UUID
			mainPhotoID  uuid.UUID
			groupPhotoID uuid.UUID
		)

		errScan := rows.Scan(&id, &mainPhotoID, &groupPhotoID)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, printError(err)
		}

		result.ID = id
		result.MainPhotoID = mainPhotoID
		result.PhotoIDs = append(result.PhotoIDs, groupPhotoID)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return &result, nil
}

func (r *PhotoRepository) CreateGroup(ctx context.Context, mainPhotoID uuid.UUID) (*entity.PhotoGroup, error) {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO photo_groups (id, main_photo_id, update_at)
		VALUES ($1, $2, $3)
	`

	group := entity.PhotoGroup{
		ID:          uuid.New(),
		MainPhotoID: mainPhotoID,
	}

	_, err := conn.Exec(ctx, query, group.ID, group.MainPhotoID, time.Now())
	if err != nil {
		return nil, printError(err)
	}

	return &group, nil
}

func (r *PhotoRepository) AddPhotoToGroup(ctx context.Context, groupID uuid.UUID, photoIDs []uuid.UUID) error {
	conn := r.getConn(ctx)

	const query = `
		INSERT INTO photo_groups_photos (photo_id, group_id)
		VALUES ($1, $2)
	`

	for _, photoID := range photoIDs {
		_, err := conn.Exec(ctx, query, photoID, groupID)
		if err != nil {
			return printError(err)
		}
	}

	return nil
}
