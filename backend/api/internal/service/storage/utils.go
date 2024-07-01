package storage

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"time"

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

func toTimePtr(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func toTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t != nil {
		return pgtype.Timestamptz{
			Time:  *t,
			Valid: true,
		}
	}
	return pgtype.Timestamptz{}
}

func isAlreadyExist(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
