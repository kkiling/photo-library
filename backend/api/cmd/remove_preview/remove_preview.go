package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// removeFilesByMask удаляет файлы в указанной директории, соответствующие маске.
func removeFilesByMask(dir, mask string) error {
	// Создаем паттерн поиска, комбинируя директорию и маску.
	pattern := filepath.Join(dir, mask)

	// Используем Glob для нахождения файлов, соответствующих паттерну.
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	// Проходимся по всем найденным файлам и удаляем их.
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			// Вместо прерывания, выводим ошибку и продолжаем с остальными файлами.
			fmt.Println("Ошибка при удалении файла:", file, err)
		} else {
			fmt.Println("Удалён файл:", file)
		}
	}

	return nil
}

func main() {
	dir := "D:\\photo_library" // Укажите нужный путь к директории
	mask := "preview_*"
	err := removeFilesByMask(dir, mask)
	if err != nil {
		fmt.Println("Ошибка при удалении файлов:", err)
	}
}
