package pgrepo

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
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

func (r *PhotoRepository) AddPhotoIDsToGroup(ctx context.Context, groupID uuid.UUID, photoIDs []uuid.UUID) error {
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

func (r *PhotoRepository) GetPhotoGroupsCount(ctx context.Context) (uint64, error) {
	conn := r.getConn(ctx)

	var counter uint64

	builder := sq.
		Select("count(1)").
		From("photo_groups").
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

func (r *PhotoRepository) getGroupPhotoIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	conn := r.getConn(ctx)

	const query = `
		SELECT photo_id
		FROM photo_groups_photos
		WHERE group_id = $1
	`

	rows, err := conn.Query(ctx, query, groupID)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result []uuid.UUID
	for rows.Next() {
		var id uuid.UUID

		errScan := rows.Scan(&id)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, errScan
		}

		if errScan != nil {
			return nil, printError(err)
		}

		result = append(result, id)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}

func (r *PhotoRepository) GetPaginatedPhotoGroups(ctx context.Context, params entity.PhotoSelectParams) ([]entity.PhotoGroup, error) {
	conn := r.getConn(ctx)

	builder := sq.
		Select("id", "main_photo_id").
		From("photo_groups").
		Offset(params.Offset).
		Limit(params.Limit).
		PlaceholderFormat(sq.Dollar)

	if params.SortOrder != entity.PhotoSortOrderNone {
		// Добавляем сортировку
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("builder.ToSql: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, printError(err)
	}
	defer rows.Close()

	var result = make([]entity.PhotoGroup, 0, params.Limit)
	for rows.Next() {
		var group entity.PhotoGroup

		errScan := rows.Scan(&group.ID, &group.MainPhotoID)
		if errScan != nil {
			if errors.Is(errScan, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, printError(err)
		}

		photoIDs, err := r.getGroupPhotoIDs(ctx, group.ID)
		if err != nil {
			return nil, fmt.Errorf("r.getGroupPhotoIDs")
		}

		group.PhotoIDs = photoIDs

		result = append(result, group)
	}

	if err := rows.Err(); err != nil {
		return nil, printError(err)
	}

	return result, nil
}
