package photo_groups_handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

// GetPhotoContent http метод для получения фотографии
func (p *PhotoGroupsHandler) GetPhotoContent(w http.ResponseWriter, r *http.Request) {
	fileKey := filepath.Base(r.URL.Path)
	photoContent, err := p.photosService.GetPhotoContent(r.Context(), fileKey)
	if err != nil {
		if serviceerr.IsNotFound(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, fmt.Errorf("p.photosService.GetPhotoContent: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	var contentType string
	switch photoContent.Extension {
	case model.PhotoExtensionJpeg:
		contentType = "image/jpeg"
	case model.PhotoExtensionPng:
		contentType = "image/png"
	default:
		http.Error(w, "Unsupported image format", http.StatusBadRequest)
		return
	}

	// Установка заголовка Content-Type и отправка изображения
	w.Header().Set("Content-Type", contentType)
	_, err = w.Write(photoContent.PhotoBody)
	if err != nil {
		p.logger.Errorf(" w.Write: %v", err)
		http.Error(w, "w.Write(photoContent.PhotoBody)", http.StatusInternalServerError)
	}
}
