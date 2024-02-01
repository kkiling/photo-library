package photoml

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Config struct {
	Url string `yaml:"url"`
}

type Service struct {
	logger log.Logger
	cfg    Config
}

func NewService(logger log.Logger, cfg Config) *Service {
	return &Service{
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Service) GetImageVector(ctx context.Context, imgBytes []byte) ([]float64, error) {
	// Создайте новый HTTP-запрос
	url := s.cfg.Url + "/get_vector_from_bytes"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(imgBytes))
	if err != nil {
		return nil, err
	}

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
