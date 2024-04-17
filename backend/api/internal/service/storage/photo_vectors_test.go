package storage

import (
	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) equalPhotoVector(new, find model.PhotoVector) {
	s.Require().Equal(new.PhotoID, find.PhotoID)
	s.Require().Equal(new.Vector, find.Vector)
	s.Require().Equal(new.Norm, find.Norm)
}

func (s *testSuite) TestAdapter_SavePhotoVector() {
	s.Run("ok save photo vector", func() {
		photoID := s.savePhoto()

		newPhotoVector := model.PhotoVector{
			PhotoID: photoID,
			Vector:  []float64{3246.1235, 4535123.561},
			Norm:    12.023345678,
		}

		err := s.storage.SavePhotoVector(s.ctx, newPhotoVector)
		s.Require().NoError(err)

		findPhotoVector, err := s.storage.GetPhotoVector(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalPhotoVector(newPhotoVector, findPhotoVector)
	})

	s.Run("error save duplicate photo vector", func() {
		photoID := s.savePhoto()

		newPhotoVector1 := model.PhotoVector{
			PhotoID: photoID,
			Vector:  []float64{3246.1235, 4535123.561},
			Norm:    12.023345678,
		}
		newPhotoVector2 := model.PhotoVector{
			PhotoID: photoID,
			Vector:  []float64{3246.1235, 4535123.561},
			Norm:    12.023345678,
		}

		err := s.storage.SavePhotoVector(s.ctx, newPhotoVector1)
		s.Require().NoError(err)
		err = s.storage.SavePhotoVector(s.ctx, newPhotoVector2)
		s.Require().Error(err)

		findExifData, err := s.storage.GetPhotoVector(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalPhotoVector(newPhotoVector1, findExifData)
	})
}

func (s *testSuite) TestAdapter_GetPhotoVector() {
	s.Run("ok get photo vector", func() {
		photoID := s.savePhoto()

		newPhotoVector := model.PhotoVector{
			PhotoID: photoID,
			Vector:  []float64{3246.1235, 4535123.561},
			Norm:    12.023345678,
		}

		err := s.storage.SavePhotoVector(s.ctx, newPhotoVector)
		s.Require().NoError(err)

		findPhotoVector, err := s.storage.GetPhotoVector(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalPhotoVector(newPhotoVector, findPhotoVector)
	})

	s.Run("get photo vector not found", func() {
		id := uuid.New()
		_, err := s.storage.GetPhotoVector(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_DeletePhotoVector() {
	s.Run("ok delete exif", func() {
		photoID := s.savePhoto()

		newPhotoVector := model.PhotoVector{
			PhotoID: photoID,
			Vector:  []float64{3246.1235, 4535123.561},
			Norm:    12.023345678,
		}

		err := s.storage.SavePhotoVector(s.ctx, newPhotoVector)
		s.Require().NoError(err)

		err = s.storage.DeletePhotoVector(s.ctx, photoID)
		s.Require().NoError(err)

		_, err = s.storage.GetPhotoVector(s.ctx, photoID)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})

	s.Run("delete exif not found", func() {
		id := uuid.New()
		err := s.storage.DeletePhotoVector(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_GetPhotoVectors() {
	photoID1 := s.savePhoto()
	newPhotoVector1 := model.PhotoVector{
		PhotoID: photoID1,
		Vector:  []float64{3246.1235, 4535123.561},
		Norm:    12.023345678,
	}

	err := s.storage.SavePhotoVector(s.ctx, newPhotoVector1)
	s.Require().NoError(err)

	photoID2 := s.savePhoto()
	newPhotoVector2 := model.PhotoVector{
		PhotoID: photoID2,
		Vector:  []float64{3246.1235, 4535123.561},
		Norm:    12.023345678,
	}

	err = s.storage.SavePhotoVector(s.ctx, newPhotoVector2)
	s.Require().NoError(err)

	photoID3 := s.savePhoto()
	newPhotoVector3 := model.PhotoVector{
		PhotoID: photoID3,
		Vector:  []float64{3246.1235, 4535123.561},
		Norm:    12.023345678,
	}

	err = s.storage.SavePhotoVector(s.ctx, newPhotoVector3)
	s.Require().NoError(err)

	findPhotoVectors, err := s.storage.GetPhotoVectors(s.ctx, model.Pagination{
		Page:    0,
		PerPage: 10,
	})
	s.Require().NoError(err)
	s.Require().True(len(findPhotoVectors) >= 3)
}
