package app

import (
	"fmt"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/fsstore"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photo_group"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photo_preview"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/similar_photos"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

// Config
const (
	ServerConfigName     = "server"
	PgConnectionName     = "pg_db"
	FsStoreName          = "fs_store"
	PhotoMLName          = "photo_ml"
	ProcessingPhotosName = "processing_photo"
	SimilarPhotosName    = "similar_photo"
	PhotoGroupName       = "photo_group"
	PhotoPreviewName     = "photo_preview"
)

func (a *App) getServerConfig() (server.Config, error) {

	var config server.Config
	err := a.cfgProvider.PopulateByKey(ServerConfigName, &config)
	if err != nil {
		return server.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPgConnConfig() (pgrepo.PgConfig, error) {

	var config pgrepo.PgConfig
	err := a.cfgProvider.PopulateByKey(PgConnectionName, &config)
	if err != nil {
		return pgrepo.PgConfig{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getFsStoreConfig() (fsstore.Config, error) {
	var config fsstore.Config
	err := a.cfgProvider.PopulateByKey(FsStoreName, &config)
	if err != nil {
		return fsstore.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPhotoMLConfig() (photoml.Config, error) {
	var config photoml.Config
	err := a.cfgProvider.PopulateByKey(PhotoMLName, &config)
	if err != nil {
		return photoml.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getProcessingPhotosConfig() (processing.Config, error) {
	var config processing.Config
	err := a.cfgProvider.PopulateByKey(ProcessingPhotosName, &config)
	if err != nil {
		return processing.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getSimilarPhotosConfig() (similar_photos.Config, error) {
	var config similar_photos.Config
	err := a.cfgProvider.PopulateByKey(SimilarPhotosName, &config)
	if err != nil {
		return similar_photos.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPhotoGroupConfig() (photo_group.Config, error) {
	var config photo_group.Config
	err := a.cfgProvider.PopulateByKey(PhotoGroupName, &config)
	if err != nil {
		return photo_group.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPhotoPreviewConfig() (photo_preview.Config, error) {
	var config photo_preview.Config
	err := a.cfgProvider.PopulateByKey(PhotoPreviewName, &config)
	if err != nil {
		return photo_preview.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}
