package storage

import (
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) equalUploadData(new, find model.PhotoUploadData) {
	s.Require().Equal(new.PhotoID, find.PhotoID)
	s.compareTime(&new.UploadAt, &find.UploadAt)
	s.Require().Equal(new.ClientID, find.ClientID)
	s.Require().Equal(new.Paths, find.Paths)
}

func (s *testSuite) createUploadData(photoID uuid.UUID) model.PhotoUploadData {
	newUploadData := model.PhotoUploadData{
		PhotoID:  photoID,
		UploadAt: time.Now(),
		Paths: []string{
			"photo_library/photo1.jpeg",
			"photo_library/photo2.jpeg",
		},
		ClientID: uuid.NewString(),
	}

	return newUploadData
}

func (s *testSuite) saveUploadData(photoID uuid.UUID) model.PhotoUploadData {
	newUploadData := s.createUploadData(photoID)

	err := s.storage.SavePhotoUploadData(s.ctx, newUploadData)
	s.Require().NoError(err)

	return newUploadData
}

func (s *testSuite) TestAdapter_SavePhotoUploadData() {
	s.Run("ok save upload data", func() {
		photoID := s.savePhoto()

		newUploadData := s.saveUploadData(photoID)

		findUploadData, err := s.storage.GetPhotoUploadData(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalUploadData(newUploadData, findUploadData)
	})

	s.Run("error save duplicate upload data", func() {
		photoID := s.savePhoto()

		newUploadData1 := s.saveUploadData(photoID)
		newUploadData2 := s.createUploadData(photoID)
		err := s.storage.SavePhotoUploadData(s.ctx, newUploadData2)
		s.Require().Error(err)

		findUploadData, err := s.storage.GetPhotoUploadData(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalUploadData(newUploadData1, findUploadData)
	})
}

func (s *testSuite) TestAdapter_GetPhotoUploadData() {
	s.Run("ok get upload data", func() {
		photoID := s.savePhoto()

		newUploadData := s.saveUploadData(photoID)

		findUploadData, err := s.storage.GetPhotoUploadData(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalUploadData(newUploadData, findUploadData)
	})

	s.Run("get upload data not found", func() {
		id := uuid.New()
		_, err := s.storage.GetPhotoUploadData(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}
