package app

import (
	"fmt"

	"github.com/kkiling/photo-library/backend/api/internal/adapter/fsstore"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
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
