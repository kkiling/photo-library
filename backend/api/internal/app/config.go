package app

import (
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/pgrepo"
	"github.com/kkiling/photo-library/backend/api/internal/service/photos"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photogroup"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/photopreview"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing/similarphotos"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/fsstore"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/photoml"
	"github.com/kkiling/photo-library/backend/api/internal/service/processing"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

const (
	ServerConfigName     = "server"
	PgConnectionName     = "pg_db"
	FsStoreName          = "fs_store"
	PhotoMLName          = "photo_ml"
	ProcessingPhotosName = "processing_photo"
	SimilarPhotosName    = "similar_photo"
	PhotoGroupName       = "photo_group"
	PhotosName           = "photos"
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

func (a *App) getSimilarPhotosConfig() (similarphotos.Config, error) {
	var config similarphotos.Config
	err := a.cfgProvider.PopulateByKey(SimilarPhotosName, &config)
	if err != nil {
		return similarphotos.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPhotoGroupConfig() (photogroup.Config, error) {
	var config photogroup.Config
	err := a.cfgProvider.PopulateByKey(PhotoGroupName, &config)
	if err != nil {
		return photogroup.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPhotosConfig() (photos.Config, error) {
	var config photos.Config
	err := a.cfgProvider.PopulateByKey(PhotosName, &config)
	if err != nil {
		return photos.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPhotoPreviewConfig() (photopreview.Config, error) {
	var config photopreview.Config
	err := a.cfgProvider.PopulateByKey(PhotoPreviewName, &config)
	if err != nil {
		return photopreview.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}
