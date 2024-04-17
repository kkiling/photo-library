package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Service struct {
	logger log.Logger
	mu     sync.Mutex
}

func NewService(logger log.Logger) *Service {
	return &Service{
		logger: logger,
		mu:     sync.Mutex{},
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("geocodeService returned %s", resp.Status)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	body := strings.Trim(string(data), " []")
	if body == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(body), obj); err != nil {
		return err
	}

	return nil
}

func (s *Service) ReverseGeocode(ctx context.Context, lat, lng float64) (*Address, error) {
	// https://operations.osmfoundation.org/policies/nominatim/
	s.mu.Lock()
	defer s.mu.Unlock()

	res := geocodeResponse{}
	geoUrl := reverseGeocodeURL(lat, lng)
	if geoUrl == "https://nominatim.openstreetmap.org/reverse?format=json&lat=57.629967&lon=39.895862" {
		fmt.Printf("")
	}
	if err := response(ctx, geoUrl, &res); err != nil {
		return nil, fmt.Errorf("ReverseGeocode (%s): %w", geoUrl, err)
	}

	return res.Address(), nil
}
