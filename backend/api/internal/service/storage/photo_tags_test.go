package storage

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) equalTagCategory(new, find model.TagCategory) {
	s.Require().Equal(new.ID, find.ID)
	s.Require().Equal(new.Type, find.Type)
	s.Require().Equal(new.Color, find.Color)
}

func (s *testSuite) equalTag(new, find model.Tag) {
	s.Require().Equal(new.ID, find.ID)
	s.Require().Equal(new.CategoryID, find.CategoryID)
	s.Require().Equal(new.PhotoID, find.PhotoID)
	s.Require().Equal(new.Name, find.Name)
}

func (s *testSuite) saveTagCategory() model.TagCategory {
	newCategory := model.TagCategory{
		ID:    uuid.New(),
		Type:  fmt.Sprintf("tag_type_%s", uuid.NewString()),
		Color: "#FF00FF",
	}
	err := s.storage.SaveTagCategory(s.ctx, newCategory)
	s.Require().NoError(err)
	return newCategory
}

func (s *testSuite) TestAdapter_SaveTagCategory() {
	s.Run("ok save tag category", func() {
		newCategory := s.saveTagCategory()

		findCategory, err := s.storage.GetTagCategory(s.ctx, newCategory.ID)
		s.Require().NoError(err)
		s.equalTagCategory(newCategory, findCategory)
	})

	s.Run("error save duplicate tag category", func() {
		typeCategory := fmt.Sprintf("tag_type_%s", uuid.NewString())
		newCategory1 := model.TagCategory{
			ID:    uuid.New(),
			Type:  typeCategory,
			Color: "#FF00FF",
		}
		err := s.storage.SaveTagCategory(s.ctx, newCategory1)
		s.Require().NoError(err)

		newCategory2 := model.TagCategory{
			ID:    uuid.New(),
			Type:  typeCategory,
			Color: "#FF00FF",
		}
		err = s.storage.SaveTagCategory(s.ctx, newCategory2)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrTagAlreadyExist)
	})
}

func (s *testSuite) TestAdapter_GetTagCategory() {
	s.Run("ok get tag category", func() {
		newCategory := model.TagCategory{
			ID:    uuid.New(),
			Type:  fmt.Sprintf("tag_type_%s", uuid.NewString()),
			Color: "#FF00FF",
		}
		err := s.storage.SaveTagCategory(s.ctx, newCategory)
		s.Require().NoError(err)

		findCategory, err := s.storage.GetTagCategory(s.ctx, newCategory.ID)
		s.Require().NoError(err)
		s.equalTagCategory(newCategory, findCategory)
	})

	s.Run("get tag category not found", func() {
		id := uuid.New()
		_, err := s.storage.GetTagCategory(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_GetTagCategoryByType() {
	s.Run("ok get tag category", func() {
		newCategory := model.TagCategory{
			ID:    uuid.New(),
			Type:  fmt.Sprintf("tag_type_%s", uuid.NewString()),
			Color: "#FF00FF",
		}
		err := s.storage.SaveTagCategory(s.ctx, newCategory)
		s.Require().NoError(err)

		findCategory, err := s.storage.GetTagCategoryByType(s.ctx, newCategory.Type)
		s.Require().NoError(err)
		s.equalTagCategory(newCategory, findCategory)
	})

	s.Run("get tag category not found", func() {
		_, err := s.storage.GetTagCategoryByType(s.ctx, uuid.NewString())
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_SaveTag() {
	s.Run("ok save tag", func() {
		photoID := s.savePhoto()
		newCategory := s.saveTagCategory()

		err := s.storage.SaveTag(s.ctx, model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategory.ID,
			PhotoID:    photoID,
			Name:       fmt.Sprintf("tag_name_%s", uuid.NewString()),
		})
		s.Require().NoError(err)
	})

	s.Run("save double tag", func() {
		photoID := s.savePhoto()
		newCategory := s.saveTagCategory()
		tagName := fmt.Sprintf("tag_name_%s", uuid.NewString())

		err := s.storage.SaveTag(s.ctx, model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategory.ID,
			PhotoID:    photoID,
			Name:       tagName,
		})
		s.Require().NoError(err)

		err = s.storage.SaveTag(s.ctx, model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategory.ID,
			PhotoID:    photoID,
			Name:       tagName,
		})
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrTagAlreadyExist)
	})

	s.Run("error save tag unknown photo", func() {
		photoID := uuid.New()
		newCategory := s.saveTagCategory()

		err := s.storage.SaveTag(s.ctx, model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategory.ID,
			PhotoID:    photoID,
			Name:       fmt.Sprintf("tag_name_%s", uuid.NewString()),
		})
		s.Require().Error(err)
	})

	s.Run("error save tag unknown category", func() {
		photoID := s.savePhoto()
		newCategoryID := uuid.New()

		err := s.storage.SaveTag(s.ctx, model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategoryID,
			PhotoID:    photoID,
			Name:       fmt.Sprintf("tag_name_%s", uuid.NewString()),
		})
		s.Require().Error(err)
	})
}

func (s *testSuite) TestAdapter_GetTags() {
	s.Run("ok get tags", func() {
		photoID := s.savePhoto()
		newCategory := s.saveTagCategory()
		tag1 := model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategory.ID,
			PhotoID:    photoID,
			Name:       fmt.Sprintf("tag_name_%s", uuid.NewString()),
		}
		tag2 := model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategory.ID,
			PhotoID:    photoID,
			Name:       fmt.Sprintf("tag_name_%s", uuid.NewString()),
		}

		err := s.storage.SaveTag(s.ctx, tag1)
		s.Require().NoError(err)

		err = s.storage.SaveTag(s.ctx, tag2)
		s.Require().NoError(err)

		tags, err := s.storage.GetTags(s.ctx, photoID)
		s.Require().NoError(err)
		s.Require().Equal(len(tags), 2)
		s.equalTag(tag1, tags[0])
		s.equalTag(tag2, tags[1])
	})
}

func (s *testSuite) TestAdapter_DeleteTags() {
	s.Run("ok delete tags", func() {
		photoID := s.savePhoto()
		newCategory := s.saveTagCategory()

		err := s.storage.SaveTag(s.ctx, model.Tag{
			ID:         uuid.New(),
			CategoryID: newCategory.ID,
			PhotoID:    photoID,
			Name:       fmt.Sprintf("tag_name_%s", uuid.NewString()),
		})
		s.Require().NoError(err)

		err = s.storage.DeleteTags(s.ctx, photoID)
		s.Require().NoError(err)

		res, err := s.storage.GetTags(s.ctx, photoID)
		s.Require().NoError(err)
		s.Require().Len(res, 0)
	})

	s.Run("delete tags not found", func() {
		id := uuid.New()
		err := s.storage.DeleteTags(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}
