package server

import (
	"context"

	"github.com/kkiling/photo-library/backend/api/internal/handler/sync_photos_handler"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

const syncPhotosSwaggerName = "sync_photos"

type SyncPhotosServer struct {
	*CustomServer
}

func NewSyncPhotosServer(
	logger log.Logger,
	cfg server.Config,
	syncPhotoService sync_photo_handler.SyncPhotosService,
) *SyncPhotosServer {
	return &SyncPhotosServer{
		CustomServer: NewCustomServer(
			logger.Named("sync_photos_server"),
			cfg,
			sync_photo_handler.NewHandler(
				logger.Named("sync_photo_handler"),
				syncPhotoService,
			),
		),
	}
}

func (s *SyncPhotosServer) Start(ctx context.Context) error {
	return s.CustomServer.Start(ctx, syncPhotosSwaggerName)
}
