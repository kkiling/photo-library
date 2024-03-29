package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/storage/entity"
)

func structToMapDBTag(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	// Если obj - указатель, получаем элемент, на который он указывает
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)

		// Если значение поля nil, пропускаем его
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		// Если значение поля - пустая строка, пропускаем ее
		if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
			continue
		}

		// Если значение поля - пустой массив или слайс, пропускаем его
		if (fieldValue.Kind() == reflect.Array || fieldValue.Kind() == reflect.Slice) && fieldValue.Len() == 0 {
			continue
		}

		tag := t.Field(i).Tag.Get("db")
		if tag == "" {
			tag = t.Field(i).Tag.Get("json")
		}

		if tag == "" {
			tag = t.Field(i).Name
		}

		m[tag] = fieldValue.Interface()
	}
	return m
}

func (r *PhotoRepository) SaveExif(ctx context.Context, data *entity.ExifPhotoData) error {
	conn := r.getConn(ctx)

	fields := structToMapDBTag(data)
	updateParts := make([]string, 0, len(fields))
	for column := range fields {
		if column == "photo_id" {
			continue
		}
		updateParts = append(updateParts, fmt.Sprintf("%s = EXCLUDED.%s", column, column))
	}

	query, args, err := sq.
		Insert("exif_photo_data").
		SetMap(fields).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return printError(err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return printError(err)
	}

	return nil
}

func (r *PhotoRepository) GetExif(ctx context.Context, photoId uuid.UUID) (*entity.ExifPhotoData, error) {
	conn := r.getConn(ctx)

	var exif entity.ExifPhotoData

	err := pgxscan.Get(ctx, conn, &exif, `SELECT * FROM exif_photo_data WHERE photo_id = $1`, photoId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, printError(err)
	}

	return &exif, nil
}
