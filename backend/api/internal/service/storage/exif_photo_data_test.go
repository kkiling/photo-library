package storage

import (
	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) equalExifData(new, find model.ExifPhotoData) {
	s.Require().Equal(new.PhotoID, new.PhotoID)
	for key, value := range new.Data {
		findValue, ok := find.Data[key]
		s.Require().True(ok)
		s.Require().Equal(value, findValue)
	}
}

func (s *testSuite) TestAdapter_SaveExif() {
	s.Run("ok save exif data", func() {
		photoID := s.savePhoto()

		newExifData := model.ExifPhotoData{
			PhotoID: photoID,
			Data: map[string]interface{}{
				model.ModelExifData:    "Pixel",
				model.DateTimeExifData: "2023-01-02",
			},
		}

		err := s.storage.SaveExif(s.ctx, newExifData)
		s.Require().NoError(err)

		findExifData, err := s.storage.GetExif(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalExifData(newExifData, findExifData)
	})

	s.Run("error save duplicate exif", func() {
		photoID := s.savePhoto()

		newExifData1 := model.ExifPhotoData{
			PhotoID: photoID,
			Data: map[string]interface{}{
				model.ModelExifData:    "Pixel",
				model.DateTimeExifData: "2023-01-02",
			},
		}
		newExifData2 := model.ExifPhotoData{
			PhotoID: photoID,
			Data: map[string]interface{}{
				model.ModelExifData:    "Xiaomi",
				model.DateTimeExifData: "2024-02-05",
			},
		}

		err := s.storage.SaveExif(s.ctx, newExifData1)
		s.Require().NoError(err)
		err = s.storage.SaveExif(s.ctx, newExifData2)
		s.Require().Error(err)

		findExifData, err := s.storage.GetExif(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalExifData(newExifData1, findExifData)
	})
}

func (s *testSuite) TestAdapter_GetExif() {
	s.Run("ok get exif", func() {
		photoID := s.savePhoto()

		newExifData := model.ExifPhotoData{
			PhotoID: photoID,
			Data: map[string]interface{}{
				model.ModelExifData:    "Pixel",
				model.DateTimeExifData: "2023-01-02",
			},
		}

		err := s.storage.SaveExif(s.ctx, newExifData)
		s.Require().NoError(err)

		findExifData, err := s.storage.GetExif(s.ctx, photoID)
		s.Require().NoError(err)
		s.equalExifData(newExifData, findExifData)
	})

	s.Run("get exif not found", func() {
		id := uuid.New()
		_, err := s.storage.GetExif(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_DeleteExif() {
	s.Run("ok delete exif", func() {
		photoID := s.savePhoto()

		newExifData := model.ExifPhotoData{
			PhotoID: photoID,
			Data: map[string]interface{}{
				model.ModelExifData:    "Pixel",
				model.DateTimeExifData: "2023-01-02",
			},
		}

		err := s.storage.SaveExif(s.ctx, newExifData)
		s.Require().NoError(err)

		err = s.storage.DeleteExif(s.ctx, photoID)
		s.Require().NoError(err)

		_, err = s.storage.GetExif(s.ctx, photoID)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})

	s.Run("delete exif not found", func() {
		id := uuid.New()
		err := s.storage.DeleteExif(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}
