package app

import (
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/pgrepo"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

const (
	ServerConfigName = "server"
	PgConnectionName = "pg_db"
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
