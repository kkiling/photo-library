package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) equalPhoto(newPhoto, findPhoto model.Photo) {
	s.Require().Equal(newPhoto.ID, findPhoto.ID)
	s.Require().Equal(newPhoto.FileKey, findPhoto.FileKey)
	s.Require().Equal(newPhoto.Hash, findPhoto.Hash)
	s.Require().Equal(newPhoto.UpdateAt.Unix(), findPhoto.UpdateAt.Unix())
	s.Require().Equal(newPhoto.Extension, findPhoto.Extension)
}

func (s *testSuite) TestAdapter_SavePhoto() {
	s.Run("ok save new photo", func() {
		id := uuid.New()
		newPhoto := model.Photo{
			ID:        id,
			FileKey:   fmt.Sprintf("%s.jpeg", id.String()),
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err := s.storage.SavePhoto(s.ctx, newPhoto)
		s.Require().NoError(err)

		findPhoto, err := s.storage.GetPhotoById(s.ctx, id)
		s.Require().NoError(err)
		s.equalPhoto(newPhoto, findPhoto)
	})

	s.Run("error save duplicate id", func() {
		id := uuid.New()

		newPhoto1 := model.Photo{
			ID:        id,
			FileKey:   fmt.Sprintf("%s.jpeg", uuid.NewString()),
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		newPhoto2 := model.Photo{
			ID:        id,
			FileKey:   fmt.Sprintf("%s.jpeg", uuid.NewString()),
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err1 := s.storage.SavePhoto(s.ctx, newPhoto1)
		s.Require().NoError(err1)

		err2 := s.storage.SavePhoto(s.ctx, newPhoto2)
		s.Require().Error(err2)
	})

	s.Run("error save duplicate filaName", func() {
		fileName := fmt.Sprintf("%s.jpeg", uuid.NewString())

		newPhoto1 := model.Photo{
			ID:        uuid.New(),
			FileKey:   fileName,
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		newPhoto2 := model.Photo{
			ID:        uuid.New(),
			FileKey:   fileName,
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err1 := s.storage.SavePhoto(s.ctx, newPhoto1)
		s.Require().NoError(err1)

		err2 := s.storage.SavePhoto(s.ctx, newPhoto2)
		s.Require().Error(err2)
	})

	s.Run("error save duplicate hash", func() {
		hash := uuid.NewString()

		newPhoto1 := model.Photo{
			ID:        uuid.New(),
			FileKey:   fmt.Sprintf("%s.jpeg", uuid.NewString()),
			Hash:      hash,
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		newPhoto2 := model.Photo{
			ID:        uuid.New(),
			FileKey:   fmt.Sprintf("%s.jpeg", uuid.NewString()),
			Hash:      hash,
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err1 := s.storage.SavePhoto(s.ctx, newPhoto1)
		s.Require().NoError(err1)

		err2 := s.storage.SavePhoto(s.ctx, newPhoto2)
		s.Require().Error(err2)
	})
}

func (s *testSuite) TestAdapter_GetPhotoById() {
	s.Run("ok get photo", func() {
		id := uuid.New()
		newPhoto := model.Photo{
			ID:        id,
			FileKey:   fmt.Sprintf("%s.jpeg", id.String()),
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err := s.storage.SavePhoto(s.ctx, newPhoto)
		s.Require().NoError(err)

		findPhoto, err := s.storage.GetPhotoById(s.ctx, id)
		s.Require().NoError(err)
		s.equalPhoto(newPhoto, findPhoto)
	})

	s.Run("get photo not found", func() {
		id := uuid.New()
		_, err := s.storage.GetPhotoById(s.ctx, id)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_GetPhotoByFilename() {
	s.Run("ok get photo", func() {
		fileName := fmt.Sprintf("%s.jpeg", uuid.NewString())
		newPhoto := model.Photo{
			ID:        uuid.New(),
			FileKey:   fileName,
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err := s.storage.SavePhoto(s.ctx, newPhoto)
		s.Require().NoError(err)

		findPhoto, err := s.storage.GetPhotoByFilename(s.ctx, fileName)
		s.Require().NoError(err)
		s.equalPhoto(newPhoto, findPhoto)
	})

	s.Run("get photo not found", func() {
		fileName := fmt.Sprintf("%s.jpeg", uuid.NewString())
		_, err := s.storage.GetPhotoByFilename(s.ctx, fileName)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_GetPhotoByHash() {
	s.Run("ok get photo", func() {
		hash := uuid.NewString()
		newPhoto := model.Photo{
			ID:        uuid.New(),
			FileKey:   fmt.Sprintf("%s.jpeg", uuid.NewString()),
			Hash:      hash,
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err := s.storage.SavePhoto(s.ctx, newPhoto)
		s.Require().NoError(err)

		findPhoto, err := s.storage.GetPhotoByHash(s.ctx, hash)
		s.Require().NoError(err)
		s.equalPhoto(newPhoto, findPhoto)
	})

	s.Run("get photo not found", func() {
		hash := uuid.NewString()
		_, err := s.storage.GetPhotoByHash(s.ctx, hash)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}

func (s *testSuite) TestAdapter_MakeNotValidPhoto() {
	s.Run("ok get photo", func() {
		newPhoto := model.Photo{
			ID:        uuid.New(),
			FileKey:   fmt.Sprintf("%s.jpeg", uuid.NewString()),
			Hash:      uuid.NewString(),
			UpdateAt:  time.Now(),
			Extension: model.PhotoExtensionJpeg,
		}

		err := s.storage.SavePhoto(s.ctx, newPhoto)
		s.Require().NoError(err)

		err = s.storage.MakeNotValidPhoto(s.ctx, newPhoto.ID, "photo not valid")
		s.Require().NoError(err)

		_, err = s.storage.GetPhotoById(s.ctx, newPhoto.ID)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)

		_, err = s.storage.GetPhotoByFilename(s.ctx, newPhoto.FileKey)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)

		_, err = s.storage.GetPhotoByHash(s.ctx, newPhoto.Hash)
		s.Require().Error(err)
		s.Require().ErrorIs(err, serviceerr.ErrNotFound)
	})
}
