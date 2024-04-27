package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) FindGroupIDByPhotoID(ctx context.Context, photoID uuid.UUID) (uuid.UUID, error) {
	queries := r.getQueries(ctx)
	res, err := queries.FindGroupIDByPhotoID(ctx, photoID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, serviceerr.ErrNotFound
		}
		return uuid.UUID{}, printError(err)
	}

	return res, nil
}

func (r *Adapter) GetGroupByID(ctx context.Context, groupID uuid.UUID) (model.PhotoGroupWithPhotoIDs, error) {
	queries := r.getQueries(ctx)
	res, err := queries.GetGroupByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PhotoGroupWithPhotoIDs{}, serviceerr.ErrNotFound
		}
		return model.PhotoGroupWithPhotoIDs{}, printError(err)
	}

	if len(res) == 0 {
		return model.PhotoGroupWithPhotoIDs{}, serviceerr.ErrNotFound
	}

	group := res[0]
	return model.PhotoGroupWithPhotoIDs{
		PhotoGroup: model.PhotoGroup{
			ID:          group.ID,
			MainPhotoID: group.MainPhotoID,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		},
		PhotoIDs: lo.FilterMap(res, func(x photo_library.GetGroupByIDRow, _ int) (uuid.UUID, bool) {
			return x.PhotoID.Bytes, x.PhotoID.Valid
		}),
	}, nil
}

func (r *Adapter) SaveGroup(ctx context.Context, group model.PhotoGroup) error {
	queries := r.getQueries(ctx)
	paramsGroup := photo_library.SaveGroupParams{
		ID:          group.ID,
		MainPhotoID: group.MainPhotoID,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
	if err := queries.SaveGroup(ctx, paramsGroup); err != nil {
		return printError(err)
	}

	return nil

	/*err := r.runTransaction(ctx, func(ctxTx context.Context) error {
		paramsGroup := photo_library.SaveGroupParams{
			ID:          group.ID,
			MainPhotoID: group.MainPhotoID,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		if err := queries.SaveGroup(ctxTx, paramsGroup); err != nil {
			return err
		}

		paramsAdd := photo_library.AddPhotoIDToGroupParams{
			PhotoID: group.MainPhotoID,
			ID: group.ID,
		}
		if err := queries.AddPhotoIDToGroup(ctxTx, paramsAdd); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return printError(err)
	}

	return nil*/
}

func (r *Adapter) AddPhotoIDsToGroup(ctx context.Context, groupID uuid.UUID, photoIDs []uuid.UUID) error {
	queries := r.getQueries(ctx)

	err := r.runTransaction(ctx, func(ctxTx context.Context) error {
		for _, photoID := range photoIDs {
			params := photo_library.AddPhotoIDToGroupParams{
				PhotoID: photoID,
				GroupID: groupID,
			}

			err := queries.AddPhotoIDToGroup(ctx, params)

			if err != nil {
				return printError(err)
			}
		}
		return nil
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetPhotoGroupsCount(ctx context.Context, _ model.PhotoGroupsFilter) (int64, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetPhotoGroupsCount(ctx)
	if err != nil {
		return 0, printError(err)
	}

	return res, nil
}

func (r *Adapter) GetPaginatedPhotoGroups(ctx context.Context, params model.PhotoGroupsParams) ([]model.PhotoGroupWithPhotoIDs, error) {
	queries := r.getQueries(ctx)
	args := photo_library.GetPaginatedPhotoGroupsParams{
		Offset: params.Paginator.GetOffset(),
		Limit:  params.Paginator.GetLimit(),
	}
	res, err := queries.GetPaginatedPhotoGroups(ctx, args)
	if err != nil {
		return nil, printError(err)
	}

	uniqRes := lo.UniqBy(res, func(item photo_library.GetPaginatedPhotoGroupsRow) uuid.UUID {
		return item.ID
	})

	result := make([]model.PhotoGroupWithPhotoIDs, 0, len(uniqRes))
	for _, item := range uniqRes {
		result = append(result, model.PhotoGroupWithPhotoIDs{
			PhotoGroup: model.PhotoGroup{
				ID:          item.ID,
				MainPhotoID: item.MainPhotoID,
				UpdatedAt:   item.UpdatedAt,
			},
			PhotoIDs: lo.FilterMap(res, func(x photo_library.GetPaginatedPhotoGroupsRow, _ int) (uuid.UUID, bool) {
				return x.PhotoID.Bytes, x.PhotoID.Valid && x.ID == item.ID
			}),
		})
	}

	return result, nil
}

func (r *Adapter) GetGroupPhotoIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	queries := r.getQueries(ctx)
	res, err := queries.GetGroupPhotoIDs(ctx, groupID)
	if err != nil {
		return nil, printError(err)
	}
	return res, nil
}

func (r *Adapter) DeletePhotoGroup(ctx context.Context, groupID uuid.UUID) error {
	queries := r.getQueries(ctx)
	err := r.runTransaction(ctx, func(ctxTx context.Context) error {
		err := queries.DeletePhotoGroupPhotos(ctxTx, groupID)
		if err != nil {
			return err
		}

		_, err = queries.DeletePhotoGroup(ctxTx, groupID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return serviceerr.ErrNotFound
			}
			return err
		}

		return nil
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) DeletePhotoGroupByPhotoID(ctx context.Context, photoID uuid.UUID) error {
	queries := r.getQueries(ctx)
	err := r.runTransaction(ctx, func(ctxTx context.Context) error {
		err := queries.DeletePhotoGroupByMainPhoto(ctxTx, photoID)
		if err != nil {
			return err
		}
		err = queries.DeletePhotoGroupPhotosByPhoto(ctxTx, photoID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *Adapter) SetPhotoGroupMainPhoto(ctx context.Context, groupID, photoID uuid.UUID) error {
	queries := r.getQueries(ctx)
	err := queries.SetPhotoGroupMainPhoto(ctx, photo_library.SetPhotoGroupMainPhotoParams{
		ID:          groupID,
		MainPhotoID: photoID,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return printError(err)
	}
	return nil
}
