package pgrepo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/kkiling/photo-library/backend/api/internal/adapter/entity"
	"reflect"
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

func (r *PhotoRepository) DeleteExif(ctx context.Context, photoId uuid.UUID) error {
	conn := r.getConn(ctx)

	query, args, err := sq.
		Delete("exif_data").
		Where(sq.Eq{"id": photoId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (r *PhotoRepository) SaveExif(ctx context.Context, data *entity.ExifData) error {
	conn := r.getConn(ctx)
	fields := structToMapDBTag(data)
	query, args, err := sq.
		Insert("exif_data").
		SetMap(fields).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (r *PhotoRepository) GetExif(ctx context.Context, photoId uuid.UUID) (*entity.ExifData, error) {
	panic("not implement")
}
