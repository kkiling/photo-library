package utils

import (
	"path/filepath"
	"strings"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

func GetPhotoExtension(path string) *model.PhotoExtension {
	// Извлекаем расширение файла из пути.
	ext := strings.ToUpper(filepath.Ext(path))
	// Удаляем точку из начала расширения.
	ext = strings.TrimPrefix(ext, ".")

	// Определяем, какому PhotoExtension соответствует извлеченное расширение.
	switch ext {
	case "JPEG", "JPG":
		photoExt := model.PhotoExtensionJpeg
		return &photoExt
	case string(model.PhotoExtensionPng):
		photoExt := model.PhotoExtensionPng
		return &photoExt
	/*case string(model.PhotoExtensionBmb):
	photoExt := model.PhotoExtensionBmb
	return &photoExt*/
	default:
		// Если расширение не соответствует известным типам, возвращаем nil.
		return nil
	}
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
