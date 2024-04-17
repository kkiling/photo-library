package storage

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func stringPtr(v pgtype.Text) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}

func dataTimePtr(v pgtype.Timestamptz) *time.Time {
	if v.Valid {
		return &v.Time
	}
	return nil
}

func pgTypeText(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{
		String: *v,
		Valid:  true,
	}
}

func pgTypeTimestamptz(v *time.Time) pgtype.Timestamptz {
	if v == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{
		Time:  *v,
		Valid: true,
	}
}

func pgTypeFloat8(v *float64) pgtype.Float8 {
	if v == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{
		Float64: *v,
		Valid:   true,
	}
}
