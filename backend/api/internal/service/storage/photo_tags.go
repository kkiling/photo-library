package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage/photo_library"
)

func (r *Adapter) GetTagCategory(ctx context.Context, id uuid.UUID) (model.TagCategory, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetTagCategory(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.TagCategory{}, serviceerr.ErrNotFound
		}
		return model.TagCategory{}, printError(err)
	}

	return model.TagCategory{
		ID:    res.ID,
		Type:  res.Type,
		Color: res.Color,
	}, nil
}

func (r *Adapter) GetTagCategoryByType(ctx context.Context, typeCategory string) (model.TagCategory, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetTagCategoryByType(ctx, typeCategory)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.TagCategory{}, serviceerr.ErrNotFound
		}
		return model.TagCategory{}, printError(err)
	}

	return model.TagCategory{
		ID:    res.ID,
		Type:  res.Type,
		Color: res.Color,
	}, nil
}

func (r *Adapter) SaveTagCategory(ctx context.Context, category model.TagCategory) error {
	queries := r.getQueries(ctx)

	err := queries.SaveTagCategory(ctx, photo_library.SaveTagCategoryParams{
		ID:    category.ID,
		Type:  category.Type,
		Color: category.Color,
	})

	if err != nil {
		if isAlreadyExist(err) {
			return fmt.Errorf("%w, %w", err, serviceerr.ErrTagAlreadyExist)
		}

		return printError(err)
	}

	return nil
}

func (r *Adapter) SaveTag(ctx context.Context, tag model.Tag) error {
	queries := r.getQueries(ctx)

	err := queries.SaveTag(ctx, photo_library.SaveTagParams{
		ID:         tag.ID,
		CategoryID: tag.CategoryID,
		PhotoID:    tag.PhotoID,
		Name:       tag.Name,
	})

	if err != nil {
		if isAlreadyExist(err) {
			return fmt.Errorf("%w, %w", err, serviceerr.ErrTagAlreadyExist)
		}
		return printError(err)
	}

	return nil
}

func (r *Adapter) GetTags(ctx context.Context, photoID uuid.UUID) ([]model.Tag, error) {
	queries := r.getQueries(ctx)

	res, err := queries.GetTags(ctx, photoID)

	if err != nil {
		return nil, printError(err)
	}

	return lo.Map(res, func(row photo_library.GetTagsRow, index int) model.Tag {
		return model.Tag{
			ID:         row.ID,
			CategoryID: row.CategoryID,
			PhotoID:    photoID,
			Name:       row.Name,
		}
	}), nil
}

func (r *Adapter) DeletePhotoTagsByCategories(ctx context.Context, photoID uuid.UUID, categoryIDs []uuid.UUID) error {
	queries := r.getQueries(ctx)

	_, err := queries.DeletePhotoTagsByCategories(ctx, photo_library.DeletePhotoTagsByCategoriesParams{
		PhotoID:     photoID,
		CategoryIds: categoryIDs,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return serviceerr.ErrNotFound
		}
		return printError(err)
	}

	return nil
}

func (r *Adapter) DeleteTag(ctx context.Context, id uuid.UUID) error {
	queries := r.getQueries(ctx)

	_, err := queries.DeleteTag(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return serviceerr.ErrNotFound
		}
		return printError(err)
	}

	return nil
}
