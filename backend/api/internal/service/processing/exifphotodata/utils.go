package exifphotodata

import (
	"fmt"

	"github.com/rwcarlsen/goexif/tiff"
)

func ratToFloat(tag *tiff.Tag, i int) (float64, error) {
	r1, r2, err := tag.Rat2(i)
	if err != nil {
		return 0, err
	}
	if r2 == 0 {
		return 0, nil
	}
	return float64(r1) / float64(r2), nil
}

func getIntFromTag(tag *tiff.Tag) (int, error) {
	if tag.Format() == tiff.IntVal {
		return tag.Int(0)
	}
	return 0, fmt.Errorf("unexpected tag format")
}

func getFloatFromTag(tag *tiff.Tag) (float64, error) {
	switch tag.Format() {
	case tiff.RatVal:
		return ratToFloat(tag, 0)
	case tiff.FloatVal:
		return tag.Float(0)
	default:
		return 0, fmt.Errorf("unexpected tag format")
	}
}

func getIntArrayFromTag(tag *tiff.Tag) ([]int, error) {
	if tag.Format() != tiff.IntVal {
		return nil, fmt.Errorf("unexpected tag format")
	}
	res := make([]int, tag.Count)
	for i := 0; i < int(tag.Count); i++ {
		val, err := tag.Int(i)
		if err != nil {
			return nil, err
		}
		res[i] = val
	}
	return res, nil
}

func getFloatArrayFromTag(tag *tiff.Tag) ([]float64, error) {
	count := int(tag.Count)
	res := make([]float64, count)

	var err error
	for i := 0; i < count; i++ {
		switch tag.Format() {
		case tiff.RatVal:
			res[i], err = ratToFloat(tag, i)
		case tiff.FloatVal:
			res[i], err = tag.Float(i)
		default:
			return nil, fmt.Errorf("unexpected tag format")
		}
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
