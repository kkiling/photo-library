package pgrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type PgConfig struct {
	ConnString string `yaml:"conn_string"`
}

func NewPgConn(ctx context.Context, cfg PgConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.ConnString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return conn, nil
}
