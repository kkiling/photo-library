package exifphotodata

import (
	"fmt"
	"reflect"
)

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	fieldVal := structValue.FieldByName(name)

	if !fieldVal.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	val := reflect.ValueOf(value)
	// Если тип поля является указателем, получим соответствующий тип без указателя
	if fieldVal.Type().Kind() == reflect.Ptr {
		if val.Type() == fieldVal.Type().Elem() {
			valPtr := reflect.New(val.Type())
			valPtr.Elem().Set(val)
			fieldVal.Set(valPtr)
			return nil
		}
	} else if fieldVal.Type() == val.Type() {
		fieldVal.Set(val)
		return nil
	}

	return fmt.Errorf("provided value type (%s) didn't match obj field type (%s)", val.Type(), fieldVal.Type())
}

func determineDataType(obj interface{}, fieldName string) dataType {
	structValue := reflect.ValueOf(obj).Elem()
	fieldVal := structValue.FieldByName(fieldName)

	kind := fieldVal.Type().Kind()
	if kind == reflect.Ptr {
		kind = fieldVal.Type().Elem().Kind()
	}

	switch kind {
	case reflect.Int:
		return dataTypeInt
	case reflect.String:
		return dataTypeString
	case reflect.Float32, reflect.Float64:
		return dataTypeFloat
	case reflect.Slice:
		kind = fieldVal.Type().Elem().Kind()
		switch kind {
		case reflect.Int:
			return dataTypeIntArray
		case reflect.Float32, reflect.Float64:
			return dataTypeFloatArray
		}
	}

	return dataTypeUnknown
}
