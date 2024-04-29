package server

import (
	"context"
	"github.com/kkiling/photo-library/backend/api/internal/interceptor"

	"github.com/kkiling/photo-library/backend/api/internal/handler/sync_photos_handler"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

const syncPhotosSwaggerName = "sync_photos"

type SyncPhotosServer struct {
	*ApiServer
}

func NewSyncPhotosServer(
	logger log.Logger,
	cfg server.Config,
	apiToken interceptor.ApiTokenService,
	syncPhotoService sync_photo_handler.SyncPhotosService,
) *SyncPhotosServer {
	return &SyncPhotosServer{
		ApiServer: NewApiServer(
			logger.Named("sync_photos_server"),
			cfg,
			apiToken,
			sync_photo_handler.NewHandler(
				logger.Named("sync_photos_handler"),
				syncPhotoService,
			),
		),
	}
}

func (s *SyncPhotosServer) Start(ctx context.Context) error {
	return s.ApiServer.Start(ctx, syncPhotosSwaggerName)
}
