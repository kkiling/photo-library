package model

import "github.com/google/uuid"

// PhotoPreview Модель превью
type PhotoPreview struct {
	ID          uuid.UUID
	PhotoID     uuid.UUID
	FileKey     string
	WidthPixel  int
	HeightPixel int
	SizePixel   int
	Original    bool
}

// PhotoPreviewDTO Превью фотографии
type PhotoPreviewDTO struct {
	// Путь до фотографии, с нужным get параметром size
	Src string
	// Ширина фото
	Width int
	// Высота фото
	Height int
	// Размер фото (max Width Height)
	Size int
}

// PhotoWithPreviewsDTO оригинал фотографии с превью
type PhotoWithPreviewsDTO struct {
	PhotoID uuid.UUID
	// Оригинал фотографии
	Original PhotoPreviewDTO
	// Превью фотографии (разные размеры)
	Previews []PhotoPreviewDTO
}
