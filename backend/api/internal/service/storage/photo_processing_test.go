package storage

import (
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

func (s *testSuite) equalPhotoProcessing(new, find model.PhotoProcessing) {
	s.Require().Equal(new.PhotoID, find.PhotoID)
	s.Require().Equal(new.ProcessingType, find.ProcessingType)
	s.Require().Equal(new.ProcessedAt.Unix(), find.ProcessedAt.Unix())
	s.Require().Equal(new.Success, find.Success)
}

func (s *testSuite) TestAdapter_AddPhotoProcessing() {
	s.Run("ok add photo processing", func() {
		photoID := s.savePhoto()

		processing1 := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now(),
			ProcessingType: model.CatalogTagsProcessing,
			Success:        true,
		}
		processing2 := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now().Add(time.Minute),
			ProcessingType: model.SimilarCoefficientProcessing,
			Success:        false,
		}
		processing3 := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now().Add(2 * time.Minute),
			ProcessingType: model.ExifDataProcessing,
			Success:        true,
		}

		err := s.storage.AddPhotoProcessing(s.ctx, processing1)
		s.Require().NoError(err)
		err = s.storage.AddPhotoProcessing(s.ctx, processing2)
		s.Require().NoError(err)
		err = s.storage.AddPhotoProcessing(s.ctx, processing3)
		s.Require().NoError(err)

		processingTypes, err := s.storage.GetPhotoProcessingTypes(s.ctx, photoID)
		s.Require().NoError(err)
		s.Require().Equal(len(processingTypes), 3)
		s.equalPhotoProcessing(processing1, processingTypes[0])
		s.equalPhotoProcessing(processing2, processingTypes[1])
		s.equalPhotoProcessing(processing3, processingTypes[2])
	})

	s.Run("error unknown photo", func() {
		photoID := uuid.New()

		processing := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now(),
			ProcessingType: model.CatalogTagsProcessing,
			Success:        true,
		}

		err := s.storage.AddPhotoProcessing(s.ctx, processing)
		s.Require().Error(err)
	})

	s.Run("error add photo double processing status", func() {
		photoID := s.savePhoto()

		processing1 := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now(),
			ProcessingType: model.CatalogTagsProcessing,
			Success:        true,
		}
		processing2 := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now().Add(time.Minute),
			ProcessingType: model.CatalogTagsProcessing,
			Success:        false,
		}

		err := s.storage.AddPhotoProcessing(s.ctx, processing1)
		s.Require().NoError(err)
		err = s.storage.AddPhotoProcessing(s.ctx, processing2)
		s.Require().Error(err)
	})
}

func (s *testSuite) TestAdapter_GetUnprocessedPhotoIDs() {
	var processingTypes = []model.ProcessingType{
		model.ExifDataProcessing,
		model.MetaDataProcessing,
	}

	var addProcessing = func(photoID uuid.UUID, processingType model.ProcessingType) {
		addProcessing := model.PhotoProcessing{
			PhotoID:        photoID,
			ProcessedAt:    time.Now().Add(2 * time.Minute),
			ProcessingType: processingType,
			Success:        true,
		}

		err := s.storage.AddPhotoProcessing(s.ctx, addProcessing)
		s.Require().NoError(err)
	}

	s.Run("ok", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()

		photoIDs, err := s.storage.GetUnprocessedPhotoIDs(s.ctx, processingTypes, 10, photoID1, photoID2)
		s.Require().NoError(err)
		s.Require().Equal(2, len(photoIDs))
		s.Require().Equal(photoID1, photoIDs[0])
		s.Require().Equal(photoID2, photoIDs[1])
	})

	s.Run("ok add unprocessed photos after add processing", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()

		addProcessing(photoID1, model.ExifDataProcessing)
		addProcessing(photoID1, model.MetaDataProcessing)

		addProcessing(photoID2, model.MetaDataProcessing)

		photoIDs, err := s.storage.GetUnprocessedPhotoIDs(s.ctx, processingTypes, 10, photoID1, photoID2)
		s.Require().NoError(err)
		s.Require().Equal(1, len(photoIDs))
		s.Require().Equal(photoID2, photoIDs[0])

		addProcessing(photoID2, model.ExifDataProcessing)

		photoIDs, err = s.storage.GetUnprocessedPhotoIDs(s.ctx, processingTypes, 10, photoID1, photoID2)
		s.Require().NoError(err)
		s.Require().Equal(0, len(photoIDs))
	})
}
