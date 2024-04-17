package exif_photo_data

import (
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type write struct {
	data map[string]interface{}
}

func (p *write) Walk(name exif.FieldName, tag *tiff.Tag) error {
	fieldName := string(name)
	switch tag.Format() {
	case tiff.IntVal:
		if tag.Count == 1 {
			if val, err := getIntFromTag(tag); err == nil {
				p.data[fieldName] = val
			}
		} else {
			val, err := getIntArrayFromTag(tag)
			if err == nil {
				p.data[fieldName] = val
			}
		}
	case tiff.FloatVal:
		if tag.Count == 1 {
			val, err := getFloatFromTag(tag)
			if err == nil {
				p.data[fieldName] = val
			}
		} else {
			val, err := getFloatArrayFromTag(tag)
			if err == nil {
				p.data[fieldName] = val
			}
		}
	case tiff.RatVal:
		if tag.Count == 1 {
			val, err := getFloatFromTag(tag)
			if err == nil {
				p.data[fieldName] = val
			}
		} else {
			val, err := getFloatArrayFromTag(tag)
			if err == nil {
				p.data[fieldName] = val
			}
		}
	case tiff.StringVal:
		val := strings.TrimSpace(strings.Trim(tag.String(), `"`))
		p.data[fieldName] = val
	case tiff.UndefVal:
		val := strings.TrimSpace(strings.Trim(tag.String(), `"`))
		p.data[fieldName] = val
	case tiff.OtherVal:
		val := strings.TrimSpace(strings.Trim(tag.String(), `"`))
		p.data[fieldName] = val
	}

	return nil
}
