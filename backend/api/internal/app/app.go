package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkiling/photo-library/backend/api/internal/adapter"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/fsstore"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/internal/service/exifphoto"
	"github.com/kkiling/photo-library/backend/api/internal/service/metaphoto"
	"github.com/kkiling/photo-library/backend/api/internal/service/similarphotos"
	"github.com/kkiling/photo-library/backend/api/internal/service/syncphotos"
	"github.com/kkiling/photo-library/backend/api/internal/service/systags"
	"github.com/kkiling/photo-library/backend/api/internal/service/tagphoto"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type App struct {
	cfgProvider config.Provider
	logger      log.Logger
	// connect
	pgxPool *pgxpool.Pool
	// adapter
	dbAdapter *adapter.DbAdapter
	pgRepo    *pgrepo.PhotoRepository
	fsStore   *fsstore.Store
	photoML   *photoml.Service
	// handler
	syncPhotoServer *handler.SyncPhotosServiceServer
	// service
	tagPhoto *tagphoto.Service

	syncPhoto     *syncphotos.Service
	exifPhoto     *exifphoto.Service
	metaPhoto     *metaphoto.Service
	sysTagPhoto   *systags.Service
	similarPhotos *similarphotos.Service
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

	serverCfg, err := a.getServerConfig()
	if err != nil {
		return fmt.Errorf("getServerConfig: %w", err)
	}

	fsStoreCfg, err := a.getFsStoreConfig()
	if err != nil {
		return fmt.Errorf("getFsStoreConfig: %w", err)
	}

	photoMlCfg, err := a.getPhotoMLConfig()
	if err != nil {
		return fmt.Errorf("getPhotoMLConfig: %w", err)
	}

	pool, err := pgrepo.NewPgConn(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("newPgConn: %w", err)
	}

	a.pgxPool = pool
	a.logger = log.NewLogger()
	a.pgRepo = pgrepo.NewPhotoRepository(a.pgxPool)
	a.dbAdapter = adapter.NewDbAdapter(a.pgRepo)
	a.fsStore = fsstore.NewStore(fsStoreCfg)
	a.photoML = photoml.NewService(
		a.logger.Named("photo_ml"),
		photoMlCfg,
	)
	a.tagPhoto = tagphoto.NewService(
		a.dbAdapter,
	)

	a.syncPhoto = syncphotos.NewService(
		a.logger.Named("sync_photo"),
		a.dbAdapter,
		a.fsStore,
	)

	a.exifPhoto = exifphoto.NewService(
		a.dbAdapter,
	)

	a.metaPhoto = metaphoto.NewService(
		a.dbAdapter,
	)

	a.sysTagPhoto = systags.NewService(
		a.logger.Named("sync_photo_service_photo"),
		a.tagPhoto,
		a.dbAdapter,
	)

	a.syncPhotoServer = handler.NewSyncPhotosServiceServer(
		a.logger.Named("sync_photo_service_photo"),
		a.syncPhoto,
		serverCfg,
	)

	a.similarPhotos = similarphotos.NewService(
		a.logger.Named("similar_photos"),
		a.tagPhoto,
		a.dbAdapter,
		a.photoML,
	)

	return nil
}

func (a *App) StartSyncPhotoServer(ctx context.Context) error {
	return a.syncPhotoServer.Start(ctx)
}

func (a *App) StopSyncPhotoServer() {
	a.syncPhotoServer.Stop()
}

func (a *App) GetExifPhoto() *exifphoto.Service {
	return a.exifPhoto
}

func (a *App) GetMetaPhoto() *metaphoto.Service {
	return a.metaPhoto
}

func (a *App) GetDbAdapter() *adapter.DbAdapter {
	return a.dbAdapter
}

func (a *App) GetFileStorage() *fsstore.Store {
	return a.fsStore
}

func (a *App) GetSysTagPhoto() *systags.Service {
	return a.sysTagPhoto
}

func (a *App) GetSimilarPhotos() *similarphotos.Service {
	return a.similarPhotos
}

func (a *App) GetLogger() log.Logger {
	return a.logger
}
