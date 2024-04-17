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
