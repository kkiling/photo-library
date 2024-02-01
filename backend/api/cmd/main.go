package main

import (
	"context"
	"fmt"
	"time"
)

type Geocoder struct {
	// Замените этот тип на ваш реальный тип Geocoder
}

func (s *Geocoder) ReverseGeocode(lat, lng float64) (string, error) {
	// Здесь должен быть ваш реальный код ReverseGeocode
	// Это просто заглушка
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("End ReverseGeocode")
	return fmt.Sprintf("Address at %.2f, %.2f", lat, lng), nil
}

func reverseGeocodeWithTimeout(ctx context.Context, s *Geocoder, lat, lng float64) (string, error) {
	// Создаем новый контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Создаем канал для получения результата и ошибки
	resultCh := make(chan reverseGeocodeResult, 1)

	// Запускаем выполнение функции в отдельной горутине
	go func() {
		address, err := s.ReverseGeocode(lat, lng)
		resultCh <- reverseGeocodeResult{address, err}
	}()

	select {
	case <-ctx.Done():
		// Время выполнения истекло
		return "", fmt.Errorf("timeout exceeded")
	case result := <-resultCh:
		// Получили результат
		return result.address, result.err
	}
}

type reverseGeocodeResult struct {
	address string
	err     error
}

func main() {
	// Замените на ваш реальный объект Geocoder
	s := &Geocoder{}

	// Установим таймаут в 1 секунду
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Вызываем ReverseGeocode с использованием контекста
	address, err := reverseGeocodeWithTimeout(ctx, s, 37.7749, -122.4194)
	if err != nil {
		fmt.Println("Error:", err)

		time.Sleep(time.Second)
		return
	}

	fmt.Println("Address:", address)
}
