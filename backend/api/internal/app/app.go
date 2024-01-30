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
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/exifphoto"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/metaphoto"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/similarphotos"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/systags"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/vectorphoto"
	"github.com/kkiling/photo-library/backend/api/internal/service/syncphotos"
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
	similarPhotos *similarphotos.Service
	syncPhoto     *syncphotos.Service
	// Processing
	tagPhoto         *tagphoto.Service
	exifPhoto        *exifphoto.Service
	metaPhoto        *metaphoto.Service
	sysTagPhoto      *systags.Service
	vectorPhoto      *vectorphoto.Service
	processingPhotos *processing.Service
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

	processingPhotosCfg, err := a.getProcessingPhotosConfig()
	if err != nil {
		return fmt.Errorf("getProcessingPhotosConfig: %w", err)
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
		a.logger.Named("exif_photo"),
		a.dbAdapter,
	)
	a.metaPhoto = metaphoto.NewService(
		a.logger.Named("meta_photo"),
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
	a.vectorPhoto = vectorphoto.NewService(
		a.logger.Named("vector_photo"),
		a.dbAdapter,
		a.photoML,
	)
	a.similarPhotos = similarphotos.NewService(
		a.logger.Named("similar_photos"),
		a.dbAdapter,
	)
	a.processingPhotos = processing.NewService(
		a.logger.Named("processing_photos"),
		processingPhotosCfg,
		a.dbAdapter,
		a.fsStore,
		map[model.PhotoProcessingStatus]processing.PhotoProcessor{
			model.PhotoProcessingNew:         nil, // Стартовая точка каждой фотографии
			model.PhotoProcessingExifData:    a.exifPhoto,
			model.PhotoProcessingMetaData:    a.metaPhoto,
			model.PhotoProcessingTagsByMeta:  a.sysTagPhoto,
			model.PhotoProcessingPhotoVector: a.vectorPhoto,
		},
	)

	return nil
}

func (a *App) StartSyncPhotoServer(ctx context.Context) error {
	return a.syncPhotoServer.Start(ctx)
}

func (a *App) StopSyncPhotoServer() {
	a.syncPhotoServer.Stop()
}

func (a *App) GetProcessingPhotos() *processing.Service {
	return a.processingPhotos
}

func (a *App) GetLogger() log.Logger {
	return a.logger
}
