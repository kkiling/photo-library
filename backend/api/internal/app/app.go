package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/kkiling/photo-library/backend/api/internal/adapter"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/fsstore"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/internal/service/syncphotos"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type App struct {
	cfgProvider     config.Provider
	logger          log.Logger
	pgxConn         *pgx.Conn
	dbAdapter       *adapter.DbAdapter
	pgRepo          *pgrepo.PhotoRepository
	fsStore         *fsstore.Store
	syncPhoto       *syncphotos.Service
	syncPhotoServer *handler.SyncPhotosServiceServer
}

func NewApp(cfgProvider config.Provider) *App {
	return &App{cfgProvider: cfgProvider}
}

func (a *App) Logger() log.Logger {
	return a.logger.Named("application")
}

func (a *App) Create(ctx context.Context) error {
	pgCfg, err := a.getPgConnConfig()
	if err != nil {
		return fmt.Errorf("getPgConnConfig: %w", err)
	}

	conn, err := pgrepo.NewPgConn(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("newPgConn: %w", err)
	}

	a.pgxConn = conn
	a.logger = log.NewLogger()
	a.pgRepo = pgrepo.NewPhotoRepository(a.pgxConn)
	a.dbAdapter = adapter.NewDbAdapter(a.pgRepo)

	a.syncPhoto = syncphotos.NewService(
		a.logger.Named("sync_photo"),
		a.dbAdapter,
		a.fsStore,
	)

	serverCfg, err := a.getServerConfig()
	if err != nil {
		return fmt.Errorf("getServerConfig: %w", err)
	}

	a.syncPhotoServer = handler.NewSyncPhotosServiceServer(
		a.logger.Named("sync_photo_service_photo"),
		a.syncPhoto,
		serverCfg,
	)

	return nil
}

func (a *App) Start(ctx context.Context) error {
	return a.syncPhotoServer.Start(ctx)
}

func (a *App) Stop() {
	a.syncPhotoServer.Stop()
}
