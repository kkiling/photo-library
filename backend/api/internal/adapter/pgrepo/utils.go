package pgrepo

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func printError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		return newErr
	}
	return err
}
