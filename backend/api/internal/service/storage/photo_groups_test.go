package storage

import (
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) equalPhotoGroup(new, find model.PhotoGroup) {
	s.Require().Equal(new.ID, find.ID)
	s.Require().Equal(new.MainPhotoID, find.MainPhotoID)
	s.Require().Equal(new.CreatedAt.Unix(), find.CreatedAt.Unix())
	s.Require().Equal(new.UpdatedAt.Unix(), find.UpdatedAt.Unix())
}

func (s *testSuite) equalPhotoGroupWithPhotoIDs(new, find model.PhotoGroupWithPhotoIDs) {
	s.equalPhotoGroup(new.PhotoGroup, find.PhotoGroup)
	s.Assert().ElementsMatch(new.PhotoIDs, find.PhotoIDs)
}

func (s *testSuite) TestAdapter_SaveGroup() {
	s.Run("ok save group", func() {
		photoID := s.savePhoto()

		newGroup := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: photoID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup)
		s.Require().NoError(err)

		findGroup, err := s.storage.GetGroupByID(s.ctx, newGroup.ID)
		s.Require().NoError(err)
		s.equalPhotoGroup(newGroup, findGroup.PhotoGroup)
	})
	s.Run("err save group without photo", func() {
		newGroup := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup)
		s.Require().Error(err)
	})
	s.Run("err save double group", func() {
		groupID := uuid.New()
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()

		newGroup1 := model.PhotoGroup{
			ID:          groupID,
			MainPhotoID: photoID1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup1)
		s.Require().NoError(err)

		newGroup2 := model.PhotoGroup{
			ID:          groupID,
			MainPhotoID: photoID2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = s.storage.SaveGroup(s.ctx, newGroup2)
		s.Require().Error(err)

		findGroup, err := s.storage.GetGroupByID(s.ctx, newGroup1.ID)
		s.Require().NoError(err)
		s.equalPhotoGroup(newGroup1, findGroup.PhotoGroup)
	})
	s.Run("err save double main photo", func() {
		photoID := s.savePhoto()

		newGroup1 := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: photoID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup1)
		s.Require().NoError(err)

		newGroup2 := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: photoID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = s.storage.SaveGroup(s.ctx, newGroup2)
		s.Require().Error(err)

		findGroup, err := s.storage.GetGroupByID(s.ctx, newGroup1.ID)
		s.Require().NoError(err)
		s.equalPhotoGroup(newGroup1, findGroup.PhotoGroup)
	})
}

func (s *testSuite) TestAdapter_FindGroupIDByPhotoID() {
	s.Run("ok find photo id", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()
		photoID3 := s.savePhoto()

		newGroup := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: photoID1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup)
		s.Require().NoError(err)

		err = s.storage.AddPhotoIDsToGroup(s.ctx, newGroup.ID, []uuid.UUID{photoID1, photoID2, photoID3})
		s.Require().NoError(err)

		id1, err := s.storage.FindGroupIDByPhotoID(s.ctx, photoID1)
		s.Require().NoError(err)
		s.Require().Equal(newGroup.ID, id1)

		id2, err := s.storage.FindGroupIDByPhotoID(s.ctx, photoID2)
		s.Require().NoError(err)
		s.Require().Equal(newGroup.ID, id2)

		id3, err := s.storage.FindGroupIDByPhotoID(s.ctx, photoID3)
		s.Require().NoError(err)
		s.Require().Equal(newGroup.ID, id3)
	})
	s.Run("not found", func() {
		_, err := s.storage.FindGroupIDByPhotoID(s.ctx, uuid.New())
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_GetGroupByID() {
	s.Run("ok get group", func() {
		photoID := s.savePhoto()

		newGroup := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: photoID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup)
		s.Require().NoError(err)

		findGroup, err := s.storage.GetGroupByID(s.ctx, newGroup.ID)
		s.Require().NoError(err)
		s.equalPhotoGroup(newGroup, findGroup.PhotoGroup)
	})
	s.Run("ok no found group", func() {
		_, err := s.storage.GetGroupByID(s.ctx, uuid.New())
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_GetGroupPhotoIDs() {
	s.Run("ok find photo ids", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()
		photoID3 := s.savePhoto()

		newGroup := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: photoID1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup)
		s.Require().NoError(err)

		err = s.storage.AddPhotoIDsToGroup(s.ctx, newGroup.ID, []uuid.UUID{photoID1, photoID2, photoID3})
		s.Require().NoError(err)

		ids, err := s.storage.GetGroupPhotoIDs(s.ctx, newGroup.ID)
		s.Require().NoError(err)
		s.Require().Equal(ids, []uuid.UUID{photoID1, photoID2, photoID3})
	})
	s.Run("not found", func() {
		ids, err := s.storage.GetGroupPhotoIDs(s.ctx, uuid.New())
		s.Require().NoError(err)
		s.Require().Len(ids, 0)
	})
}

func (s *testSuite) TestAdapter_DeletePhotoGroup() {
	s.Run("ok", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()
		photoID3 := s.savePhoto()

		newGroup := model.PhotoGroup{
			ID:          uuid.New(),
			MainPhotoID: photoID1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.storage.SaveGroup(s.ctx, newGroup)
		s.Require().NoError(err)

		err = s.storage.AddPhotoIDsToGroup(s.ctx, newGroup.ID, []uuid.UUID{photoID2, photoID3})
		s.Require().NoError(err)

		err = s.storage.DeletePhotoGroup(s.ctx, newGroup.ID)
		s.Require().NoError(err)

		ids, err := s.storage.GetGroupPhotoIDs(s.ctx, newGroup.ID)
		s.Require().NoError(err)
		s.Require().Len(ids, 0)

		_, err = s.storage.GetGroupByID(s.ctx, newGroup.ID)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
	s.Run("not found", func() {
		err := s.storage.DeletePhotoGroup(s.ctx, uuid.New())
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_GetPhotoGroupsCount() {
	s.Run("ok get count", func() {
		count := 10
		for i := 0; i < count; i++ {
			photoID := s.savePhoto()
			err := s.storage.SaveGroup(s.ctx, model.PhotoGroup{
				ID:          uuid.New(),
				MainPhotoID: photoID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
			s.Require().NoError(err)
		}
		totalCount, err := s.storage.GetPhotoGroupsCount(s.ctx, model.PhotoGroupsFilter{
			// TODO: Будет дорабатываться
		})
		s.Require().NoError(err)

		// TODO: доработать
		s.Require().Less(int64(count), totalCount)
	})
}

func (s *testSuite) TestAdapter_GetPaginatedPhotoGroups() {
	s.Run("ok get count", func() {
		count := 10
		for i := 0; i < count; i++ {
			photoID := s.savePhoto()
			group := model.PhotoGroup{
				ID:          uuid.New(),
				MainPhotoID: photoID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			err := s.storage.SaveGroup(s.ctx, group)
			s.Require().NoError(err)

			err = s.storage.AddPhotoIDsToGroup(s.ctx, group.ID, []uuid.UUID{group.MainPhotoID})
			s.Require().NoError(err)
		}

		items, err := s.storage.GetPaginatedPhotoGroups(s.ctx, model.PhotoGroupsParams{
			Paginator: model.Pagination{
				Page:    0,
				PerPage: 5,
			},
			//SortOrder:  "",
			// SortDirect: "",
			Filter: model.PhotoGroupsFilter{
				// TODO: Будет дорабатываться
			},
		})
		s.Require().NoError(err)
		s.Require().Len(items, 5)
		// TODO: доработать
	})
}
