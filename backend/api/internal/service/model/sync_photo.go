package model

import "time"

// SyncPhotoRequest данные новой фотографии которые прислал клиент синхронизатор
type SyncPhotoRequest struct {
	// Пути фотографий которые загружаем (может быть несколько если фото одинаковые)
	Paths []string
	// Рассчитанный на клиенте хеш фотографии
	Hash string
	// Данные фото
	Body []byte
	// Информация о последнем изменении фото
	PhotoUpdatedAt time.Time
	// Информация о клиенте, загрузившего фотографии
	ClientInfo string
}

type SyncPhotoResponse struct {
	// Фото было загружено ранее
	HasBeenUploadedBefore bool
	// Рассчитанный на клиенте хеш фотографии
	Hash string
}
