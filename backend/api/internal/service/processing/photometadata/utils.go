package photometadata

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"

	"github.com/kkiling/photo-library/backend/api/internal/service/model"
)

func getImageDetails(photoBody []byte) (width int, height int, err error) {
	reader := bytes.NewReader(photoBody)
	// Получение размера изображения в пикселях
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, serviceerr.InvalidInputErr(err, "image.DecodeConfig")
	}

	return img.Width, img.Height, nil
}

func parseDate(s string) (time.Time, error) {
	if s == invalidTime {
		return time.Time{}, fmt.Errorf("invalid date format")
	}

	t, err := time.Parse(timeLayout, s)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

type regexToDate struct {
	pattern     *regexp.Regexp
	transformer func([]string) string
}

var rules = []regexToDate{
	{
		regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2}) (\d{2})-(\d{2})-(\d{2})`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^IMG-(\d{4})(\d{2})(\d{2})-WA\d+`),
		func(matches []string) string { return matches[1] + "-" + matches[2] + "-" + matches[3] },
	},
	{
		regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2}) (\d{2})-(\d{2})-(\d{2})_\d+`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})_(\d{2})(\d{2})(\d{2})`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^IMG_(\d{4})(\d{2})(\d{2})_(\d{2})(\d{2})(\d{2})`),
		func(matches []string) string {
			return matches[1] + "-" + matches[2] + "-" + matches[3] + " " + matches[4] + ":" + matches[5] + ":" + matches[6]
		},
	},
	{
		regexp.MustCompile(`^IMG_(\d{4})(\d{2})(\d{2})`),
		func(matches []string) string { return matches[1] + "-" + matches[2] + "-" + matches[3] },
	},
}

func fileNameToTime(path string) (time.Time, error) {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[0 : len(filename)-len(ext)]

	for _, rule := range rules {
		if matches := rule.pattern.FindStringSubmatch(nameWithoutExt); matches != nil {
			formattedDate := rule.transformer(matches)
			if len(formattedDate) == 10 { // "2006-01-02"
				return time.Parse("2006-01-02", formattedDate)
			}
			return time.Parse("2006-01-02 15:04:05", formattedDate)
		}
	}

	// Handle timestamp
	if ts, err := strconv.ParseInt(nameWithoutExt, 10, 64); err == nil {
		return time.Unix(ts/1000, 0), nil // convert from milliseconds
	}

	return time.Time{}, fmt.Errorf("could not parse date from filename")
}

// Convert GPS coordinates from degrees, minutes, and seconds format to decimal format
func dmsToDecimal(degrees float64, minutes float64, seconds float64) float64 {
	return degrees + minutes/60 + seconds/3600
}

func convertToGeo(latitude, longitude []float64) (*model.Geo, error) {
	if len(latitude) != 3 || len(longitude) != 3 {
		return nil, fmt.Errorf("invalid input. Ensure that each slice has 3 elements representing degrees, minutes, and seconds respectively")
	}

	return &model.Geo{
		Latitude:  dmsToDecimal(latitude[0], latitude[1], latitude[2]),
		Longitude: dmsToDecimal(longitude[0], longitude[1], longitude[2]),
	}, nil
}
