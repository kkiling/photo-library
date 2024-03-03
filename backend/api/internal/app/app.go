package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/fsstore"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/geo"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage"
	pgrepo2 "github.com/kkiling/photo-library/backend/api/internal/adapter/storage/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/handler"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/catalogtags"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/exifphotodata"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/metatags"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photometadata"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/similarphotos"
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
	storageAdapter *storage.Adapter
	pgRepo         *pgrepo2.PhotoRepository
	fsStore        *fsstore.Store
	photoML        *photoml.Service
	// handler
	syncPhotoServer *handler.SyncPhotosServiceServer
	// service
	similarPhotos *similarphotos.Service
	syncPhoto     *syncphotos.Service
	geoService    *geo.Service
	// Processing
	tagPhoto         *tagphoto.Service
	exifPhoto        *exifphotodata.Service
	metaPhoto        *photometadata.Service
	metaTagsPhoto    *metatags.Service
	catalogTagsPhoto *catalogtags.Service
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

	getSimilarPhotosCfg, err := a.getSimilarPhotosConfig()
	if err != nil {
		return fmt.Errorf("getSimilarPhotosConfig: %w", err)
	}

	pool, err := pgrepo2.NewPgConn(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("newPgConn: %w", err)
	}

	a.pgxPool = pool
	a.logger = log.NewLogger()
	a.pgRepo = pgrepo2.NewPhotoRepository(a.pgxPool)
	a.storageAdapter = storage.NewStorageAdapter(a.pgRepo)
	a.fsStore = fsstore.NewStore(fsStoreCfg)
	a.geoService = geo.NewService(a.logger.Named("geo_service"))
	a.photoML = photoml.NewService(
		a.logger.Named("photo_ml"),
		photoMlCfg,
	)
	a.tagPhoto = tagphoto.NewService(
		a.storageAdapter,
	)
	a.syncPhoto = syncphotos.NewService(
		a.logger.Named("sync_photo"),
		a.storageAdapter,
		a.fsStore,
	)
	a.exifPhoto = exifphotodata.NewService(
		a.logger.Named("exif_photo"),
		a.storageAdapter,
	)
	a.metaPhoto = photometadata.NewService(
		a.logger.Named("meta_photo"),
		a.storageAdapter,
	)
	a.metaTagsPhoto = metatags.NewService(
		a.logger.Named("meta_photo_service_photo"),
		a.tagPhoto,
		a.storageAdapter,
		a.geoService,
	)
	a.catalogTagsPhoto = catalogtags.NewService(
		a.logger.Named("catalog_photo_service_photo"),
		a.tagPhoto,
		a.storageAdapter,
	)
	a.syncPhotoServer = handler.NewSyncPhotosServiceServer(
		a.logger.Named("sync_photo_service_photo"),
		a.syncPhoto,
		serverCfg,
	)
	a.vectorPhoto = vectorphoto.NewService(
		a.logger.Named("vector_photo"),
		a.storageAdapter,
		a.photoML,
	)
	a.similarPhotos = similarphotos.NewService(
		a.logger.Named("similar_photos"),
		getSimilarPhotosCfg,
		a.storageAdapter,
	)
	a.processingPhotos = processing.NewService(
		a.logger.Named("processing_photos"),
		processingPhotosCfg,
		a.storageAdapter,
		a.fsStore,
		map[model.PhotoProcessingStatus]processing.PhotoProcessor{
			model.ExifDataProcessing:           a.exifPhoto,
			model.MetaDataProcessing:           a.metaPhoto,
			model.CatalogTagsProcessing:        a.catalogTagsPhoto,
			model.MetaTagsProcessing:           a.metaTagsPhoto,
			model.PhotoVectorProcessing:        a.vectorPhoto,
			model.SimilarCoefficientProcessing: a.similarPhotos,
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
