package main

import (
	"fmt"
	"github.com/codingsince1985/geo-golang/openstreetmap"
)

func main() {
	geocoder := openstreetmap.Geocoder()
	location, _ := geocoder.ReverseGeocode(55.878, 37.653)
	fmt.Printf("%+v\n", location)
}
