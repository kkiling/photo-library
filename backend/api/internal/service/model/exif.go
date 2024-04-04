package model

import "github.com/google/uuid"

const (
	GPSLongitudeExifData     = "GPSLongitude"
	GPSLatitudeExifData      = "GPSLatitude"
	ModelExifData            = "Model"
	MakeExifData             = "Make"
	DateTimeExifData         = "DateTime"
	DateTimeOriginalExifData = "DateTimeOriginal"
)

type ExifPhotoData struct {
	PhotoID uuid.UUID
	Data    map[string]interface{}
}

func (e *ExifPhotoData) GetString(fileName string) (string, bool) {
	v, ok := e.Data[fileName]
	if !ok {
		return "", false
	}
	res, ok := v.(string)
	return res, ok
}

func (e *ExifPhotoData) GetFloat(fileName string) (float64, bool) {
	v, ok := e.Data[fileName]
	if !ok {
		return 0, false
	}
	res, ok := v.(float64)
	return res, ok
}

func (e *ExifPhotoData) GetFloatArray(fileName string) ([]float64, bool) {
	v, ok := e.Data[fileName]
	if !ok {
		return nil, false
	}
	resInterface, ok := v.([]interface{})

	res := make([]float64, 0, len(resInterface))

	for _, r := range resInterface {
		vv, ok := r.(float64)
		if ok {
			res = append(res, vv)
		} else {
			return nil, false
		}
	}

	return res, true
}
