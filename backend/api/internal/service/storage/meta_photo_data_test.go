package storage

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) equalPhotoMetadata(new, find model.PhotoMetadata) {
	s.Require().Equal(new.PhotoID, find.PhotoID)
	s.Require().Equal(new.ModelInfo, find.ModelInfo)
	s.Require().Equal(new.SizeBytes, find.SizeBytes)
	s.Require().Equal(new.WidthPixel, find.WidthPixel)
	s.Require().Equal(new.HeightPixel, find.HeightPixel)
	s.compareTime(new.DateTime, find.DateTime)
	s.compareTime(&new.PhotoUpdatedAt, &find.PhotoUpdatedAt)
	if new.Geo != nil && find.Geo != nil {
		s.Require().Equal(new.Geo.Longitude, find.Geo.Longitude)
		s.Require().Equal(new.Geo.Latitude, find.Geo.Latitude)
	} else {
		s.Require().Equal(new.Geo, find.Geo)
	}
}

func (s *testSuite) createMetaData(photoID uuid.UUID) model.PhotoMetadata {
	newMetadata := model.PhotoMetadata{
		PhotoID:        photoID,
		ModelInfo:      lo.ToPtr("Model Info"),
		SizeBytes:      1024,
		WidthPixel:     512,
		HeightPixel:    256,
		DateTime:       lo.ToPtr(time.Now()),
		PhotoUpdatedAt: time.Now(),
		Geo: &model.Geo{
			Latitude:  52.1234,
			Longitude: 53.56,
		},
	}

	return newMetadata
}

func (s *testSuite) saveMetaData(photoID uuid.UUID) model.PhotoMetadata {
	newMetadata := s.createMetaData(photoID)

	err := s.storage.SaveMetadata(s.ctx, newMetadata)
	s.Require().NoError(err)

	return newMetadata
}

func (s *testSuite) TestAdapter_SaveMetadata() {
	s.Run("ok save meta data", func() {
		photoID := s.savePhoto()
		newMetadata := s.saveMetaData(photoID)

		findMetadata, err := s.storage.GetMetadata(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalPhotoMetadata(newMetadata, findMetadata)
	})

	s.Run("error save duplicate metadata", func() {
		photoID := s.savePhoto()

		newMetadata1 := s.createMetaData(photoID)
		err := s.storage.SaveMetadata(s.ctx, newMetadata1)
		s.Require().NoError(err)

		newMetadata2 := s.createMetaData(photoID)
		err = s.storage.SaveMetadata(s.ctx, newMetadata2)
		s.Require().Error(err)

		findMetaData, err := s.storage.GetMetadata(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalPhotoMetadata(newMetadata1, findMetaData)
	})
}

func (s *testSuite) TestAdapter_GetMetadata() {
	s.Run("ok get metadata", func() {
		photoID := s.savePhoto()

		newMetaData := s.saveMetaData(photoID)

		findMetaData, err := s.storage.GetMetadata(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalPhotoMetadata(newMetaData, findMetaData)
	})

	s.Run("get metadata not found", func() {
		id := uuid.New()
		_, err := s.storage.GetMetadata(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_DeleteMetadata() {
	s.Run("ok delete exif", func() {
		photoID := s.savePhoto()
		_ = s.saveMetaData(photoID)

		err := s.storage.DeleteMetadata(s.ctx, photoID)
		s.Require().NoError(err)

		_, err = s.storage.GetMetadata(s.ctx, photoID)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})

	s.Run("delete metadata not found", func() {
		id := uuid.New()
		err := s.storage.DeleteMetadata(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}
