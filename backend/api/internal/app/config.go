package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/kkiling/photo-library/backend/api/pkg/common/server"
)

const (
	ServerConfigName = "server"
	PgConnectionName = "pg_db"
)

type PgConfig struct {
	ConnString string `yaml:"conn_string"`
}

func (a *App) getServerConfig() (server.Config, error) {

	var config server.Config
	err := a.cfgProvider.PopulateByKey(ServerConfigName, &config)
	if err != nil {
		return server.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getPgConnConfig() (PgConfig, error) {

	var config PgConfig
	err := a.cfgProvider.PopulateByKey(PgConnectionName, &config)
	if err != nil {
		return PgConfig{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) newPgConn(ctx context.Context) (*pgx.Conn, error) {
	ctg, err := a.getPgConnConfig()
	if err != nil {
		return nil, fmt.Errorf("getPgConnConfig: %w", err)
	}

	conn, err := pgx.Connect(ctx, ctg.ConnString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}
