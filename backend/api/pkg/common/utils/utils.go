package utils

import (
	"reflect"
	"runtime"
	"strings"
)

// GetFunctionName получить имя функции
func GetFunctionName(function any) (name string) {
	fullName := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	elements := strings.Split(fullName, "/")
	shortName := elements[len(elements)-1]
	fnName := strings.TrimSuffix(shortName, "-fm")
	return fnName
}

func TransformToName(str string) string {
	str = strings.TrimSpace(str)
	// Если строка пуста, просто возвращаем её
	if len(str) == 0 {
		return str
	}

	// Преобразуем первое слово в строке
	firstWord := strings.SplitN(str, " ", 2)[0]
	capitalizedFirst := strings.Title(strings.ToLower(firstWord))

	// Преобразуем оставшуюся часть строки в нижний регистр
	rest := strings.ToLower(str[len(firstWord):])

	// Объединяем и возвращаем результат
	return capitalizedFirst + rest
}

func TransformToNamePtr(strPtr *string) *string {
	if strPtr == nil {
		return nil
	}

	str := TransformToName(*strPtr)

	if str == "" {
		return nil
	}
	return &str
}
