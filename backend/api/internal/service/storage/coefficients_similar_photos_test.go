package storage

import (
	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

func (s *testSuite) TestAdapter_FindCoefficientSimilarPhoto() {
	s.Run("ok save coefficient", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()

		newCoefficient := model.CoefficientSimilarPhoto{
			PhotoID1:    photoID1,
			PhotoID2:    photoID2,
			Coefficient: 0.42421,
		}

		err := s.storage.SaveCoefficientSimilarPhotos(s.ctx, newCoefficient)
		s.Require().NoError(err)
	})

	s.Run("error save duplicate coefficient", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()

		newCoefficient1 := model.CoefficientSimilarPhoto{
			PhotoID1:    photoID1,
			PhotoID2:    photoID2,
			Coefficient: 0.42421,
		}

		err := s.storage.SaveCoefficientSimilarPhotos(s.ctx, newCoefficient1)
		s.Require().NoError(err)

		newCoefficient2 := model.CoefficientSimilarPhoto{
			PhotoID1:    photoID1,
			PhotoID2:    photoID2,
			Coefficient: 0.52421,
		}

		err = s.storage.SaveCoefficientSimilarPhotos(s.ctx, newCoefficient2)
		s.Require().Error(err)
	})
}

func (s *testSuite) TestAdapter_SaveCoefficientSimilarPhotos() {
	photoID1 := s.savePhoto()
	photoID2 := s.savePhoto()
	photoID3 := s.savePhoto()

	newCoefficient1 := model.CoefficientSimilarPhoto{
		PhotoID1:    photoID1,
		PhotoID2:    photoID2,
		Coefficient: 0.42421,
	}

	err := s.storage.SaveCoefficientSimilarPhotos(s.ctx, newCoefficient1)
	s.Require().NoError(err)

	newCoefficient2 := model.CoefficientSimilarPhoto{
		PhotoID1:    photoID1,
		PhotoID2:    photoID3,
		Coefficient: 0.52421,
	}

	err = s.storage.SaveCoefficientSimilarPhotos(s.ctx, newCoefficient2)
	s.Require().NoError(err)

	newCoefficient3 := model.CoefficientSimilarPhoto{
		PhotoID1:    photoID2,
		PhotoID2:    photoID3,
		Coefficient: 0.52421,
	}

	err = s.storage.SaveCoefficientSimilarPhotos(s.ctx, newCoefficient3)
	s.Require().NoError(err)

	find, err := s.storage.FindCoefficientSimilarPhoto(s.ctx, photoID1)
	s.Require().NoError(err)
	s.Require().Equal(2, len(find))
	s.Require().Equal(find[0].PhotoID1, newCoefficient1.PhotoID1)
	s.Require().Equal(find[0].PhotoID2, newCoefficient1.PhotoID2)
	s.Require().Equal(find[0].Coefficient, newCoefficient1.Coefficient)

	s.Require().Equal(find[1].PhotoID1, newCoefficient2.PhotoID1)
	s.Require().Equal(find[1].PhotoID2, newCoefficient2.PhotoID2)
	s.Require().Equal(find[1].Coefficient, newCoefficient2.Coefficient)
}

func (s *testSuite) TestAdapter_DeleteCoefficientSimilarPhoto() {
	s.Run("ok delete CoefficientSimilarPhoto", func() {
		photoID1 := s.savePhoto()
		photoID2 := s.savePhoto()

		newCoefficient := model.CoefficientSimilarPhoto{
			PhotoID1:    photoID1,
			PhotoID2:    photoID2,
			Coefficient: 0.42421,
		}

		err := s.storage.SaveCoefficientSimilarPhotos(s.ctx, newCoefficient)
		s.Require().NoError(err)

		err = s.storage.DeleteCoefficientSimilarPhoto(s.ctx, photoID2)
		s.Require().NoError(err)

		find1, err := s.storage.FindCoefficientSimilarPhoto(s.ctx, photoID1)
		s.Require().NoError(err)
		s.Require().Len(find1, 0)

		find2, err := s.storage.FindCoefficientSimilarPhoto(s.ctx, photoID2)
		s.Require().NoError(err)
		s.Require().Len(find2, 0)
	})

	s.Run("delete CoefficientSimilarPhoto not found", func() {
		err := s.storage.DeleteCoefficientSimilarPhoto(s.ctx, uuid.New())
		s.Require().NoError(err)
	})
}
