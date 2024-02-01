package geo

import (
	"strings"
)

const url = "https://nominatim.openstreetmap.org/"

type geocodeResponse struct {
	DisplayName string     `json:"display_name"`
	Addr        osmAddress `json:"address"`
}

func (r *geocodeResponse) Address() *Address {
	return &Address{
		FormattedAddress: r.DisplayName,
		HouseNumber:      r.Addr.HouseNumber,
		Street:           r.Addr.Street(),
		Postcode:         r.Addr.Postcode,
		City:             r.Addr.Locality(),
		Suburb:           r.Addr.Suburb,
		State:            r.Addr.State,
		Country:          r.Addr.Country,
		CountryCode:      strings.ToUpper(r.Addr.CountryCode),
	}
}

type osmAddress struct {
	HouseNumber   string `json:"house_number"`
	Road          string `json:"road"`
	Pedestrian    string `json:"pedestrian"`
	Footway       string `json:"footway"`
	Cycleway      string `json:"cycleway"`
	Highway       string `json:"highway"`
	Path          string `json:"path"`
	Suburb        string `json:"suburb"`
	City          string `json:"city"`
	Town          string `json:"town"`
	Village       string `json:"village"`
	Hamlet        string `json:"hamlet"`
	County        string `json:"county"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
	State         string `json:"state"`
	StateDistrict string `json:"state_district"`
	Postcode      string `json:"postcode"`
}

// Locality checks different fields for the locality name
func (a osmAddress) Locality() string {
	var locality string

	if a.City != "" {
		locality = a.City
	} else if a.Town != "" {
		locality = a.Town
	} else if a.Village != "" {
		locality = a.Village
	} else if a.Hamlet != "" {
		locality = a.Hamlet
	}

	return locality
}

// Street checks different fields for the street name
func (a osmAddress) Street() string {
	var street string

	if a.Road != "" {
		street = a.Road
	} else if a.Pedestrian != "" {
		street = a.Pedestrian
	} else if a.Path != "" {
		street = a.Path
	} else if a.Cycleway != "" {
		street = a.Cycleway
	} else if a.Footway != "" {
		street = a.Footway
	} else if a.Highway != "" {
		street = a.Highway
	}

	return street
}

type Address struct {
	FormattedAddress string
	Street           string
	HouseNumber      string
	Suburb           string
	Postcode         string
	State            string
	StateCode        string
	StateDistrict    string
	County           string
	Country          string
	CountryCode      string
	City             string
}
