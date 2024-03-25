package utils

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"path/filepath"
	"strings"
)

func GetPhotoExtension(path string) *model.PhotoExtension {
	// Извлекаем расширение файла из пути.
	ext := strings.ToUpper(filepath.Ext(path))
	// Удаляем точку из начала расширения.
	ext = strings.TrimPrefix(ext, ".")

	// Определяем, какому PhotoExtension соответствует извлеченное расширение.
	switch ext {
	case string(model.PhotoExtensionJpg), string(model.PhotoExtensionJpeg):
		photoExt := model.PhotoExtensionJpeg
		return &photoExt
	case string(model.PhotoExtensionPng):
		photoExt := model.PhotoExtensionPng
		return &photoExt
	case string(model.PhotoExtensionBmb):
		photoExt := model.PhotoExtensionBmb
		return &photoExt
	default:
		// Если расширение не соответствует известным типам, возвращаем nil.
		return nil
	}
}

func FileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
