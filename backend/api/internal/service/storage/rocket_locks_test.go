package storage

import (
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

func (s *testSuite) TestStorage_RocketLock() {
	s.Run("ok", func() {
		key := uuid.NewString()
		ttl := 10000 * time.Millisecond

		// locks first time
		_, err := s.storage.RocketLock(s.ctx, key, ttl)
		s.Require().NoError(err)

		// Key is locked
		_, err = s.storage.RocketLock(s.ctx, key, ttl)
		s.ErrorIs(err, serviceerr.ErrAlreadyLocked)
	})
	s.Run("ttl", func() {
		key := uuid.NewString()
		ttl := 10 * time.Millisecond

		// locks first time
		_, err := s.storage.RocketLock(s.ctx, key, ttl)
		s.Require().NoError(err)
		time.Sleep(11 * time.Millisecond)

		_, err = s.storage.RocketLock(s.ctx, key, ttl)
		s.Require().NoError(err)
	})
}

func (s *testSuite) TestStorage_RocketDeleteLock() {
	s.Run("ok", func() {
		key := uuid.NewString()
		ttl := time.Hour

		id, err := s.storage.RocketLock(s.ctx, key, ttl)
		s.Require().NoError(err)
		s.Require().NotEmpty(id)

		// preserves others locks by lock id
		err = s.storage.RocketLockDelete(s.ctx, model.RocketLockID{Key: key, Ts: 0})
		s.Require().NoError(err)

		_, err = s.storage.RocketLock(s.ctx, key, ttl)
		s.Require().Error(err)

		err = s.storage.RocketLockDelete(s.ctx, id)
		s.Require().NoError(err)

		// can lock again before ttl
		_, err = s.storage.RocketLock(s.ctx, key, ttl)
		s.Require().NoError(err)
	})

	s.Run("double", func() {
		key := uuid.NewString()
		ttl := time.Hour

		id, err := s.storage.RocketLock(s.ctx, key, ttl)
		s.Require().NoError(err)
		s.Require().NotEmpty(id)

		err = s.storage.RocketLockDelete(s.ctx, id)
		s.Require().NoError(err)

		err = s.storage.RocketLockDelete(s.ctx, id)
		s.Require().NoErrorf(err, "double delete is not error")
	})
}
