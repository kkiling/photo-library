package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"io"
	"net/http"
	"strings"
)

type Service struct {
	logger log.Logger
	//mu     sync.Mutex
}

func NewService(logger log.Logger) *Service {
	return &Service{
		logger: logger,
		//mu:     sync.Mutex{},
	}
}

func reverseGeocodeURL(lat, lng float64) string {
	return url + "reverse?" + fmt.Sprintf("format=json&lat=%f&lon=%f", lat, lng)
}

// Response gets response from url
func response(ctx context.Context, url string, obj *geocodeResponse) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	body := strings.Trim(string(data), " []")
	//DebugLogger.Printf("Received response: %s\n", body)
	if body == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(body), obj); err != nil {
		//ErrLogger.Printf("Error unmarshalling response: %s\n", err.Error())
		return err
	}

	return nil
}

func (s *Service) ReverseGeocode(ctx context.Context, lat, lng float64) (*Address, error) {
	// s.mu.Lock()
	res := geocodeResponse{}
	geoUrl := reverseGeocodeURL(lat, lng)
	if err := response(ctx, geoUrl, &res); err != nil {
		return nil, fmt.Errorf("ReverseGeocode: %w", err)
	}
	// s.mu.Unlock()

	return res.Address(), nil
}
