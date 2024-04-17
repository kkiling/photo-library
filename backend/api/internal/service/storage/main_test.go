package storage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

type testSuite struct {
	suite.Suite
	ctx     context.Context
	ctrl    *gomock.Controller
	storage *Adapter
}

func TestStorage(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
}

func (s *testSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *testSuite) SetupSuite() {
	s.ctx = context.Background()
	pool, err := pgrepo.NewPgConn(s.ctx, pgrepo.PgConfig{
		ConnString: getDsn(),
	})
	if err != nil {
		panic(fmt.Errorf("pgrepo.NewPgConn: %w", err))
	}
	s.storage = NewStorageAdapter(pool)
}

func getDsn() string {
	if dsn := os.Getenv("PG_DSN"); dsn != "" {
		return dsn
	}
	return "postgresql://postgres@localhost:5432/photo_library_test?sslmode=disable"
}

func (s *testSuite) savePhoto() uuid.UUID {
	photoID := uuid.New()

	newPhoto := model.Photo{
		ID:        photoID,
		FileKey:   fmt.Sprintf("%s.jpeg", photoID.String()),
		Hash:      uuid.NewString(),
		UpdateAt:  time.Now(),
		Extension: model.PhotoExtensionJpeg,
	}

	err := s.storage.SavePhoto(s.ctx, newPhoto)
	s.Require().NoError(err)
	return photoID
}

func (s *testSuite) compareTime(t1, t2 *time.Time) {
	if t1 != nil && t2 != nil {
		s.Require().Equal(t1.Unix(), t2.Unix())
		return
	}
	s.Require().Equal(t1, t2)
}
