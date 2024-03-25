package photopreview

import (
	"bytes"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
)

// Пример функции для поворота на 90 градусов по часовой стрелке
func rotate90CW(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(bounds.Max.Y-y-1, x, img.At(x, y))
		}
	}
	return dst
}

// Пример функции для поворота на 90 градусов против часовой стрелки
func rotate90CCW(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(y, bounds.Max.X-x-1, img.At(x, y))
		}
	}
	return dst
}

// Пример функции для отражения по горизонтали
func flipH(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(bounds.Max.X-x-1, y, img.At(x, y))
		}
	}
	return dst
}

// Пример функции для отражения по вертикали
func flipV(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, bounds.Max.Y-y-1, img.At(x, y))
		}
	}
	return dst
}

// Пример функции для поворота на 180 градусов
func rotate180(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
		}
	}
	return dst
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		// Нормальная ориентация, ничего не делаем
		return img
	case 2:
		// Отражение по горизонтали
		return flipH(img)
	case 3:
		// Поворот на 180 градусов
		return rotate180(img)
	case 4:
		// Отражение по вертикали
		return flipV(img)
	case 6:
		// Поворот на 90 градусов по часовой стрелке
		return rotate90CW(img)
	case 8:
		// Поворот на 90 градусов против часовой стрелки
		return rotate90CCW(img)
	default:
		// Неизвестная ориентация, не применяем изменения
		return img
	}
}

type imagePreview struct {
	photoBody []byte
	width     int
	height    int
}

func createImagePreview(originalImage image.Image, extension model.PhotoExtension, orientation int, maxSize int) (imagePreview, error) {
	// Вычисление новых размеров для сохранения пропорций
	originalWidth := originalImage.Bounds().Dx()
	originalHeight := originalImage.Bounds().Dy()

	newWidth, newHeight := originalWidth, originalHeight
	if originalWidth > maxSize || originalHeight > maxSize {
		if originalWidth > originalHeight {
			newWidth = maxSize
			newHeight = (originalHeight * maxSize) / originalWidth
		} else {
			newHeight = maxSize
			newWidth = (originalWidth * maxSize) / originalHeight
		}
	}

	// Создание нового изображения с новыми размерами
	targetRect := image.Rect(0, 0, newWidth, newHeight)
	resizedImg := image.NewRGBA(targetRect)
	draw.CatmullRom.Scale(resizedImg, targetRect, originalImage, originalImage.Bounds(), draw.Src, nil)

	// Применение ориентации к измененному по размеру изображению
	resizedAndOrientedImg := applyOrientation(resizedImg, orientation)

	// Кодирование в соответствующий формат
	var buf bytes.Buffer
	switch extension {
	case model.PhotoExtensionJpg, model.PhotoExtensionJpeg:
		if err := jpeg.Encode(&buf, resizedAndOrientedImg, nil); err != nil {
			return imagePreview{}, serviceerr.MakeErr(err, "failed to encode jpeg")
		}
	case model.PhotoExtensionPng:
		if err := png.Encode(&buf, resizedAndOrientedImg); err != nil {
			return imagePreview{}, serviceerr.MakeErr(err, "failed to encode png")
		}
	default:
		return imagePreview{}, serviceerr.NotFoundError("unsupported format")
	}

	if orientation == 6 || orientation == 8 {
		newWidth, newHeight = newHeight, newWidth
	}

	return imagePreview{
		photoBody: buf.Bytes(),
		width:     newWidth,
		height:    newHeight,
	}, nil
}
