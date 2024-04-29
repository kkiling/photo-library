package app

import (
	"context"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth/api_token"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth/codes"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth/jwt_helper"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth/password"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth/session_manager"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/fsstore"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/geo"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/service/lock"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos/photo_groups_service"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos/photo_metadata_service"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos/photo_tags_service"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/catalog_tags"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/exif_photo_data"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/meta_tags"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photo_group"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photo_metadata"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photo_preview"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/similar_photos"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/vector_photo"
	"github.com/kkiling/photo-library/backend/api/internal/service/storage"
	"github.com/kkiling/photo-library/backend/api/internal/service/sync_photos"
	"github.com/kkiling/photo-library/backend/api/internal/service/tags"
	"github.com/kkiling/photo-library/backend/api/pkg/common/config"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

type App struct {
	cfgProvider config.Provider
	logger      log.Logger
	// cfg
	serverCfg server.Config
	// connect
	pgxPool *pgxpool.Pool
	// adapter
	storageAdapter *storage.Adapter
	fsStore        *fsstore.Store
	photoML        *photoml.Service
	geoService     *geo.Service
	// service
	lockService           *lock.Service
	tagsServices          *tags.Service
	photoGroupService     *photo_groups_service.Service
	photoMetadataService  *photo_metadata_service.Service
	photoTagsService      *photo_tags_service.Service
	syncPhotoService      *sync_photos.Service
	authService           *auth.Service
	confirmCodeService    *codes.Service
	passwordService       *password.Service
	sessionManagerService *session_manager.SessionManager
	jwtHelper             *jwt_helper.JwtHelper
	apiTokenService       *api_token.Service
	// Processing
	similarPhotos    *similar_photos.Processing
	exifPhoto        *exif_photo_data.Processing
	metaPhoto        *photo_metadata.Processing
	metaTagsPhoto    *meta_tags.Processing
	catalogTagsPhoto *catalog_tags.Processing
	vectorPhoto      *vector_photo.Processing
	photoGroup       *photo_group.Processing
	photoPreview     *photo_preview.Processing
	processingPhotos *processing.Service
}

func NewApp(cfgProvider config.Provider) *App {
	return &App{cfgProvider: cfgProvider}
}

func (a *App) Create(ctx context.Context) error {
	pgCfg, err := a.getPgConnConfig()
	if err != nil {
		return fmt.Errorf("getPgConnConfig: %w", err)
	}

	a.serverCfg, err = a.getServerConfig()
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

	photoGroupCfg, err := a.getPhotoGroupConfig()
	if err != nil {
		return fmt.Errorf("getPhotoGroupConfig: %w", err)
	}

	photoPreviewCfg, err := a.getPhotoPreviewConfig()
	if err != nil {
		return fmt.Errorf("getPhotoPreviewConfig: %w", err)
	}

	authCfg, err := a.getAuthConfig()
	if err != nil {
		return fmt.Errorf("getAuthConfig: %w", err)
	}

	sessionManagerCfg, err := a.getSessionManagerConfig()
	if err != nil {
		return fmt.Errorf("getSessionManagerConfig: %w", err)
	}

	jwtHelperCfg, err := a.getJwtHelperConfig()
	if err != nil {
		return fmt.Errorf("getJwtHelperConfig: %w", err)
	}

	pool, err := pgrepo.NewPgConn(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("newPgConn: %w", err)
	}

	a.jwtHelper, err = jwt_helper.NewHelper(jwtHelperCfg)
	if err != nil {
		return fmt.Errorf("jwt_helper.NewHelper: %w", err)
	}

	a.pgxPool = pool
	a.logger = log.NewLogger()
	a.storageAdapter = storage.NewStorageAdapter(a.pgxPool)
	a.fsStore = fsstore.NewStore(fsStoreCfg)
	a.geoService = geo.NewService(a.logger.Named("geo_service"))
	a.lockService = lock.NewService(a.storageAdapter)
	a.confirmCodeService = codes.NewService(
		a.logger.Named("confirm_code_service"),
		a.storageAdapter,
	)
	a.apiTokenService = api_token.NewService(
		a.logger.Named("api_token_service"),
		a.storageAdapter,
	)
	a.passwordService = password.NewService(a.logger.Named("password_service"))
	a.sessionManagerService = session_manager.NewSessionManager(
		a.logger.Named("confirm_code_service"),
		sessionManagerCfg,
		a.jwtHelper,
	)
	a.photoML = photoml.NewService(
		a.logger.Named("photo_ml"),
		photoMlCfg,
	)
	a.tagsServices = tags.NewService(
		a.storageAdapter,
	)
	a.photoGroupService = photo_groups_service.NewService(
		a.logger.Named("photo_groups_service"),
		a.fsStore,
		a.storageAdapter,
	)
	a.photoMetadataService = photo_metadata_service.NewService(
		a.logger.Named("photo_metadata_service"),
		a.storageAdapter,
	)
	a.photoTagsService = photo_tags_service.NewService(
		a.logger.Named("photo_tags_service"),
		a.tagsServices,
		a.storageAdapter,
	)
	a.authService = auth.NewService(
		a.logger.Named("auth_service"),
		a.storageAdapter,
		authCfg,
		a.confirmCodeService,
		a.passwordService,
		a.sessionManagerService,
	)
	a.syncPhotoService = sync_photos.NewService(
		a.logger.Named("sync_photo"),
		a.storageAdapter,
		a.fsStore,
	)
	a.exifPhoto = exif_photo_data.NewService(
		a.logger.Named("exif_photo"),
		a.storageAdapter,
	)
	a.metaPhoto = photo_metadata.NewService(
		a.logger.Named("meta_photo"),
		a.storageAdapter,
	)
	a.metaTagsPhoto = meta_tags.NewService(
		a.logger.Named("meta_photo_service_photo"),
		a.tagsServices,
		a.storageAdapter,
		a.geoService,
	)
	a.catalogTagsPhoto = catalog_tags.NewService(
		a.logger.Named("catalog_photo_service_photo"),
		a.tagsServices,
		a.storageAdapter,
	)
	a.vectorPhoto = vector_photo.NewService(
		a.logger.Named("vector_photo"),
		a.storageAdapter,
		a.photoML,
	)
	a.photoGroup = photo_group.NewService(
		a.logger.Named("photo_group"),
		photoGroupCfg,
		a.storageAdapter,
		a.lockService,
	)
	a.photoPreview = photo_preview.NewService(
		a.logger.Named("photo_preview"),
		photoPreviewCfg,
		a.storageAdapter,
		a.fsStore,
	)
	a.similarPhotos = similar_photos.NewService(
		a.logger.Named("similar_photos"),
		getSimilarPhotosCfg,
		a.storageAdapter,
		a.lockService,
	)
	a.processingPhotos = processing.NewService(
		a.logger.Named("processing_photos"),
		processingPhotosCfg,
		a.storageAdapter,
		a.fsStore,
		a.lockService,
		map[model.ProcessingType]processing.PhotoProcessor{
			model.ExifDataProcessing:           a.exifPhoto,
			model.MetaDataProcessing:           a.metaPhoto,
			model.CatalogTagsProcessing:        a.catalogTagsPhoto,
			model.MetaTagsProcessing:           a.metaTagsPhoto,
			model.PhotoVectorProcessing:        a.vectorPhoto,
			model.SimilarCoefficientProcessing: a.similarPhotos,
			model.PhotoPreviewProcessing:       a.photoPreview,
			model.PhotoGroupProcessing:         a.photoGroup,
		},
	)

	return nil
}

func (a *App) GetLogger() log.Logger {
	return a.logger
}

func (a *App) GetSessionManager() *session_manager.SessionManager {
	return a.sessionManagerService
}

func (a *App) GetProcessingPhotos() *processing.Service {
	return a.processingPhotos
}

func (a *App) GetServerConfig() server.Config {
	return a.serverCfg
}

func (a *App) GetLockService() *lock.Service {
	return a.lockService
}

func (a *App) GetTagsServices() *tags.Service {
	return a.tagsServices
}

func (a *App) GetPhotoGroupService() *photo_groups_service.Service {
	return a.photoGroupService
}

func (a *App) GetPhotoMetadataService() *photo_metadata_service.Service {
	return a.photoMetadataService
}

func (a *App) GetAuthService() *auth.Service {
	return a.authService
}

func (a *App) ApiTokenService() *api_token.Service {
	return a.apiTokenService
}

func (a *App) GetPhotoTagsService() *photo_tags_service.Service {
	return a.photoTagsService
}

func (a *App) GetSyncPhotoService() *sync_photos.Service {
	return a.syncPhotoService
}
