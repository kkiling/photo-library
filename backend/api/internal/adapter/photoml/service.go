package photoml

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

var ErrInternalServerError = errors.New("internal server error")

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

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, ErrInternalServerError
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	var vector []float64
	if unmarshalErr := json.Unmarshal(body, &vector); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return vector, nil
}
