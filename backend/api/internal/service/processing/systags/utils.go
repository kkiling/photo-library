package systags

import (
	"math"
	"path/filepath"
	"strings"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

const earthRadiusKm = 6371.0

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func getDirectories(path string) []string {
	catalogs := filepath.Dir(path)
	parts := strings.Split(catalogs, string(filepath.Separator))
	return parts
}

func distance(a, b model.Geo) float64 {
	return haversine(a.Latitude, a.Longitude, b.Latitude, b.Longitude)
}

func meanGeo(coordinates []model.Geo) model.Geo {
	var latSum, lonSum float64
	for _, coord := range coordinates {
		latSum += coord.Latitude
		lonSum += coord.Longitude
	}
	return model.Geo{Latitude: latSum / float64(len(coordinates)), Longitude: lonSum / float64(len(coordinates))}
}

func groupByMean(coordinates []model.Geo, radius float64) [][]model.Geo {
	var groups [][]model.Geo

	for len(coordinates) > 0 {
		currentCoord := coordinates[0]
		closeGroup := []model.Geo{currentCoord}

		// Удаляем текущую координату из списка
		coordinates = coordinates[1:]

		i := 0
		for i < len(coordinates) {
			if distance(currentCoord, coordinates[i]) <= radius {
				closeGroup = append(closeGroup, coordinates[i])

				// Удаляем эту координату из списка
				coordinates = append(coordinates[:i], coordinates[i+1:]...)
			} else {
				i++
			}
		}

		// Находим среднее значение для группы и добавляем его
		groups = append(groups, closeGroup)
	}

	return groups
}
