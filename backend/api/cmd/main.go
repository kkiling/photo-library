package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/floats"
	"io"
	"net/http"
	"os"
)

func GetFileBody(ctx context.Context, filePath string) ([]byte, error) {
	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create new file: %w", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create new file: %w", err)
	}

	return body, nil
}

func GetImageVector(ctx context.Context, imgBytes []byte) ([]float64, error) {
	// Создайте новый HTTP-запрос
	url := "http://127.0.0.1:5000/get_vector_from_bytes"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(imgBytes))
	if err != nil {
		return nil, err
	}

	// Вы можете установить дополнительные заголовки, если это необходимо
	req.Header.Set("Content-Type", "image/jpeg")

	// Выполните запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читайте ответ сервера
	body, _ := io.ReadAll(resp.Body)
	// Декодируем JSON-ответ в slice float64
	var vector []float64
	if err := json.Unmarshal(body, &vector); err != nil {
		return nil, err
	}

	return vector, nil
}

func similarity(vector1, vector2 []float64) float64 {
	d1 := floats.Norm(vector1, 2)              // L2 норма (евклидова норма) для vector1
	d2 := floats.Norm(vector2, 2)              // L2 норма для vector2
	dotProduct := floats.Dot(vector1, vector2) // Скалярное произведение двух векторов
	return dotProduct / (d1 * d2)
}

func main() {
	ctx := context.Background()
	//imgPath1 := "/Users/kkiling/photos/f983d019-5776-407e-b51c-0b26d4f4edd4.jpeg"
	imgPath1 := "/Users/kkiling/Desktop/vlad2016/IMG_20160501_134833.jpg"
	imgBytes1, err := GetFileBody(ctx, imgPath1)
	if err != nil {
		panic(err)
	}

	//imgPath2 := "/Users/kkiling/photos/f565797d-b650-4cfa-aa67-f811fe31c517.jpeg"
	imgPath2 := "/Users/kkiling/Desktop/vlad2016/IMG_20160501_134836.jpg"
	imgBytes2, err := GetFileBody(ctx, imgPath2)
	if err != nil {
		panic(err)
	}

	vector1, err := GetImageVector(ctx, imgBytes1)
	if err != nil {
		return
	}

	vector2, err := GetImageVector(ctx, imgBytes2)
	if err != nil {
		return
	}

	k := similarity(vector1, vector2)
	fmt.Println(k)
}
