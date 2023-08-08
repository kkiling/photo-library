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
	UpdateAt time.Time
	// Идентификатор клиента, загрузившего фотографии
	ClientId string
}

type SyncPhotoResponse struct {
	// Фото было загружено ранее
	HasBeenUploadedBefore bool
	// Рассчитанный на клиенте хеш фотографии
	Hash string
	// Дата когда фотография была загружена на сервер
	UploadAt time.Time
}
